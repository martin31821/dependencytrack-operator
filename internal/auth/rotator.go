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
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

const (
	secretKeyUsername    = "username"
	secretKeyPassword    = "password"
	secretKeyPasswordNew = "password-new"
	minPasswordLength    = 30
)

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;create;update;patch

// PasswordRotationRunnable checks and rotates the DependencyTrack admin
// password stored in a Kubernetes Secret. It runs once per leader election.
type PasswordRotationRunnable struct {
	Client    client.Client
	Namespace string
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

	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: r.Namespace, Name: secretName}, secret); err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "failed to read credentials secret", "secret", secretName)
			return nil
		}
		log.Info("credentials secret not found, creating with default admin:admin")
		if err := r.createDefaultSecret(ctx, secretName); err != nil {
			log.Error(err, "failed to create credentials secret")
			return nil
		}
		// Re-read after creation so the rest of the flow sees the new secret.
		if err := r.Client.Get(ctx, types.NamespacedName{Namespace: r.Namespace, Name: secretName}, secret); err != nil {
			log.Error(err, "failed to read credentials secret after creation", "secret", secretName)
			return nil
		}
	}

	username := string(secret.Data[secretKeyUsername])
	password := string(secret.Data[secretKeyPassword])
	pendingPassword := string(secret.Data[secretKeyPasswordNew])

	apiClient := deptrack.NewAPIClient(deptrackURL)

	// Crash recovery: if password-new exists a previous rotation was interrupted.
	// Try the pending password against DependencyTrack first.
	if pendingPassword != "" {
		log.Info("detected in-progress rotation, attempting to resume")
		if _, _, err := apiClient.UserAPI.ValidateCredentials(ctx).
			Username(username).
			Password(pendingPassword).
			Execute(); err == nil {
			log.Info("DependencyTrack already accepted the new password, finalizing secret")
			return r.finalizeSecret(ctx, secret, pendingPassword)
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
		log.Error(err, "failed to generate new password")
		return nil
	}

	// Step A: persist the new password as a temporary key so we can recover if
	// the operator crashes before finishing.
	patch := secret.DeepCopy()
	patch.Data[secretKeyPasswordNew] = []byte(newPassword)
	if err := r.Client.Patch(ctx, patch, client.MergeFrom(secret)); err != nil {
		log.Error(err, "failed to store pending password in secret")
		return nil
	}
	secret = patch

	// Step B: update the password in DependencyTrack.
	if _, err := apiClient.UserAPI.ForceChangePassword(ctx).
		Username(username).
		Password(password).
		NewPassword(newPassword).
		ConfirmPassword(newPassword).
		Execute(); err != nil {
		log.Error(err, "failed to update password in DependencyTrack; secret retains password-new for recovery")
		return nil
	}

	// Step C: promote the new password to the main key and remove the temp key.
	return r.finalizeSecret(ctx, secret, newPassword)
}

func (r *PasswordRotationRunnable) finalizeSecret(ctx context.Context, secret *corev1.Secret, newPassword string) error {
	log := logf.FromContext(ctx).WithName("password-rotation")

	patch := secret.DeepCopy()
	patch.Data[secretKeyPassword] = []byte(newPassword)
	delete(patch.Data, secretKeyPasswordNew)
	if err := r.Client.Patch(ctx, patch, client.MergeFrom(secret)); err != nil {
		log.Error(err, "failed to finalize secret after password rotation")
		return nil
	}

	log.Info("password rotation completed successfully")
	return nil
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
