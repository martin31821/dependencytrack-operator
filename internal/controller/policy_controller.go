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

	"github.com/google/uuid"
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
	policyFinalizer = "dependencytrack.mko.dev/policy-finalizer"
	conditionReady  = "Ready"

	// Reason values for the Ready condition.
	reasonCredentialsError      = "CredentialsError"
	reasonPolicyListFailed      = "PolicyListFailed"
	reasonPolicyGetFailed       = "PolicyGetFailed"
	reasonPolicyCreateFailed    = "PolicyCreateFailed"
	reasonPolicyUpdateFailed    = "PolicyUpdateFailed"
	reasonConditionCheckFailed  = "ConditionCheckFailed"
	reasonConditionCreateFailed = "ConditionCreateFailed"
	reasonOwnershipConflict     = "OwnershipConflict"
	reasonPolicySynced          = "PolicySynced"
	reasonPolicyDeleteFailed    = "PolicyDeleteFailed"
)

// PolicyReconciler reconciles a Policy object.
type PolicyReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	DTProvider deptrack.ClientProviderInterface
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=policies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=policies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=policies/finalizers,verbs=update

// Reconcile implements the controller reconcile loop.
func (r *PolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	policy := &dependencytrackv1alpha1.Policy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add finalizer if not present.
	if !controllerutil.ContainsFinalizer(policy, policyFinalizer) {
		controllerutil.AddFinalizer(policy, policyFinalizer)
		if err := r.Update(ctx, policy); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("added finalizer to policy", "name", policy.Name)
		return ctrl.Result{}, nil
	}

	// Handle deletion.
	if !policy.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, policy)
	}

	return r.reconcileUpsert(ctx, policy)
}

func (r *PolicyReconciler) reconcileUpsert(ctx context.Context, policy *dependencytrackv1alpha1.Policy) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Authenticate with DependencyTrack.
	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonCredentialsError,
			"Failed to authenticate with DependencyTrack: "+err.Error())
		return ctrl.Result{}, err
	}

	// Get or create the remote policy.
	dtPolicy, err := r.getOrCreateDTPolicy(ctx, authCtx, apiClient, policy)
	if err != nil {
		// Ownership-conflict is a non-retryable policy violation.
		if strings.Contains(err.Error(), "already owned by") {
			r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonOwnershipConflict,
				err.Error())
			log.Error(err, "policy name is already owned in DependencyTrack — refusing to create",
				"policyName", policy.Spec.Name)
			return ctrl.Result{}, nil // non-retryable, stop reconciling
		}
		return ctrl.Result{}, err
	}

	// Persist the confirmed UUID so that condition sync can look up the policy.
	policy.Status.UUID = dtPolicy.GetUuid()
	policy.Status.Name = policy.Spec.Name

	// Sync inline conditions.
	if err := r.syncPolicyConditions(ctx, authCtx, apiClient, policy, dtPolicy); err != nil {
		policy.Status.UUID = ""
		policy.Status.Name = ""
		return ctrl.Result{}, err
	}

	r.setReadyCondition(ctx, policy, metav1.ConditionTrue, reasonPolicySynced,
		fmt.Sprintf("Policy %q successfully reconciled in DependencyTrack (uuid=%s)",
			policy.Spec.Name, dtPolicy.GetUuid()))

	if err := r.Status().Update(ctx, policy); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// getOrCreateDTPolicy fetches the existing DT policy by recorded UUID, checks
// for same-name conflicts when no UUID is set, or creates a new global policy.
func (r *PolicyReconciler) getOrCreateDTPolicy(
	ctx context.Context,
	authCtx context.Context,
	apiClient *dtapi.APIClient,
	policy *dependencytrackv1alpha1.Policy,
) (*dtapi.Policy, error) {
	log := logf.FromContext(ctx)

	var dtPolicy *dtapi.Policy

	if policy.Status.UUID != "" {
		existing, httpResp, err := apiClient.PolicyAPI.GetPolicy(authCtx, policy.Status.UUID).Execute()
		switch {
		case err == nil:
			dtPolicy = existing
		case httpResp != nil && httpResp.StatusCode == http.StatusNotFound:
			log.Info("policy not found in DependencyTrack, will recreate", "uuid", policy.Status.UUID)
			dtPolicy = nil
		default:
			r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonPolicyGetFailed,
				"Failed to get policy from DependencyTrack: "+err.Error())
			_ = r.Status().Update(ctx, policy)
			return nil, fmt.Errorf("getting policy from DependencyTrack: %w", err)
		}
	}

	if dtPolicy == nil {
		// No UUID-owned policy found. If we have no stored UUID, check for
		// same-name conflicts before creating — this prevents accidental
		// duplicate policies when a DT operator created a policy with the
		// same name manually or from another controller.
		if policy.Status.UUID == "" {
			allPolicies, _, err := apiClient.PolicyAPI.GetPolicies(authCtx).Execute()
			if err != nil {
				r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonPolicyListFailed,
					"Failed to list policies from DependencyTrack: "+err.Error())
				_ = r.Status().Update(ctx, policy)
				return nil, fmt.Errorf("listing policies from DependencyTrack: %w", err)
			}
			if nameExists, conflictUUID := findPolicyByName(allPolicies, policy.Spec.Name); nameExists {
				return nil, fmt.Errorf("policy name %q is already owned by DependencyTrack policy uuid=%s", policy.Spec.Name, conflictUUID)
			}
		}

		// Create a new global policy.
		policyUUID := uuid.New().String()
		dtPolicy = &dtapi.Policy{
			Name:           policy.Spec.Name,
			Global:         boolPtr(true),
			Operator:       string(policy.Spec.Operator),
			Uuid:           policyUUID,
			ViolationState: string(policy.Spec.ViolationState),
		}
		created, httpResp, err := apiClient.PolicyAPI.CreatePolicy(authCtx).Policy(*dtPolicy).Execute()
		if err != nil {
			msg := err.Error()
			if gerr, ok := err.(interface{ Body() []byte }); ok && len(gerr.Body()) > 0 {
				msg = string(gerr.Body())
			}
			srvStatus := 0
			if httpResp != nil {
				srvStatus = httpResp.StatusCode
			}
			log.Error(err, "create policy failed", "status", srvStatus, "body", msg)
			r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonPolicyCreateFailed,
				"Failed to create policy in DependencyTrack: "+msg)
			return nil, fmt.Errorf("creating policy in DependencyTrack: %w", err)
		}
		dtPolicy = created
		log.Info("created policy in DependencyTrack", "uuid", dtPolicy.GetUuid())
		r.Recorder.Eventf(policy, "Normal", "PolicyCreated",
			"Created policy %q in DependencyTrack (uuid=%s)",
			policy.Spec.Name, dtPolicy.GetUuid())
	} else {
		// Policy exists with matching UUID. Check for field drift and converge.
		dtExpectedOperator := string(policy.Spec.Operator)
		dtExpectedState := string(policy.Spec.ViolationState)
		needsUpdate := dtPolicy.GetName() != policy.Spec.Name ||
			dtPolicy.GetOperator() != dtExpectedOperator ||
			dtPolicy.GetViolationState() != dtExpectedState
		if needsUpdate {
			updated, _, err := apiClient.PolicyAPI.UpdatePolicy(authCtx).Policy(dtapi.Policy{
				Name:           policy.Spec.Name,
				Uuid:           dtPolicy.GetUuid(),
				Global:         boolPtr(true),
				Operator:       dtExpectedOperator,
				ViolationState: dtExpectedState,
			}).Execute()
			if err != nil {
				r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonPolicyUpdateFailed,
					"Failed to update policy fields in DependencyTrack: "+err.Error())
				return nil, fmt.Errorf("updating policy fields: %w", err)
			}
			dtPolicy = updated
			log.Info("converged policy fields on UUID-owned policy",
				"uuid", dtPolicy.GetUuid())
		}
	}

	return dtPolicy, nil
}

// findPolicyByName searches a list of DT policies for one matching the given name.
// Returns (true, uuid) when found, (false, "") otherwise.
func findPolicyByName(policies []dtapi.Policy, name string) (bool, string) {
	for _, p := range policies {
		if p.GetName() == name {
			return true, p.GetUuid()
		}
	}
	return false, ""
}

// syncPolicyConditions ensures the DT policy has exactly the conditions declared
// in the K8s spec — creating new ones and deleting stale ones.
func (r *PolicyReconciler) syncPolicyConditions(
	ctx context.Context,
	authCtx context.Context,
	apiClient *dtapi.APIClient,
	policy *dependencytrackv1alpha1.Policy,
	dtPolicy *dtapi.Policy,
) error {
	log := logf.FromContext(ctx)

	// Build a set of existing condition UUIDs.
	existingConditionUUIDs := make(map[string]struct{})
	for _, cond := range dtPolicy.GetPolicyConditions() {
		existingConditionUUIDs[cond.GetUuid()] = struct{}{}
	}

	// Track which conditions from the spec are already present by UUID.
	createdUUIDs := make(map[string]struct{})
	for _, specCond := range policy.Spec.Conditions {
		dtCond, found, err := r.findDTCondition(authCtx, apiClient, policy.Status.UUID, specCond)
		if err != nil {
			r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonConditionCheckFailed,
				"Failed to check policy condition: "+err.Error())
			return err
		}
		if found {
			createdUUIDs[dtCond.GetUuid()] = struct{}{}
			continue
		}

		// Condition doesn't exist — create it.
		if err := r.createDTCondition(ctx, authCtx, apiClient, policy, dtPolicy.GetUuid(), specCond); err != nil {
			return err
		}
	}

	// Delete conditions that are no longer in the spec.
	for dtUUID := range existingConditionUUIDs {
		if _, keep := createdUUIDs[dtUUID]; !keep {
			if _, err := apiClient.PolicyConditionAPI.DeletePolicyCondition(authCtx, dtUUID).Execute(); err != nil {
				log.Error(err, "failed to delete policy condition", "uuid", dtUUID)
				// Non-fatal — continue reconciling.
			}
		}
	}

	return nil
}

// findDTCondition looks up whether a spec condition already exists in the DT policy.
// Returns the DT condition if found, or (nil, false, nil) if not found.
func (r *PolicyReconciler) findDTCondition(
	authCtx context.Context,
	apiClient *dtapi.APIClient,
	policyUUID string,
	specCond dependencytrackv1alpha1.PolicyCondition,
) (*dtapi.PolicyCondition, bool, error) {
	// Fetch the full policy to get its conditions.
	dtPolicy, _, err := apiClient.PolicyAPI.GetPolicy(authCtx, policyUUID).Execute()
	if err != nil {
		return nil, false, err
	}

	dtCond := policyConditionToDT(specCond)

	for _, existing := range dtPolicy.GetPolicyConditions() {
		if existing.GetOperator() == dtCond.GetOperator() &&
			existing.GetSubject() == dtCond.GetSubject() &&
			existing.GetValue() == dtCond.GetValue() {
			return &existing, true, nil
		}
	}

	return nil, false, nil
}

// createDTCondition creates a new policy condition in DependencyTrack.
func (r *PolicyReconciler) createDTCondition(
	ctx context.Context,
	authCtx context.Context,
	apiClient *dtapi.APIClient,
	policy *dependencytrackv1alpha1.Policy,
	policyUUID string,
	specCond dependencytrackv1alpha1.PolicyCondition,
) error {
	log := logf.FromContext(ctx)

	dtCond := policyConditionToDT(specCond)
	dtCond.Uuid = uuid.New().String()

	created, _, err := apiClient.PolicyConditionAPI.CreatePolicyCondition(authCtx, policyUUID).
		PolicyCondition(dtCond).Execute()
	if err != nil {
		r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonConditionCreateFailed,
			"Failed to create policy condition in DependencyTrack: "+err.Error())
		return fmt.Errorf("creating policy condition: %w", err)
	}

	log.Info("created policy condition", "uuid", created.GetUuid(), "subject", created.GetSubject())
	return nil
}

// policyConditionToDT converts a K8s PolicyCondition spec to a DT PolicyCondition model.
func policyConditionToDT(specCond dependencytrackv1alpha1.PolicyCondition) dtapi.PolicyCondition {
	return dtapi.PolicyCondition{
		Operator: string(specCond.Operator),
		Subject:  string(specCond.Subject),
		Value:    specCond.Value,
	}
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}

// reconcileDelete handles policy deletion from DependencyTrack. It retains the
// finalizer on any transient failure so that controller-runtime retries, and
// treats a 404 during delete as confirmed absence — removing the finalizer
// without raising an error.
func (r *PolicyReconciler) reconcileDelete(ctx context.Context, policy *dependencytrackv1alpha1.Policy) error {
	log := logf.FromContext(ctx)

	if policy.Status.UUID != "" {
		authCtx, apiClient, err := r.DTProvider.Get(ctx)
		if err != nil {
			log.Error(err, "failed to authenticate with DependencyTrack during deletion")
			r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonCredentialsError,
				"Failed to authenticate with DependencyTrack during deletion: "+err.Error())
			return err
		}

		// Delete associated conditions first (best-effort).
		existing, _, _ := apiClient.PolicyAPI.GetPolicy(authCtx, policy.Status.UUID).Execute()
		if existing != nil {
			for _, cond := range existing.GetPolicyConditions() {
				if _, err := apiClient.PolicyConditionAPI.DeletePolicyCondition(authCtx, cond.GetUuid()).Execute(); err != nil {
					log.Error(err, "failed to delete policy condition", "uuid", cond.GetUuid())
				}
			}
		}

		// Delete the policy itself.
		resp, err := apiClient.PolicyAPI.DeletePolicy(authCtx, policy.Status.UUID).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				// Policy already absent — treat as confirmed absence and
				// proceed to finalizer removal.
				log.Info("policy already absent in DependencyTrack, cleanup confirmed",
					"uuid", policy.Status.UUID)
			} else {
				r.setReadyCondition(ctx, policy, metav1.ConditionFalse, reasonPolicyDeleteFailed,
					"Failed to delete policy from DependencyTrack: "+err.Error())
				return err
			}
		}

		log.Info("deleted policy from DependencyTrack", "uuid", policy.Status.UUID)
		r.Recorder.Eventf(policy, "Normal", "PolicyDeleted",
			"Deleted policy %q from DependencyTrack (uuid=%s)",
			policy.Spec.Name, policy.Status.UUID)
	}

	controllerutil.RemoveFinalizer(policy, policyFinalizer)
	return r.Update(ctx, policy)
}

// setReadyCondition upserts the Ready condition on the policy's status and
// attempts to update the status subresource. It classifies transient API errors
// (HTTP 5xx) as non-blocking so that kube-rbac-proxy retry budgets are not
// exhausted by a single failing Policy resource.
func (r *PolicyReconciler) setReadyCondition(ctx context.Context, policy *dependencytrackv1alpha1.Policy,
	status metav1.ConditionStatus, reason, message string,
) {
	log := logf.FromContext(ctx)
	meta.SetStatusCondition(&policy.Status.Conditions, metav1.Condition{
		Type:               conditionReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: policy.Generation,
	})
	if err := r.Status().Update(ctx, policy); err != nil {
		log.Error(err,
			"failed to update policy status (retryable)",
			"policy", fmt.Sprintf("%s/%s", policy.Namespace, policy.Name))
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.Policy{}).
		Named("policy").
		Complete(r)
}
