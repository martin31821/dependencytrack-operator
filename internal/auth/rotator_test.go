/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

const (
	testNamespace  = "dependencytrack"
	testSecretName = "deptrack-credentials"
)

// fakeDependencyTrack implements the two endpoints the rotator calls. The
// first failForceChange/failLogin calls answer 503 to simulate an instance
// (or its TLS frontend) that is not ready yet.
type fakeDependencyTrack struct {
	mu               sync.Mutex
	password         string
	failForceChange  int
	failLogin        int
	forceChangeCalls int
}

func (f *fakeDependencyTrack) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch req.URL.Path {
	case "/api/v1/user/login":
		if f.failLogin > 0 {
			f.failLogin--
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if req.Form.Get("password") != f.password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("token"))
	case "/api/v1/user/forceChangePassword":
		f.forceChangeCalls++
		if f.failForceChange > 0 {
			f.failForceChange--
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if req.Form.Get("password") != f.password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		f.password = req.Form.Get("newPassword")
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (f *fakeDependencyTrack) currentPassword() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.password
}

var _ = Describe("PasswordRotationRunnable", func() {
	var (
		ctx    context.Context
		dt     *fakeDependencyTrack
		k8s    client.Client
		runner *PasswordRotationRunnable
	)

	readSecret := func() *corev1.Secret {
		s := &corev1.Secret{}
		Expect(k8s.Get(ctx, types.NamespacedName{Namespace: testNamespace, Name: testSecretName}, s)).To(Succeed())
		return s
	}

	startRunner := func() <-chan error {
		done := make(chan error, 1)
		go func() { done <- runner.Start(ctx) }()
		return done
	}

	BeforeEach(func() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		DeferCleanup(cancel)

		dt = &fakeDependencyTrack{password: "admin"}
		srv := httptest.NewServer(dt)
		DeferCleanup(srv.Close)

		Expect(os.Setenv("DEPTRACK_URL", srv.URL)).To(Succeed())
		Expect(os.Setenv("DEPTRACK_CREDENTIALS_SECRET", testSecretName)).To(Succeed())
		DeferCleanup(func() {
			_ = os.Unsetenv("DEPTRACK_URL")
			_ = os.Unsetenv("DEPTRACK_CREDENTIALS_SECRET")
		})

		k8s = fake.NewClientBuilder().Build()
		runner = &PasswordRotationRunnable{
			Client:    k8s,
			Namespace: testNamespace,
			Backoff:   wait.Backoff{Duration: 10 * time.Millisecond, Factor: 1, Steps: 1000},
		}
	})

	It("retries until DependencyTrack accepts the rotation", func() {
		dt.failForceChange = 3

		Eventually(startRunner()).Should(Receive(BeNil()))

		secret := readSecret()
		Expect(secret.Data).NotTo(HaveKey(secretKeyPasswordNew))
		Expect(len(secret.Data[secretKeyPassword])).To(BeNumerically(">=", minPasswordLength))
		Expect(string(secret.Data[secretKeyPassword])).To(Equal(dt.currentPassword()))
		Expect(dt.forceChangeCalls).To(Equal(4))
	})

	It("keeps password-new when the pending password cannot be checked", func() {
		// DependencyTrack already holds the pending password, but the first
		// login checks fail with 503. The rotator must not overwrite it.
		const pending = "pending-password-that-dt-already-accepted"
		dt.password = pending
		dt.failLogin = 2
		Expect(k8s.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: testSecretName, Namespace: testNamespace},
			Data: map[string][]byte{
				secretKeyUsername:    []byte("admin"),
				secretKeyPassword:    []byte("admin"),
				secretKeyPasswordNew: []byte(pending),
			},
		})).To(Succeed())

		Eventually(startRunner()).Should(Receive(BeNil()))

		secret := readSecret()
		Expect(string(secret.Data[secretKeyPassword])).To(Equal(pending))
		Expect(secret.Data).NotTo(HaveKey(secretKeyPasswordNew))
		Expect(dt.forceChangeCalls).To(BeZero())
	})

	It("rotates again when DependencyTrack rejects the pending password", func() {
		Expect(k8s.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: testSecretName, Namespace: testNamespace},
			Data: map[string][]byte{
				secretKeyUsername:    []byte("admin"),
				secretKeyPassword:    []byte("admin"),
				secretKeyPasswordNew: []byte("never-applied"),
			},
		})).To(Succeed())

		Eventually(startRunner()).Should(Receive(BeNil()))

		secret := readSecret()
		Expect(secret.Data).NotTo(HaveKey(secretKeyPasswordNew))
		Expect(string(secret.Data[secretKeyPassword])).To(Equal(dt.currentPassword()))
		Expect(dt.currentPassword()).NotTo(Equal("never-applied"))
	})

	It("stops retrying when the context is cancelled", func() {
		dt.failForceChange = 1 << 30
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)

		done := startRunner()
		Eventually(func() int {
			dt.mu.Lock()
			defer dt.mu.Unlock()
			return dt.forceChangeCalls
		}).Should(BeNumerically(">=", 2))
		cancel()

		Eventually(done).Should(Receive(BeNil()))
	})
})
