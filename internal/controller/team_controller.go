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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

const (
	teamFinalizer       = "dependencytrack.mko.dev/finalizer"
	conditionReconciled = "Reconciled"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	DTProvider *deptrack.ClientProvider
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams/finalizers,verbs=update

func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	team := &dependencytrackv1alpha1.Team{}
	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure our finalizer is registered before doing any external work.
	if !controllerutil.ContainsFinalizer(team, teamFinalizer) {
		controllerutil.AddFinalizer(team, teamFinalizer)
		if err := r.Update(ctx, team); err != nil {
			return ctrl.Result{}, err
		}
		// Requeue; the update will trigger a new reconcile with the finalizer in place.
		return ctrl.Result{}, nil
	}

	if !team.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, team)
	}

	return r.reconcileUpsert(ctx, team)
}

func (r *TeamReconciler) reconcileUpsert(ctx context.Context, team *dependencytrackv1alpha1.Team) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		setCondition(team, metav1.ConditionFalse, "CredentialsError", err.Error())
		_ = r.Status().Update(ctx, team)
		return ctrl.Result{}, err
	}

	var dtUUID string

	if team.Status.UUID != "" {
		// Team may already exist — fetch it by the UUID we recorded last time.
		existing, httpResp, err := apiClient.TeamAPI.GetTeam(authCtx, team.Status.UUID).Execute()
		switch {
		case err == nil:
			// Team still exists; rename if the spec changed.
			if existing.Name != team.Spec.Name {
				log.Info("updating team name", "uuid", existing.Uuid, "oldName", existing.Name, "newName", team.Spec.Name)
				updated, _, err := apiClient.TeamAPI.UpdateTeam(authCtx).Team(dtapi.Team{
					Name: team.Spec.Name,
					Uuid: existing.Uuid,
				}).Execute()
				if err != nil {
					return r.failStatus(ctx, team, "APIError", "failed to update team name: "+err.Error(), err)
				}
				dtUUID = updated.Uuid
			} else {
				dtUUID = existing.Uuid
			}
		case httpResp != nil && httpResp.StatusCode == http.StatusNotFound:
			// Team was deleted from DependencyTrack out-of-band; recreate it.
			log.Info("team not found in DependencyTrack, recreating", "uuid", team.Status.UUID)
		default:
			return r.failStatus(ctx, team, "APIError", "failed to get team from DependencyTrack: "+err.Error(), err)
		}
	}

	if dtUUID == "" {
		createTeam := dtapi.Team{
			Name: team.Spec.Name,
		}
		created, _, err := apiClient.TeamAPI.CreateTeam(authCtx).Team(createTeam).Execute()
		if err != nil {
			return r.failStatus(ctx, team, "APIError", "failed to create team: "+err.Error(), err)
		}
		dtUUID = created.Uuid
		log.Info("created team in DependencyTrack", "uuid", dtUUID)
		r.Recorder.Eventf(team, "Normal", "TeamCreated", "Created team %q in DependencyTrack (uuid=%s)", team.Spec.Name, dtUUID)
	}

	team.Status.UUID = dtUUID
	setCondition(team, metav1.ConditionTrue, "TeamSynced", "Team successfully reconciled in DependencyTrack")
	if err := r.Status().Update(ctx, team); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TeamReconciler) reconcileDelete(ctx context.Context, team *dependencytrackv1alpha1.Team) error {
	log := logf.FromContext(ctx)

	if team.Status.UUID != "" {
		authCtx, apiClient, err := r.DTProvider.Get(ctx)
		if err != nil {
			log.Error(err, "failed to authenticate with DependencyTrack during deletion")
			return err
		}

		httpResp, err := apiClient.TeamAPI.DeleteTeam(authCtx).Team(dtapi.Team{
			Name: team.Spec.Name,
			Uuid: team.Status.UUID,
		}).Execute()
		if err != nil && (httpResp == nil || httpResp.StatusCode != http.StatusNotFound) {
			log.Error(err, "failed to delete team from DependencyTrack", "uuid", team.Status.UUID)
			return err
		}
		log.Info("deleted team from DependencyTrack", "uuid", team.Status.UUID)
		r.Recorder.Eventf(team, "Normal", "TeamDeleted", "Deleted team %q from DependencyTrack (uuid=%s)", team.Spec.Name, team.Status.UUID)
	}

	controllerutil.RemoveFinalizer(team, teamFinalizer)
	return r.Update(ctx, team)
}

// failStatus sets a failed condition, persists the status, and returns the error so the reconcile loop requeues.
func (r *TeamReconciler) failStatus(ctx context.Context, team *dependencytrackv1alpha1.Team, reason, msg string, cause error) (ctrl.Result, error) {
	logf.FromContext(ctx).Error(cause, msg)
	setCondition(team, metav1.ConditionFalse, reason, msg)
	_ = r.Status().Update(ctx, team)
	return ctrl.Result{}, cause
}

// setCondition upserts the Reconciled condition on the team's status.
func setCondition(team *dependencytrackv1alpha1.Team, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&team.Status.Conditions, metav1.Condition{
		Type:               conditionReconciled,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: team.Generation,
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.Team{}).
		Named("team").
		Complete(r)
}
