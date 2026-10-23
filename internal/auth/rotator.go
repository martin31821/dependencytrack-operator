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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/martin31821/dependencytrack-operator/internal/dependencytrack"
)

const (
	secretKeyUsername    = "username"
	secretKeyPassword    = "password"
	secretKeyPasswordNew = "password-new"
	minPasswordLength    = 30
)

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;create;update;patch

// defaultRotationBackoff retries quickly at first (DependencyTrack or its TLS
// certificate is often still coming up right after install) and settles at
// one attempt every five minutes.
var defaultRotationBackoff = wait.Backoff{
	Duration: 5 * time.Second,
	Factor:   2,
	Jitter:   0.1,
	Steps:    math.MaxInt32,
	Cap:      5 * time.Minute,
}

// PasswordRotationRunnable checks and rotates the DependencyTrack admin
// password stored in a Kubernetes Secret. It runs once per leader election and
// retries with exponential backoff until the rotation succeeds or the context
// is cancelled. When ClientProvider is set, the runnable invalidates its auth
// cache after rotation so controllers immediately pick up the new credentials.
type PasswordRotationRunnable struct {
	Client    client.Client
	Namespace string
	// ClientProvider, if set, is invalidated after rotation so controllers
	// re-authenticate with the new password on their next reconciliation.
	ClientProvider dependencytrack.ClientProviderInterface
	// Backoff controls the retry schedule for failed attempts. The zero value
	// selects defaultRotationBackoff.
	Backoff wait.Backoff
}

func (r *PasswordRotationRunnable) NeedLeaderElection() bool { return true }

func (r *PasswordRotationRunnable) Start(ctx context.Context) error {
	log := logf.FromContext(ctx).WithName("password-rotation")

	deptrackURL := os.Getenv("DEPTRACK_URL")
	secretName := os.Getenv("DEPTRACK_CREDENTIALS_SECRET")
	if deptrackURL == "" || secretName == "" {
		log.Error(fmt.Errorf("missing required env vars"), "DEPTRACK_URL and DEPTRACK_CREDENTIALS_SECRET must be set")
		return nil
	}

	backoff := r.Backoff
	if backoff.Steps == 0 {
		backoff = defaultRotationBackoff
	}

	for attempt := 1; ; attempt++ {
		err := r.rotate(ctx, deptrackURL, secretName)
		if err == nil {
			return nil
		}
		delay := backoff.Step()
		log.Error(err, "password rotation attempt failed, retrying", "attempt", attempt, "retryIn", delay.String())
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
		}
	}
}

// rotate performs one rotation attempt. It re-reads the Secret on every call,
// so a retry resumes from whatever state the previous attempt left behind.
// A non-nil error means the attempt should be retried.
func (r *PasswordRotationRunnable) rotate(ctx context.Context, deptrackURL, secretName string) error {
	log := logf.FromContext(ctx).WithName("password-rotation")

	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: r.Namespace, Name: secretName}, secret); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("read credentials secret %s: %w", secretName, err)
		}
		log.Info("credentials secret not found, creating with default admin:admin")
		if err := r.createDefaultSecret(ctx, secretName); err != nil {
			return err
		}
		// We just created the secret, so populate it directly instead of
		// re-reading (the cache informer is not yet running at this point in
		// manager startup).
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: r.Namespace,
			},
			Data: map[string][]byte{
				secretKeyUsername: []byte("admin"),
				secretKeyPassword: []byte("admin"),
			},
		}
	}

	username := string(secret.Data[secretKeyUsername])
	password := string(secret.Data[secretKeyPassword])
	pendingPassword := string(secret.Data[secretKeyPasswordNew])

	apiClient := dependencytrack.NewAPIClient(deptrackURL)

	// Crash recovery: if password-new exists a previous rotation was interrupted.
	// Try the pending password against DependencyTrack first.
	if pendingPassword != "" {
		log.Info("detected in-progress rotation, attempting to resume")
		_, resp, err := apiClient.UserAPI.ValidateCredentials(ctx).
			Username(username).
			Password(pendingPassword).
			Execute()
		if err == nil {
			log.Info("DependencyTrack already accepted the new password, finalizing secret")
			if err := r.finalizeSecret(ctx, secret, pendingPassword); err != nil {
				return err
			}
			r.invalidateProvider()
			return nil
		}
		// Only an explicit rejection proves the pending password was never
		// applied. On any other failure (TLS, network, 5xx) DependencyTrack may
		// already hold it, and overwriting password-new would lose the only
		// valid credential, so retry the check instead.
		if resp == nil || resp.StatusCode != http.StatusUnauthorized {
			return fmt.Errorf("check pending password against DependencyTrack: %w", err)
		}
		log.Info("new password not yet applied to DependencyTrack, retrying rotation")
		// Fall through and rotate again using the original password.
	}

	if len(password) >= minPasswordLength {
		log.Info("password meets length requirement, no rotation needed")
		return nil
	}

	log.Info("password shorter than minimum, starting rotation", "minLength", minPasswordLength)

	newPassword, err := generatePassword()
	if err != nil {
		return fmt.Errorf("generate new password: %w", err)
	}

	// Step A: persist the new password as a temporary key so we can recover if
	// the operator crashes before finishing.
	patch := secret.DeepCopy()
	if patch.Data == nil {
		patch.Data = map[string][]byte{}
	}
	patch.Data[secretKeyPasswordNew] = []byte(newPassword)
	if err := r.Client.Patch(ctx, patch, client.MergeFrom(secret)); err != nil {
		return fmt.Errorf("store pending password in secret: %w", err)
	}
	secret = patch

	// Step B: update the password in DependencyTrack.
	if _, err := apiClient.UserAPI.ForceChangePassword(ctx).
		Username(username).
		Password(password).
		NewPassword(newPassword).
		ConfirmPassword(newPassword).
		Execute(); err != nil {
		return fmt.Errorf("update password in DependencyTrack (secret retains password-new for recovery): %w", err)
	}

	// Step C: promote the new password to the main key and remove the temp key.
	if err := r.finalizeSecret(ctx, secret, newPassword); err != nil {
		return err
	}
	r.invalidateProvider()
	return nil
}

func (r *PasswordRotationRunnable) finalizeSecret(ctx context.Context, secret *corev1.Secret, newPassword string) error {
	log := logf.FromContext(ctx).WithName("password-rotation")

	patch := secret.DeepCopy()
	patch.Data[secretKeyPassword] = []byte(newPassword)
	delete(patch.Data, secretKeyPasswordNew)
	if err := r.Client.Patch(ctx, patch, client.MergeFrom(secret)); err != nil {
		return fmt.Errorf("finalize secret after password rotation: %w", err)
	}

	log.Info("password rotation completed successfully")
	return nil
}

// invalidateProvider clears the cached auth token so controllers re-authenticate
// with the newly rotated password on their next reconciliation.
func (r *PasswordRotationRunnable) invalidateProvider() {
	if r.ClientProvider != nil {
		r.ClientProvider.Invalidate()
	}
}

// generatePassword returns a cryptographically random URL-safe string
// long enough to satisfy the minimum length requirement.
func generatePassword() (string, error) {
	// base64 encoding expands 3 bytes → 4 chars, so ceil(minPasswordLength * 3 / 4)
	// bytes gives us at least minPasswordLength encoded chars.
	n := (minPasswordLength*3 + 3) / 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// createDefaultSecret creates a credentials secret with admin:admin
// so the operator can bootstrap a connection to a fresh DependencyTrack
// instance. The password-rotation logic will then upgrade the password.
func (r *PasswordRotationRunnable) createDefaultSecret(ctx context.Context, name string) error {
	log := logf.FromContext(ctx).WithName("password-rotation")

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: r.Namespace,
		},
		Data: map[string][]byte{
			secretKeyUsername: []byte("admin"),
			secretKeyPassword: []byte("admin"),
		},
	}
	if err := r.Client.Create(ctx, secret); err != nil {
		return fmt.Errorf("create credentials secret: %w", err)
	}
	log.Info("created credentials secret with default admin:admin")
	return nil
}
