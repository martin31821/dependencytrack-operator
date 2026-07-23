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

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

// ---- notification-pkg-level mocks (namespaced to avoid collision with policy_controller_test.go) ----

// notificationMockClientProvider implements deptrack.ClientProviderInterface without needing a K8s client.
type notificationMockClientProvider struct {
	url    string
	mu     sync.Mutex
	api    *dtapi.APIClient
	getErr error
}

func (p *notificationMockClientProvider) Get(ctx context.Context) (context.Context, *dtapi.APIClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.getErr != nil {
		return nil, nil, p.getErr
	}

	if p.api == nil {
		p.api = deptrack.NewAPIClient(p.url)
	}
	return context.WithValue(ctx, dtapi.ContextAccessToken, "mock-token"), p.api, nil
}

func (p *notificationMockClientProvider) Invalidate() {
	p.mu.Lock()
	p.api = nil
	p.mu.Unlock()
}

// notificationMockDTServer is a standalone httptest server for notification API endpoints.
type notificationMockDTServer struct {
	mu                sync.Mutex
	publisher         *dtapi.NotificationPublisher
	publishers        map[string]*dtapi.NotificationPublisher // uuid -> publisher
	rules             map[string]*dtapi.NotificationRule      // uuid -> rule
	failCreate        bool
	failList          bool
	failUpdate        bool
	failDelete        bool
	lastCreateRequest *dtapi.CreateNotificationPublisherRequest
}

const (
	notifPublisherPath       = "/v1/notification/publisher"
	notifPublisherPathPrefix = notifPublisherPath + "/"
	notifValidateCredsPath   = "/v1/user/validate-credentials"
)

func (s *notificationMockDTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// The generated client sends paths like /api/v1/... (server URL includes /api).
	// The route constants don't include /api, so strip it.
	stripped, found := strings.CutPrefix(r.URL.Path, "/api")
	if !found {
		stripped = r.URL.Path
	}

	switch {
	case stripped == notifPublisherPath && (r.Method == http.MethodPost || r.Method == http.MethodPut):
		s.handleCreateOrUpdate(w, r)
	case stripped == notifPublisherPath && r.Method == http.MethodGet:
		s.handleList(w)
	case strings.HasPrefix(stripped, notifPublisherPathPrefix) && r.Method == http.MethodDelete:
		s.handleDelete(w, r)
	case stripped == "/v1/notification/rule" && r.Method == http.MethodGet:
		s.handleRuleList(w)
	case strings.HasPrefix(stripped, "/v1/notification/rule/") && r.Method == http.MethodDelete:
		s.handleRuleDelete(w, r)
	case stripped == notifValidateCredsPath:
		writeJSON(w, "mock-token")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *notificationMockDTServer) handleCreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	// The generated UpdateNotificationPublisherRequest sends PUT /v1/notification/publisher
	// with the UUID in the body.  CreateNotificationPublisherRequest sends PUT to the
	// same path with no UUID in the body.  We peek at the raw body to distinguish.
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer func() { r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) }()

	// Try creating a CreateNotificationPublisherRequest first; it has no uuid field.
	// If that succeeds and the decoded struct has an empty UUID, it's a create.
	var createReq dtapi.CreateNotificationPublisherRequest
	if err := json.Unmarshal(bodyBytes, &createReq); err == nil {
		// CreateNotificationPublisherRequest has no uuid field, so any "uuid" key in
		// the JSON would be silently discarded.  We must therefore also try the
		// UpdateNotificationPublisherRequest to see if the body contains a uuid.
		var updateReq dtapi.UpdateNotificationPublisherRequest
		if json.Unmarshal(bodyBytes, &updateReq) == nil && updateReq.Uuid != "" {
			// Body contains a uuid — it's an update.
			s.handleUpdateFromBody(w, bytes.NewReader(bodyBytes))
			return
		}
		// No uuid detected — treat as create.
		s.handleCreateFromBody(w, bytes.NewReader(bodyBytes))
		return
	}
	// Not a valid create request — must be an update.
	s.handleUpdateFromBody(w, bytes.NewReader(bodyBytes))
}

func (s *notificationMockDTServer) handleCreateFromBody(w http.ResponseWriter, r io.Reader) {
	if s.failCreate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.CreateNotificationPublisherRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.lastCreateRequest = &req
	pub := &dtapi.NotificationPublisher{
		Name:          req.Name,
		ExtensionName: req.ExtensionName,
		Uuid:          "pub-created-uuid-123",
	}
	if req.Description != nil {
		pub.Description = req.Description
	}
	s.publisher = pub
	s.publishers[pub.Uuid] = pub
	writeJSON(w, pub)
}

func (s *notificationMockDTServer) handleUpdateFromBody(w http.ResponseWriter, r io.Reader) {
	if s.failUpdate {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var req dtapi.UpdateNotificationPublisherRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uuid := req.Uuid
	pub, exists := s.publishers[uuid]
	if !exists {
		// Fallback: try the single-slot publisher (used by tests that pre-seed
		// only mockServer.publisher without putting it in the map).
		pub = s.publisher
		if pub == nil || pub.Uuid != uuid {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	pub.Name = req.Name
	pub.ExtensionName = req.ExtensionName
	if req.Description != nil {
		pub.Description = req.Description
	}
	// Always sync the single-slot publisher when the UUID matches,
	// even if we found it in the map.
	if s.publisher == nil || s.publisher.Uuid != uuid {
		s.publisher = pub
	}
	s.publishers[uuid] = pub
	writeJSON(w, pub)
}

func (s *notificationMockDTServer) handleList(w http.ResponseWriter) {
	if s.failList {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	list := make([]dtapi.NotificationPublisher, 0, len(s.publishers))
	for _, pub := range s.publishers {
		list = append(list, *pub)
	}
	writeJSON(w, list)
}

func (s *notificationMockDTServer) handleRuleList(w http.ResponseWriter) {
	list := make([]dtapi.NotificationRule, 0, len(s.rules))
	for _, rule := range s.rules {
		list = append(list, *rule)
	}
	writeJSON(w, list)
}

func (s *notificationMockDTServer) handleRuleDelete(w http.ResponseWriter, r *http.Request) {
	stripped, _ := strings.CutPrefix(r.URL.Path, "/api")
	path := "/" + stripped
	uuid, _ := strings.CutPrefix(path, "/v1/notification/rule/")

	if _, exists := s.rules[uuid]; exists {
		delete(s.rules, uuid)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *notificationMockDTServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	if s.failDelete {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// The generated client calls DELETE /v1/notification/publisher/{uuid}.
	stripped, _ := strings.CutPrefix(r.URL.Path, "/api")
	uuid, _ := strings.CutPrefix(stripped, notifPublisherPathPrefix)

	// The controller passes publisher.Status.UUID to the generated API.
	// Check the publishers map first.
	if _, exists := s.publishers[uuid]; exists {
		delete(s.publishers, uuid)
		if s.publisher != nil && s.publisher.Uuid == uuid {
			s.publisher = nil
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Fallback: single-slot publisher.
	if s.publisher != nil && s.publisher.Uuid == uuid {
		s.publisher = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// ---- Ginkgo tests ----

const publisherNS = "default"

var _ = Describe("NotificationPublisher Controller", func() {
	var (
		server       *httptest.Server
		mockServer   *notificationMockDTServer
		provider     *notificationMockClientProvider
		ctx          context.Context
		fakeRecorder *record.FakeRecorder
		ctrl         *NotificationPublisherReconciler
	)

	BeforeEach(func() {
		mockServer = &notificationMockDTServer{
			publishers: make(map[string]*dtapi.NotificationPublisher),
			rules:      make(map[string]*dtapi.NotificationRule),
		}
		server = httptest.NewServer(mockServer)
		provider = &notificationMockClientProvider{url: server.URL}
		DeferCleanup(server.Close)

		fakeRecorder = record.NewFakeRecorder(10)
		ctrl = &NotificationPublisherReconciler{
			Client:     k8sClient,
			Scheme:     k8sClient.Scheme(),
			Recorder:   fakeRecorder,
			DTProvider: provider,
		}
		ctx = context.Background()
	})

	// Helpers
	createNP := func(name, ns string) *dependencytrackv1alpha1.NotificationPublisher {
		return &dependencytrackv1alpha1.NotificationPublisher{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
			},
			Spec: dependencytrackv1alpha1.NotificationPublisherSpec{
				ExtensionName: "webhook",
				Name:          name + "-ext",
			},
		}
	}
	createNPWithFinalizer := func(name, ns string) *dependencytrackv1alpha1.NotificationPublisher {
		p := createNP(name, ns)
		controllerutil.AddFinalizer(p, notificationPublisherFinalizer)
		return p
	}
	setNPUUID := func(name, ns, uuid string) {
		p := &dependencytrackv1alpha1.NotificationPublisher{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, p)).To(Succeed())
		p.Status.UUID = uuid
		Expect(k8sClient.Status().Update(ctx, p)).To(Succeed())
	}
	deleteNP := func(name, ns string) {
		p := &dependencytrackv1alpha1.NotificationPublisher{}
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, p); err == nil {
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())
		}
	}

	// ---- Finalizer ----

	Context("When reconciling a new resource", func() {
		const name = "basic-publisher"
		BeforeEach(func() {
			p := createNP(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
		})
		It("should add the finalizer on first reconcile", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(p, notificationPublisherFinalizer)).To(BeTrue())
		})
	})

	// ---- Create ----

	Context("When reconciling with no DT UUID (create)", func() {
		const name = "create-publisher"
		BeforeEach(func() {
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
		})
		It("should create the publisher in DT and set UUID", func() {
			// 1st reconcile = finalizer (envtest strips it on Create).
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			// 2nd reconcile = actual DT work.
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))
			Expect(p.Status.Name).To(Equal("create-publisher-ext"))
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PublisherSynced"))
			Eventually(fakeRecorder.Events).Should(Receive(ContainSubstring("PublisherCreated")))
		})
	})

	// ---- Update (noop — fields match) ----

	Context("When reconciling with an existing UUID that matches spec (noop)", func() {
		const name = "noop-publisher"
		BeforeEach(func() {
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "noop-publisher-ext", ExtensionName: "webhook", Uuid: "existing-uuid-match",
			}
			mockServer.publishers["existing-uuid-match"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "existing-uuid-match")
		})
		It("should not call update when fields already match", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("existing-uuid-match"))
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(cond.Reason).To(Equal("PublisherSynced"))
			Consistently(fakeRecorder.Events).ShouldNot(Receive())
		})
	})

	// ---- Update (drift) ----

	Context("When reconciling with a UUID and drifted spec fields", func() {
		const name = "drift-publisher"
		BeforeEach(func() {
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "old-name-ext", ExtensionName: "email", Uuid: "drift-uuid-abc",
			}
			mockServer.publishers["drift-uuid-abc"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			p.Spec.ExtensionName = "webhook"
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "drift-uuid-abc")
		})
		It("should update the publisher in DT when spec drifts", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("drift-uuid-abc"))
			Expect(p.Status.Name).To(Equal("drift-publisher-ext"))
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			mockServer.mu.Lock()
			defer mockServer.mu.Unlock()
			Expect(mockServer.publisher.Name).To(Equal("drift-publisher-ext"))
			Expect(mockServer.publisher.ExtensionName).To(Equal("webhook"))
		})
	})

	// ---- Recreate (404) ----

	Context("When the publisher was deleted from DT (recreate)", func() {
		const name = "recreate-publisher"
		BeforeEach(func() {
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "deleted-uuid-xyz")
		})
		It("should recreate the publisher when DT returns 404", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))
		})
	})

	// ---- Delete ----

	Context("When a publisher is deleted (finalizer path)", func() {
		It("should remove the finalizer after successful DT deletion", func() {
			const name = "del-success-pub"
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "del-pub-ext", ExtensionName: "webhook", Uuid: "del-success-uuid",
			}
			mockServer.publishers["del-success-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			setNPUUID(name, publisherNS, "del-success-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			Eventually(fakeRecorder.Events).Should(Receive(ContainSubstring("PublisherDeleted")))
		})

		It("should recreate the publisher when DT does not have it and keep the finalizer", func() {
			const name = "recreate-404-pub"
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			setNPUUID(name, publisherNS, "nonexistent-uuid-999")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			// Controller recreates missing DT publisher — finalizer stays.
			Expect(controllerutil.ContainsFinalizer(obj, notificationPublisherFinalizer)).To(BeTrue())
			Expect(obj.Status.UUID).To(Equal("pub-created-uuid-123"))
			cond := meta.FindStatusCondition(obj.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		})

		It("should retain the finalizer when DT delete fails", func() {
			const name = "del-fail-pub"
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "fail-pub-ext", ExtensionName: "webhook", Uuid: "del-fail-uuid",
			}
			mockServer.publishers["del-fail-uuid"] = mockServer.publisher
			mockServer.failDelete = true
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			setNPUUID(name, publisherNS, "del-fail-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			obj = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationPublisherFinalizer)).To(BeTrue())
		})
	})

	// ---- Negative tests ----

	Context("Negative tests — failure modes", func() {
		It("should report CredentialsError when auth fails", func() {
			const name = "neg-auth-fail"
			badProvider := &notificationMockClientProvider{getErr: fmt.Errorf("auth failed: unreachable")}
			badCtrl := &NotificationPublisherReconciler{
				Client: k8sClient, Scheme: k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(5), DTProvider: badProvider,
			}
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			_, err := badCtrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("CredentialsError"))
		})

		It("should report PublisherListFailed when LIST returns HTTP 500", func() {
			const name = "neg-list-fail"
			mockServer.failList = true
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "list-fail-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherListFailed"))
		})

		It("should report PublisherCreateFailed when CREATE returns HTTP 500", func() {
			const name = "neg-create-fail"
			mockServer.failCreate = true
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			// Envtest does not strip finalizers, so the first reconcile goes
			// straight to reconcileUpsert.  Since UUID is empty the controller
			// tries to create and hits the 500.
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherCreateFailed"))
		})

		It("should report PublisherUpdateFailed when UPDATE returns HTTP 500", func() {
			const name = "neg-update-fail"
			mockServer.mu.Lock()
			mockServer.failUpdate = true
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "old-updated-name", ExtensionName: "old-ext", Uuid: "update-fail-uuid",
			}
			mockServer.publishers["update-fail-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "update-fail-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal("PublisherUpdateFailed"))
		})

		It("should retain the finalizer when a transient DT delete error occurs", func() {
			const name = "neg-transient-delete"
			mockServer.mu.Lock()
			mockServer.failDelete = true
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "neg-del-ext", ExtensionName: "webhook", Uuid: "neg-del-uuid",
			}
			mockServer.publishers["neg-del-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "neg-del-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			obj = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationPublisherFinalizer)).To(BeTrue())
		})

		It("should report PublisherDeleteFailed when DT returns HTTP 500 during cleanup", func() {
			const name = "neg-delete-fail"
			mockServer.mu.Lock()
			mockServer.failDelete = true
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "neg-delete-fail-ext", ExtensionName: "webhook", Uuid: "neg-delete-fail-uuid",
			}
			mockServer.publishers["neg-delete-fail-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "neg-delete-fail-uuid")
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).To(HaveOccurred())
			// Verify the finalizer is retained (reconcile error means delete didn't complete)
			obj = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(obj, notificationPublisherFinalizer)).To(BeTrue())
		})
	})

	// ---- Description field ----

	Context("When a publisher has a description set", func() {
		const name = "desc-publisher"
		BeforeEach(func() {
			p := createNPWithFinalizer(name, publisherNS)
			p.Spec.Description = "Test publisher description"
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
		})
		It("should pass the description to the DT create request", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))
			Expect(mockServer.lastCreateRequest).NotTo(BeNil())
			Expect(mockServer.lastCreateRequest.Description).ToNot(BeNil())
			Expect(*mockServer.lastCreateRequest.Description).To(Equal("Test publisher description"))
		})
	})

	// ---- Dependency-blocked deletion ----

	Context("When a publisher has a dependent NotificationRule", func() {
		It("should block deletion and set DependencyBlocked condition", func() {
			const name = "dep-blocked-pub"
			// Pre-seed publisher in mock DT
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: "dep-pub-ext", ExtensionName: "webhook", Uuid: "dep-blocked-uuid",
			}
			mockServer.publishers["dep-blocked-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "dep-blocked-uuid")
			// Create a referencing rule in the same namespace
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ref-rule-" + name,
					Namespace: publisherNS,
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:         "ref-rule-name",
					Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:        dependencytrackv1alpha1.NotificationRuleLevelInfo,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: name},
				},
			}
			Expect(k8sClient.Create(ctx, rule)).To(Succeed())
			DeferCleanup(func() {
				_ = k8sClient.Delete(ctx, rule)
			})
			// Reconcile — should find the publisher and reconcile Upsert successfully
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			// Mark for deletion — should be blocked by the dependent rule
			Expect(k8sClient.Delete(ctx, p)).To(Succeed())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			// Verify: finalizer retained, condition set
			p = &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(p, notificationPublisherFinalizer)).To(BeTrue())
			cond = meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal(conditionReasonDependencyBlocked))
			// The DT publisher should still exist — no delete request was sent
			mockServer.mu.Lock()
			_, exists := mockServer.publishers["dep-blocked-uuid"]
			mockServer.mu.Unlock()
			Expect(exists).To(BeTrue())
		})

		It("should allow deletion after the dependent rule is retargeted", func() {
			const pubName = "dep-release-pub"
			// Pre-seed publisher in mock DT
			mockServer.mu.Lock()
			mockServer.publisher = &dtapi.NotificationPublisher{
				Name: pubName + "-ext", ExtensionName: "webhook", Uuid: "dep-release-uuid",
			}
			mockServer.publishers["dep-release-uuid"] = mockServer.publisher
			mockServer.mu.Unlock()
			p := createNPWithFinalizer(pubName, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, pubName, publisherNS)
			setNPUUID(pubName, publisherNS, "dep-release-uuid")
			// Create a referencing rule
			rule := &dependencytrackv1alpha1.NotificationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ref-rule-" + pubName,
					Namespace: publisherNS,
				},
				Spec: dependencytrackv1alpha1.NotificationRuleSpec{
					Name:         "ref-rule-name",
					Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
					TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
					Level:        dependencytrackv1alpha1.NotificationRuleLevelInfo,
					PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: pubName},
				},
			}
			Expect(k8sClient.Create(ctx, rule)).To(Succeed())
			DeferCleanup(func() {
				_ = k8sClient.Delete(ctx, rule)
			})
			// Retarget the rule to a different publisher
			rule.Spec.PublisherRef.Name = "other-publisher"
			Expect(k8sClient.Update(ctx, rule)).To(Succeed())
			// Reconcile the publisher to ensure it's properly set up (drift-check passes since UUID already set)
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: pubName, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			// Verify the publisher UUID is set before deleting
			pub := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: pubName, Namespace: publisherNS}, pub)).To(Succeed())
			Expect(pub.Status.UUID).To(Equal("dep-release-uuid"), "K8s publisher UUID should match mock server")
			// Delete the publisher
			Expect(k8sClient.Delete(ctx, pub)).To(Succeed())
			// Reconcile — should now succeed
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: pubName, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			// Verify DT publisher was deleted (GC will remove K8s object after finalizer removal)
			mockServer.mu.Lock()
			_, exists := mockServer.publishers["dep-release-uuid"]
			mockServer.mu.Unlock()
			Expect(exists).To(BeFalse())
			// After finalizer removal the K8s object is garbage-collected, so Get returns NotFound
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: pubName, Namespace: publisherNS}, obj)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})

		It("should treat remote 404 as successful deletion", func() {
			const name = "absent-pub"
			// Publisher exists in K8s but not in DT
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
			setNPUUID(name, publisherNS, "absent-uuid")
			// Delete the publisher from K8s
			obj := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, obj)).To(Succeed())
			Expect(k8sClient.Delete(ctx, obj)).To(Succeed())
			// Reconcile — should treat 404 as success and remove finalizer
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			// After finalizer removal the K8s object is garbage-collected
			np := &dependencytrackv1alpha1.NotificationPublisher{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, np)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})
	})

	// ---- ObservedGeneration ----

	Context("ObservedGeneration is set on Ready status", func() {
		const name = "obs-gen-publisher"
		BeforeEach(func() {
			p := createNPWithFinalizer(name, publisherNS)
			Expect(k8sClient.Create(ctx, p)).To(Succeed())
			DeferCleanup(deleteNP, name, publisherNS)
		})
		It("should set ObservedGeneration in the Ready condition", func() {
			_, err := ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			_, err = ctrl.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: publisherNS}})
			Expect(err).NotTo(HaveOccurred())
			p := &dependencytrackv1alpha1.NotificationPublisher{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: publisherNS}, p)).To(Succeed())
			Expect(p.Status.UUID).To(Equal("pub-created-uuid-123"))
			cond := meta.FindStatusCondition(p.Status.Conditions, conditionReady)
			Expect(cond).NotTo(BeNil())
			Expect(cond.ObservedGeneration).To(Equal(p.Generation))
		})
	})
})
