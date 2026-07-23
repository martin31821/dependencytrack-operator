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
	"context"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

const (
	apiKeyFinalizer  = "dependencytrack.mko.dev/finalizer"
	secretKeyAPIKey  = "api-key"
	teamNotReadyWait = 15 * time.Second
)

// APIKeyReconciler reconciles a APIKey object
type APIKeyReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	DTProvider *deptrack.ClientProvider
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=apikeys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=apikeys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=apikeys/finalizers,verbs=update
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *APIKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	apiKey := &dependencytrackv1alpha1.APIKey{}
	if err := r.Get(ctx, req.NamespacedName, apiKey); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(apiKey, apiKeyFinalizer) {
		controllerutil.AddFinalizer(apiKey, apiKeyFinalizer)
		if err := r.Update(ctx, apiKey); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !apiKey.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, apiKey)
	}

	return r.reconcileUpsert(ctx, apiKey)
}

// teamUUID looks up the referenced Team CR and returns its DependencyTrack UUID.
// It returns ("", requeue, error): when requeue is true the caller should
// requeue after teamNotReadyWait without treating it as a hard error.
func (r *APIKeyReconciler) teamUUID(ctx context.Context, apiKey *dependencytrackv1alpha1.APIKey) (string, bool, error) {
	log := logf.FromContext(ctx)

	team := &dependencytrackv1alpha1.Team{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: apiKey.Namespace, Name: apiKey.Spec.TeamRef}, team); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("referenced Team not found, requeueing", "team", apiKey.Spec.TeamRef)
			setAPIKeyCondition(apiKey, metav1.ConditionFalse, "TeamNotFound",
				"referenced Team "+apiKey.Spec.TeamRef+" does not exist")
			_ = r.Status().Update(ctx, apiKey)
			return "", true, nil
		}
		return "", false, err
	}

	// Block until the Team has been successfully reconciled.
	if !meta.IsStatusConditionTrue(team.Status.Conditions, conditionReady) {
		log.Info("referenced Team not yet reconciled, requeueing", "team", apiKey.Spec.TeamRef)
		setAPIKeyCondition(apiKey, metav1.ConditionFalse, "TeamNotReady",
			"referenced Team "+apiKey.Spec.TeamRef+" has not been successfully reconciled yet")
		_ = r.Status().Update(ctx, apiKey)
		return "", true, nil
	}

	if team.Status.UUID == "" {
		log.Info("referenced Team has no UUID yet, requeueing", "team", apiKey.Spec.TeamRef)
		setAPIKeyCondition(apiKey, metav1.ConditionFalse, "TeamNotReady",
			"referenced Team "+apiKey.Spec.TeamRef+" has not reported its UUID yet")
		_ = r.Status().Update(ctx, apiKey)
		return "", true, nil
	}

	return team.Status.UUID, false, nil
}

func (r *APIKeyReconciler) reconcileUpsert(ctx context.Context, apiKey *dependencytrackv1alpha1.APIKey) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	teamUUID, requeue, err := r.teamUUID(ctx, apiKey)
	if err != nil {
		return ctrl.Result{}, err
	}
	if requeue {
		return ctrl.Result{RequeueAfter: teamNotReadyWait}, nil
	}

	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		setAPIKeyCondition(apiKey, metav1.ConditionFalse, "CredentialsError", err.Error())
		_ = r.Status().Update(ctx, apiKey)
		return ctrl.Result{}, err
	}

	var publicID string

	if apiKey.Status.PublicID != "" {
		// Verify the key still exists by checking the team's key list.
		team, httpResp, err := apiClient.TeamAPI.GetTeam(authCtx, teamUUID).Execute()
		switch {
		case err != nil && httpResp != nil && httpResp.StatusCode == http.StatusNotFound:
			// Team itself is gone — the TeamReconciler will sort that out; requeue.
			log.Info("team not found in DependencyTrack, requeueing")
			return ctrl.Result{RequeueAfter: teamNotReadyWait}, nil
		case err != nil:
			return r.failAPIKeyStatus(ctx, apiKey, "APIError", "failed to get team from DependencyTrack: "+err.Error(), err)
		default:
			found := false
			for _, k := range team.GetApiKeys() {
				if k.GetPublicId() == apiKey.Status.PublicID {
					found = true
					break
				}
			}
			if found {
				publicID = apiKey.Status.PublicID
			} else {
				log.Info("API key no longer present in DependencyTrack, recreating", "publicId", apiKey.Status.PublicID)
			}
		}
	}

	if publicID == "" {
		// Generate a new API key for the team.
		created, _, err := apiClient.TeamAPI.GenerateApiKey(authCtx, teamUUID).Execute()
		if err != nil {
			log.Error(err, "failed to generate API key in DependencyTrack")
			return r.failAPIKeyStatus(ctx, apiKey, "APIError", "failed to generate API key: "+err.Error(), err)
		}
		publicID = created.GetPublicId()
		log.Info("generated API key", "publicId", publicID)
		r.Recorder.Eventf(apiKey, "Normal", "APIKeyCreated", "Generated API key (publicId=%s) for team %s", publicID, apiKey.Spec.TeamRef)

		// Persist the key value in the target Secret immediately so it is not lost.
		if err := r.upsertSecret(ctx, apiKey, created.GetKey()); err != nil {
			// The key exists in DT but we failed to store it — record publicID so we
			// don't generate a second orphaned key on the next reconcile.
			apiKey.Status.PublicID = publicID
			setAPIKeyCondition(apiKey, metav1.ConditionFalse, "SecretWriteError", err.Error())
			_ = r.Status().Update(ctx, apiKey)
			return ctrl.Result{}, err
		}
	}

	// Sync the comment if one is specified.
	if apiKey.Spec.Comment != "" {
		if _, _, err := apiClient.TeamAPI.UpdateApiKeyComment(authCtx, publicID).Body(apiKey.Spec.Comment).Execute(); err != nil {
			return r.failAPIKeyStatus(ctx, apiKey, "APIError", "failed to update API key comment: "+err.Error(), err)
		}
	}

	apiKey.Status.PublicID = publicID
	setAPIKeyCondition(apiKey, metav1.ConditionTrue, "KeySynced", "API key successfully reconciled in DependencyTrack")
	if err := r.Status().Update(ctx, apiKey); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *APIKeyReconciler) reconcileDelete(ctx context.Context, apiKey *dependencytrackv1alpha1.APIKey) error {
	log := logf.FromContext(ctx)

	if apiKey.Status.PublicID != "" {
		authCtx, apiClient, err := r.DTProvider.Get(ctx)
		if err != nil {
			log.Error(err, "failed to authenticate with DependencyTrack during deletion")
			return err
		}

		httpResp, err := apiClient.TeamAPI.DeleteApiKey(authCtx, apiKey.Status.PublicID).Execute()
		if err != nil && (httpResp == nil || httpResp.StatusCode != http.StatusNotFound) {
			log.Error(err, "failed to delete API key from DependencyTrack", "publicId", apiKey.Status.PublicID)
			return err
		}
		log.Info("deleted API key from DependencyTrack", "publicId", apiKey.Status.PublicID)
		r.Recorder.Eventf(apiKey, "Normal", "APIKeyDeleted", "Deleted API key (publicId=%s)", apiKey.Status.PublicID)
	}

	// Remove the managed Secret.
	secret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Namespace: apiKey.Namespace, Name: apiKey.Spec.SecretName}, secret)
	if err == nil {
		if err := r.Delete(ctx, secret); err != nil && !apierrors.IsNotFound(err) {
			log.Error(err, "failed to delete API key secret", "secret", apiKey.Spec.SecretName)
			return err
		}
	} else if !apierrors.IsNotFound(err) {
		return err
	}

	controllerutil.RemoveFinalizer(apiKey, apiKeyFinalizer)
	return r.Update(ctx, apiKey)
}

// upsertSecret creates or updates the target Secret with the API key value.
func (r *APIKeyReconciler) upsertSecret(ctx context.Context, apiKey *dependencytrackv1alpha1.APIKey, keyValue string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiKey.Spec.SecretName,
			Namespace: apiKey.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if err := controllerutil.SetControllerReference(apiKey, secret, r.Scheme); err != nil {
			return err
		}
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		secret.Data[secretKeyAPIKey] = []byte(keyValue)
		return nil
	})
	return err
}

// failAPIKeyStatus records a failed condition, persists the status, and returns the error.
func (r *APIKeyReconciler) failAPIKeyStatus(ctx context.Context, apiKey *dependencytrackv1alpha1.APIKey, reason, msg string, cause error) (ctrl.Result, error) {
	logf.FromContext(ctx).Error(cause, msg)
	setAPIKeyCondition(apiKey, metav1.ConditionFalse, reason, msg)
	_ = r.Status().Update(ctx, apiKey)
	return ctrl.Result{}, cause
}

func setAPIKeyCondition(apiKey *dependencytrackv1alpha1.APIKey, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&apiKey.Status.Conditions, metav1.Condition{
		Type:               conditionReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: apiKey.Generation,
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.APIKey{}).
		Owns(&corev1.Secret{}).
		Named("apikey").
		Complete(r)
}
