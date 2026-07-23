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
	"fmt"
	"net/http"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	handler "sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

const (
	notificationPublisherFinalizer       = "dependencytrack.mko.dev/finalizer"
	conditionReasonDependencyBlocked     = "DependencyBlocked"
	conditionReasonPublisherDeleteFailed = "PublisherDeleteFailed"
)

// NotificationPublisherReconciler reconciles a NotificationPublisher object.
type NotificationPublisherReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	DTProvider deptrack.ClientProviderInterface
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationpublishers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationpublishers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationpublishers/finalizers,verbs=update
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationrules,verbs=get;list;watch

func (r *NotificationPublisherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	publisher := &dependencytrackv1alpha1.NotificationPublisher{}
	if err := r.Get(ctx, req.NamespacedName, publisher); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure our finalizer is registered before doing any external work.
	if !controllerutil.ContainsFinalizer(publisher, notificationPublisherFinalizer) {
		controllerutil.AddFinalizer(publisher, notificationPublisherFinalizer)
		if err := r.Update(ctx, publisher); err != nil {
			return ctrl.Result{}, err
		}
		// Requeue; the update will trigger a new reconcile with the finalizer in place.
		return ctrl.Result{}, nil
	}

	if !publisher.DeletionTimestamp.IsZero() {
		result, err := r.reconcileDelete(ctx, publisher)
		return result, err
	}

	return r.reconcileUpsert(ctx, publisher)
}

func (r *NotificationPublisherReconciler) reconcileUpsert(ctx context.Context, publisher *dependencytrackv1alpha1.NotificationPublisher) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		setPublisherCondition(publisher, metav1.ConditionFalse, "CredentialsError",
			"failed to authenticate with DependencyTrack")
		_ = r.Status().Update(ctx, publisher)
		return ctrl.Result{}, err
	}

	var dtUUID string

	if publisher.Status.UUID != "" {
		// Publisher may already exist — find it by UUID in the list.
		// DependencyTrack v5 does not expose a GET-by-UUID endpoint for publishers,
		// so we list all and locate the matching UUID.
		existing, _, err := apiClient.NotificationAPI.GetAllNotificationPublishers(authCtx).Execute()
		if err != nil {
			return r.failPublisherStatus(ctx, publisher, "PublisherListFailed",
				"failed to list publishers from DependencyTrack: "+err.Error(), err)
		}

		found := false
		for _, pub := range existing {
			if pub.Uuid == publisher.Status.UUID {
				dtUUID = pub.Uuid
				found = true

				// Check for spec drift and update if needed.
				if pub.Name != publisher.Spec.Name || pub.ExtensionName != publisher.Spec.ExtensionName {
					log.Info("updating publisher due to spec drift", "uuid", pub.Uuid, "oldName", pub.Name, "newName", publisher.Spec.Name, "oldExt", pub.ExtensionName, "newExt", publisher.Spec.ExtensionName)
					updateReq := dtapi.UpdateNotificationPublisherRequest{
						Uuid:          pub.Uuid,
						ExtensionName: publisher.Spec.ExtensionName,
						Name:          publisher.Spec.Name,
					}
					if publisher.Spec.Description != "" {
						updateReq.Description = &publisher.Spec.Description
					}
					if _, _, err := apiClient.NotificationAPI.UpdateNotificationPublisher(authCtx).UpdateNotificationPublisherRequest(updateReq).Execute(); err != nil {
						return r.failPublisherStatus(ctx, publisher, "PublisherUpdateFailed",
							"failed to update publisher: "+err.Error(), err)
					}
				}
				break
			}
		}

		if !found {
			// Publisher was deleted out-of-band; recreate it.
			log.Info("publisher not found in DependencyTrack, recreating", "uuid", publisher.Status.UUID)
		}
	}

	if dtUUID == "" {
		createReq := dtapi.CreateNotificationPublisherRequest{
			ExtensionName: publisher.Spec.ExtensionName,
			Name:          publisher.Spec.Name,
		}
		if publisher.Spec.Description != "" {
			createReq.Description = &publisher.Spec.Description
		}
		created, httpRes, err := apiClient.NotificationAPI.CreateNotificationPublisher(authCtx).CreateNotificationPublisherRequest(createReq).Execute()
		if err != nil {
			// 409 Conflict means a publisher with the same extension already exists.
			// Try to adopt the existing publisher instead of failing outright.
			if httpRes != nil && httpRes.StatusCode == http.StatusConflict {
				log.Info("publisher already exists in DependencyTrack (409), attempting to adopt",
					"extensionName", publisher.Spec.ExtensionName, "name", publisher.Spec.Name)
				adopted, adoptErr := r.adoptExistingPublisher(authCtx, apiClient, publisher)
				if adoptErr != nil {
					return r.failPublisherStatus(ctx, publisher, "PublisherCreateFailed",
						fmt.Sprintf("failed to create publisher in DependencyTrack: 409 Conflict (publisher already exists) and adoption failed: %v", adoptErr), adoptErr)
				}
				dtUUID = adopted
				log.Info("adopted existing publisher in DependencyTrack", "uuid", dtUUID)
				r.Recorder.Eventf(publisher, "Normal", "PublisherAdopted",
					"Adopted existing publisher %q in DependencyTrack (uuid=%s)", publisher.Spec.Name, dtUUID)
			} else {
				return r.failPublisherStatus(ctx, publisher, "PublisherCreateFailed",
					"failed to create publisher in DependencyTrack: "+err.Error(), err)
			}
		} else {
			dtUUID = created.Uuid
			log.Info("created publisher in DependencyTrack", "uuid", dtUUID)
			r.Recorder.Eventf(publisher, "Normal", "PublisherCreated", "Created publisher %q in DependencyTrack (uuid=%s)", publisher.Spec.Name, dtUUID)
		}
	}

	publisher.Status.UUID = dtUUID
	publisher.Status.Name = publisher.Spec.Name
	publisher.Status.ObservedGeneration = publisher.Generation

	setPublisherCondition(publisher, metav1.ConditionTrue, "PublisherSynced", "Publisher successfully reconciled in DependencyTrack")
	if err := r.Status().Update(ctx, publisher); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// hasDependentRules checks whether any same-namespace NotificationRule
// references this publisher by its spec.name via PublisherRef.Name.
// Returns the names of dependent rules when found so that status conditions
// can identify the blocking resource safely.
func (r *NotificationPublisherReconciler) hasDependentRules(ctx context.Context, publisher *dependencytrackv1alpha1.NotificationPublisher) (bool, string, error) {
	log := logf.FromContext(ctx)
	rules := &dependencytrackv1alpha1.NotificationRuleList{}
	if err := r.List(ctx, rules, client.InNamespace(publisher.Namespace)); err != nil {
		log.Error(err, "failed to list notification rules for dependency check")
		return false, "", err
	}
	var dependentNames []string
	for i := range rules.Items {
		rule := &rules.Items[i]
		if rule.Spec.PublisherRef.Name == publisher.Name {
			log.Info("found dependent rule blocking publisher deletion",
				"ruleNamespace", rule.Namespace, "ruleName", rule.Name)
			dependentNames = append(dependentNames, rule.Name)
		}
	}
	if len(dependentNames) == 0 {
		return false, "", nil
	}
	return true, strings.Join(dependentNames, ", "), nil
}

func (r *NotificationPublisherReconciler) reconcileDelete(ctx context.Context, publisher *dependencytrackv1alpha1.NotificationPublisher) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Check for same-namespace NotificationRules that reference this publisher.
	// Block deletion while any dependent rule exists to prevent orphaning remote rules.
	hasDeps, dependentRules, err := r.hasDependentRules(ctx, publisher)
	if err != nil {
		return ctrl.Result{}, err
	}
	if hasDeps {
		setPublisherCondition(publisher, metav1.ConditionFalse,
			conditionReasonDependencyBlocked,
			"deletion blocked: same-namespace NotificationRule(s) \""+dependentRules+"\" reference this publisher")
		if updateErr := r.Status().Update(ctx, publisher); updateErr != nil {
			log.Error(updateErr, "failed to update publisher status during dependency-block")
			return ctrl.Result{}, updateErr
		}
		log.Info("publisher deletion blocked by dependent rule(s)", "publisher", publisher.Name, "namespace", publisher.Namespace, "dependentRules", dependentRules)
		// Requeue periodically so the controller detects when a dependent
		// rule is deleted.  The controller does not watch NotificationRules
		// directly; it must actively poll to notice deletion.
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if publisher.Status.UUID != "" {
		authCtx, apiClient, err := r.DTProvider.Get(ctx)
		if err != nil {
			log.Error(err, "failed to authenticate with DependencyTrack during deletion")
			return ctrl.Result{}, err
		}

		resp, err := apiClient.NotificationAPI.DeleteNotificationPublisher(authCtx, publisher.Status.UUID).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				log.Info("publisher already absent from DependencyTrack (404)", "uuid", publisher.Status.UUID)
			} else {
				log.Error(err, "failed to delete publisher from DependencyTrack", "uuid", publisher.Status.UUID)
				setPublisherCondition(publisher, metav1.ConditionFalse,
					conditionReasonPublisherDeleteFailed,
					"failed to delete publisher from DependencyTrack")
				_ = r.Status().Update(ctx, publisher)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("deleted publisher from DependencyTrack", "uuid", publisher.Status.UUID)
		}
		r.Recorder.Eventf(publisher, "Normal", "PublisherDeleted", "Deleted publisher %q from DependencyTrack (uuid=%s)", publisher.Spec.Name, publisher.Status.UUID)
	}

	controllerutil.RemoveFinalizer(publisher, notificationPublisherFinalizer)
	if err := r.Update(ctx, publisher); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// adoptExistingPublisher lists all publishers from DependencyTrack and finds
// one matching the desired extensionName. Returns its UUID, or an error if
// no matching publisher is found.
func (r *NotificationPublisherReconciler) adoptExistingPublisher(ctx context.Context, apiClient *dtapi.APIClient, publisher *dependencytrackv1alpha1.NotificationPublisher) (string, error) {
	log := logf.FromContext(ctx)

	allPubs, _, err := apiClient.NotificationAPI.GetAllNotificationPublishers(ctx).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list publishers for adoption: %w", err)
	}

	// Prefer exact name+extension match.
	for _, pub := range allPubs {
		if pub.ExtensionName == publisher.Spec.ExtensionName && pub.Name == publisher.Spec.Name {
			log.Info("found existing publisher matching name+extension",
				"uuid", pub.Uuid, "name", pub.Name, "extension", pub.ExtensionName)
			return pub.Uuid, nil
		}
	}

	// Fallback: match by extensionName only (handles built-in publishers
	// like the default webhook that DT creates with a system-assigned name).
	for _, pub := range allPubs {
		if pub.ExtensionName == publisher.Spec.ExtensionName {
			log.Info("adopted existing publisher by extensionName (name mismatch)",
				"uuid", pub.Uuid, "dtName", pub.Name, "k8sName", publisher.Spec.Name,
				"extension", pub.ExtensionName)
			return pub.Uuid, nil
		}
	}

	return "", fmt.Errorf("no existing publisher found with extensionName %q", publisher.Spec.ExtensionName)
}

// failPublisherStatus sets a failed condition, persists the status, and returns
// the error so the reconcile loop retries.
func (r *NotificationPublisherReconciler) failPublisherStatus(ctx context.Context, publisher *dependencytrackv1alpha1.NotificationPublisher, reason, msg string, cause error) (ctrl.Result, error) {
	logf.FromContext(ctx).Error(cause, msg)
	setPublisherCondition(publisher, metav1.ConditionFalse, reason, msg)
	_ = r.Status().Update(ctx, publisher)
	return ctrl.Result{}, cause
}

// setPublisherCondition upserts the Ready condition on the publisher's status.
func setPublisherCondition(publisher *dependencytrackv1alpha1.NotificationPublisher, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&publisher.Status.Conditions, metav1.Condition{
		Type:               conditionReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: publisher.Generation,
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *NotificationPublisherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.NotificationPublisher{}).
		Watches(
			&dependencytrackv1alpha1.NotificationRule{},
			handler.EnqueueRequestsFromMapFunc(func(_ context.Context, obj client.Object) []ctrl.Request {
				rule, ok := obj.(*dependencytrackv1alpha1.NotificationRule)
				if !ok {
					return nil
				}
				if rule.Spec.PublisherRef.Name == "" {
					return nil
				}
				return []ctrl.Request{{
					NamespacedName: types.NamespacedName{
						Name:      rule.Spec.PublisherRef.Name,
						Namespace: rule.Namespace,
					},
				}}
			}),
		).
		Named("notificationpublisher").
		Complete(r)
}
