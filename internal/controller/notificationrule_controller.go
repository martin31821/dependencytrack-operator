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
	handler "sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/deptrack"
)

const (
	notificationRuleFinalizer = "dependencytrack.mko.dev/notificationrule-finalizer"
)

// NotificationRuleReconciler reconciles a NotificationRule object.
type NotificationRuleReconciler struct {
	client.Client
	Scheme                   *runtime.Scheme
	Recorder                 record.EventRecorder
	DTProvider               deptrack.ClientProviderInterface
	PublisherConfigValidator *PublisherConfigValidator
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationrules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationrules/finalizers,verbs=update
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=notificationpublishers,verbs=get;list;watch

func (r *NotificationRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rule := &dependencytrackv1alpha1.NotificationRule{}
	if err := r.Get(ctx, req.NamespacedName, rule); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure our finalizer is registered before doing any external work.
	if !controllerutil.ContainsFinalizer(rule, notificationRuleFinalizer) {
		controllerutil.AddFinalizer(rule, notificationRuleFinalizer)
		if err := r.Update(ctx, rule); err != nil {
			return ctrl.Result{}, err
		}
		// Requeue; the update will trigger a new reconcile with the finalizer in place.
		return ctrl.Result{}, nil
	}

	if !rule.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, rule)
	}

	return r.reconcileUpsert(ctx, rule)
}

func (r *NotificationRuleReconciler) reconcileUpsert(ctx context.Context, rule *dependencytrackv1alpha1.NotificationRule) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Resolve the publisher reference from the same namespace.
	publisher, err := r.resolvePublisher(ctx, rule)
	if err != nil {
		setRuleCondition(rule, metav1.ConditionFalse, "PublisherRefError", err.Error())
		_ = r.Status().Update(ctx, rule)
		return ctrl.Result{}, err
	}

	if publisher.Status.UUID == "" {
		msg := "publisher is not yet Ready — it has no remote UUID"
		setRuleCondition(rule, metav1.ConditionFalse, "PublisherRefError", msg)
		_ = r.Status().Update(ctx, rule)
		return ctrl.Result{}, nil
	}

	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		setRuleCondition(rule, metav1.ConditionFalse, "CredentialsError", err.Error())
		_ = r.Status().Update(ctx, rule)
		return ctrl.Result{}, err
	}

	var dtUUID string

	if rule.Status.UUID != "" {
		// Rule may already exist — find it by UUID in the list.
		existing, _, err := apiClient.NotificationAPI.GetAllNotificationRules(authCtx).Execute()
		if err != nil {
			return r.failRuleStatus(ctx, rule, "RuleListFailed",
				"failed to list rules from DependencyTrack: "+err.Error(), err)
		}

		found := false
		for _, ruleInList := range existing {
			if ruleInList.Uuid == rule.Status.UUID {
				dtUUID = ruleInList.Uuid
				found = true

				// Check for spec drift and update if needed.
				// Compare required fields and optional fields (enabled, notifyOn,
				// filterExpression) so that the staged-create-then-update path
				// correctly detects drift on a follow-up reconcile.
				if ruleInList.Name != rule.Spec.Name ||
					ruleInList.Scope != string(rule.Spec.Scope) ||
					ruleInList.TriggerType != string(rule.Spec.TriggerType) ||
					driftEnabled(ruleInList.Enabled, rule.Spec.Enabled) ||
					driftNotifyOn(ruleInList.NotifyOn, rule.Spec.NotifyOn) ||
					driftFilterExpression(ruleInList.FilterExpression, rule.Spec.FilterExpression) ||
					driftLogSuccessfulPublish(ruleInList.LogSuccessfulPublish, rule.Spec.LogSuccessfulPublish) ||
					driftNotifyChildren(ruleInList.NotifyChildren, rule.Spec.NotifyChildren) ||
					driftScheduleCron(ruleInList.ScheduleCron, rule.Spec.ScheduleCron) ||
					driftScheduleSkipUnchanged(ruleInList.ScheduleSkipUnchanged, rule.Spec.ScheduleSkipUnchanged) {
					log.Info("updating rule due to spec drift", "uuid", ruleInList.Uuid)

					// Validate and converge publisher config before updating.
					publisherConfig, err := r.validatePublisherConfig(authCtx, rule, apiClient, publisher.Status.UUID)
					if err != nil {
						return r.failRuleStatus(ctx, rule, "PublisherConfigValidationFailed",
							err.Error(), err)
					}

					updateReq := dtapi.UpdateNotificationRuleRequest{
						Uuid:  ruleInList.Uuid,
						Level: string(rule.Spec.Level),
						Name:  rule.Spec.Name,
						Scope: string(rule.Spec.Scope),
					}
					if rule.Spec.Enabled != nil {
						updateReq.Enabled = rule.Spec.Enabled
					}
					if rule.Spec.FilterExpression != "" {
						updateReq.FilterExpression = &rule.Spec.FilterExpression
					}
					if len(rule.Spec.NotifyOn) > 0 {
						updateReq.NotifyOn = notifyOnToStrings(rule.Spec.NotifyOn)
					}
					if publisherConfig != "" {
						updateReq.PublisherConfig = &publisherConfig
					}
					if rule.Spec.LogSuccessfulPublish != nil {
						updateReq.LogSuccessfulPublish = rule.Spec.LogSuccessfulPublish
					}
					if rule.Spec.NotifyChildren != nil {
						updateReq.NotifyChildren = rule.Spec.NotifyChildren
					}
					if rule.Spec.ScheduleCron != "" {
						updateReq.ScheduleCron = &rule.Spec.ScheduleCron
					}
					if rule.Spec.ScheduleSkipUnchanged != nil {
						updateReq.ScheduleSkipUnchanged = rule.Spec.ScheduleSkipUnchanged
					}
					if _, _, err := apiClient.NotificationAPI.UpdateNotificationRule(authCtx).UpdateNotificationRuleRequest(updateReq).Execute(); err != nil {
						// DT v5.0.2 sometimes returns a plain-text 400 response
						// for publisher-configuration validation errors, causing
						// the OpenAPI client to fail JSON parsing with
						// "invalid character 'T'".  Detect this pattern and treat
						// it as a non-blocking config-warning rather than a
						// hard failure — the rule UUID is already set so the
						// rule is considered adopted.
						var hasPublisherConfigIssue bool
						if strings.Contains(err.Error(), "invalid character") &&
							strings.Contains(err.Error(), "publisher requires configuration") {
							hasPublisherConfigIssue = true
						} else if openAPIErr, ok := err.(*dtapi.GenericOpenAPIError); ok {
							if body := openAPIErr.Body(); strings.Contains(string(body), "publisher requires configuration") {
								hasPublisherConfigIssue = true
							}
						}
						if hasPublisherConfigIssue {
							log.V(1).Info("skipping rule update — publisher config required but not provided; rule UUID already set", "uuid", ruleInList.Uuid)
						} else {
							msg := "failed to update rule: " + err.Error()
							// Surface the raw server response body when the error is
							// a GenericOpenAPIError with a body that couldn't be parsed
							// as the expected error model (happens when the API returns
							// a non-RFC-7807 400 response).
							if openAPIErr, ok := err.(*dtapi.GenericOpenAPIError); ok {
								if body := openAPIErr.Body(); len(body) > 0 {
									msg = "failed to update rule: " + err.Error() +
										" — server response: " + string(body)
								}
							}
							return r.failRuleStatus(ctx, rule, "RuleUpdateFailed", msg, err)
						}
					}
				}
				break
			}
		}

		if !found {
			// Rule was deleted out-of-band; fall through to recreate.
			log.Info("rule not found in DependencyTrack, recreating", "uuid", rule.Status.UUID)
			dtUUID = ""
		}
	}

	if dtUUID == "" {
		// Create via the generated API — this is the staged-create path.
		// The create request only needs the four required fields.
		// Optional fields (enabled, notifyOn, filterExpression) are applied
		// in a follow-up update reconcile, proving the UUID is authoritative
		// and resumes correctly after a failure.
		createReq := dtapi.CreateNotificationRuleRequest{
			Level: string(rule.Spec.Level),
			Name:  rule.Spec.Name,
			Publisher: dtapi.Publisher{
				Uuid: publisher.Status.UUID,
			},
			Scope: string(rule.Spec.Scope),
		}
		created, _, err := apiClient.NotificationAPI.CreateNotificationRule(authCtx).CreateNotificationRuleRequest(createReq).Execute()
		if err != nil {
			// Create failed — the name might already exist in DependencyTrack.
			// DT's PUT-to-create handler can throw an unhandled JDO exception
			// (500) instead of a proper 409 Conflict when the name is already
			// taken.  Always attempt adoption first; only surface the original
			// error if adoption also fails.
			if adoptedUUID, adoptErr := r.adoptExistingRule(ctx, apiClient, authCtx, rule); adoptErr == nil {
				dtUUID = adoptedUUID
				log.Info("adopted existing rule in DependencyTrack", "uuid", dtUUID)
			} else {
				return r.failRuleStatus(ctx, rule, "RuleCreateFailed",
					"failed to create rule in DependencyTrack: "+err.Error(), err)
			}
		} else {
			dtUUID = created.Uuid
			log.Info("created rule in DependencyTrack", "uuid", dtUUID)
			r.Recorder.Eventf(rule, "Normal", "RuleCreated", "Created rule %q in DependencyTrack (uuid=%s)", rule.Spec.Name, dtUUID)
		}
	}

	// Converge team and project associations after the rule exists (update or create path).
	if len(rule.Spec.Teams) > 0 {
		teams, err := resolveTeamRefs(ctx, r.Client, apiClient, authCtx, rule.Namespace, rule.Spec.Teams)
		if err != nil {
			return r.failRuleStatus(ctx, rule, "TeamResolutionFailed",
				"failed to resolve team references: "+err.Error(), err)
		}
		if err := convergeTeamAssociations(apiClient, authCtx, dtUUID, teams); err != nil {
			return r.failRuleStatus(ctx, rule, "TeamAssociationFailed",
				"failed to converge team associations: "+err.Error(), err)
		}
	}

	if len(rule.Spec.Projects) > 0 {
		projects := buildProjectAssociations(rule.Spec.Projects)
		if err := convergeProjectAssociations(apiClient, authCtx, dtUUID, projects); err != nil {
			return r.failRuleStatus(ctx, rule, "ProjectAssociationFailed",
				"failed to converge project associations: "+err.Error(), err)
		}
	}

	rule.Status.UUID = dtUUID
	rule.Status.Name = rule.Spec.Name
	rule.Status.ObservedGeneration = rule.Generation

	setRuleCondition(rule, metav1.ConditionTrue, "RuleSynced", "Rule successfully reconciled in DependencyTrack")
	if err := r.Status().Update(ctx, rule); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// conditionReasonRuleDeleteFailed is the Ready reason for credential-safe
// rule-deletion failure conditions.
const conditionReasonRuleDeleteFailed = "RuleDeleteFailed"

func (r *NotificationRuleReconciler) reconcileDelete(ctx context.Context, rule *dependencytrackv1alpha1.NotificationRule) error {
	log := logf.FromContext(ctx)

	if rule.Status.UUID != "" {
		authCtx, apiClient, err := r.DTProvider.Get(ctx)
		if err != nil {
			log.Error(err, "failed to authenticate with DependencyTrack during deletion")
			setRuleCondition(rule, metav1.ConditionFalse, "CredentialsError",
				"failed to authenticate during deletion: "+err.Error())
			_ = r.Status().Update(ctx, rule)
			return err
		}

		// Build a NotificationRule struct for the DELETE body. The generated
		// model serialises required fields (name, scope, triggerType) even when
		// empty, so we must populate them — the server identifies the rule by UUID
		// and ignores the other field values for the delete operation.
		deleteRule := dtapi.NotificationRule{
			Uuid:        rule.Status.UUID,
			Name:        rule.Spec.Name,
			Scope:       string(rule.Spec.Scope),
			TriggerType: string(rule.Spec.TriggerType),
		}
		httpResp, err := apiClient.NotificationAPI.DeleteNotificationRule(authCtx).NotificationRule(deleteRule).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				// Rule is already absent from DependencyTrack — treat as
				// successful cleanup absence.
				log.Info("rule already absent from DependencyTrack (404)", "uuid", rule.Status.UUID)
			} else {
				msg := "failed to delete rule from DependencyTrack: " + err.Error()
				// Surface the raw server response body for 400 errors.
				if openAPIErr, ok := err.(*dtapi.GenericOpenAPIError); ok {
					if body := openAPIErr.Body(); len(body) > 0 {
						msg += " — server response: " + string(body)
					}
				}
				log.Error(err, "failed to delete rule from DependencyTrack", "uuid", rule.Status.UUID)
				setRuleCondition(rule, metav1.ConditionFalse, conditionReasonRuleDeleteFailed,
					msg)
				_ = r.Status().Update(ctx, rule)
				return err
			}
		} else {
			log.Info("deleted rule from DependencyTrack", "uuid", rule.Status.UUID)
			r.Recorder.Eventf(rule, "Normal", "RuleDeleted", "Deleted rule %q from DependencyTrack (uuid=%s)", rule.Spec.Name, rule.Status.UUID)
		}
	}

	controllerutil.RemoveFinalizer(rule, notificationRuleFinalizer)
	return r.Update(ctx, rule)
}

// resolvePublisher looks up the NotificationPublisher referenced by the rule
// in the same namespace.
func (r *NotificationRuleReconciler) resolvePublisher(ctx context.Context, rule *dependencytrackv1alpha1.NotificationRule) (*dependencytrackv1alpha1.NotificationPublisher, error) {
	publisher := &dependencytrackv1alpha1.NotificationPublisher{}
	if err := r.Get(ctx, types.NamespacedName{Name: rule.Spec.PublisherRef.Name, Namespace: rule.Namespace}, publisher); err != nil {
		return nil, err
	}
	return publisher, nil
}

// adoptExistingRule looks for a NotificationRule already existing in
// DependencyTrack that matches the CR's spec.name.  This is a fallback
// used when CreateNotificationRule fails with a 500 (DT's JDO layer
// throws on unique-name violation instead of returning 409).
func (r *NotificationRuleReconciler) adoptExistingRule(ctx context.Context, apiClient *dtapi.APIClient, dtCtx context.Context, rule *dependencytrackv1alpha1.NotificationRule) (string, error) {
	log := logf.FromContext(ctx)

	allRules, _, err := apiClient.NotificationAPI.GetAllNotificationRules(dtCtx).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list rules for adoption: %w", err)
	}

	for _, rt := range allRules {
		if rt.Name == rule.Spec.Name {
			log.Info("adopted existing rule in DependencyTrack", "uuid", rt.Uuid, "name", rt.Name)
			return rt.Uuid, nil
		}
	}

	return "", fmt.Errorf("no existing rule found with name %q", rule.Spec.Name)
}

// failRuleStatus sets a failed condition, persists the status, and returns
// the error so the reconcile loop retries.
func (r *NotificationRuleReconciler) failRuleStatus(ctx context.Context, rule *dependencytrackv1alpha1.NotificationRule, reason, msg string, cause error) (ctrl.Result, error) {
	logf.FromContext(ctx).Error(cause, msg)
	setRuleCondition(rule, metav1.ConditionFalse, reason, msg)
	_ = r.Status().Update(ctx, rule)
	return ctrl.Result{}, cause
}

// setRuleCondition upserts the Ready condition on the rule's status.
func setRuleCondition(rule *dependencytrackv1alpha1.NotificationRule, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&rule.Status.Conditions, metav1.Condition{
		Type:               conditionReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: rule.Generation,
	})
}

// driftEnabled returns true when the spec and remote values disagree.
func driftEnabled(remote *bool, spec *bool) bool {
	if remote == nil && spec == nil {
		return false
	}
	if remote == nil || spec == nil {
		return true
	}
	return *remote != *spec
}

// notifyOnToStrings converts a typed NotificationRuleTriggerEvent slice to []string
// for use with the Dependency-Track API client.
func notifyOnToStrings(events []dependencytrackv1alpha1.NotificationRuleTriggerEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = string(e)
	}
	return out
}

// driftNotifyOn returns true when the spec and remote notifyOn sets disagree.
func driftNotifyOn(remote []string, spec []dependencytrackv1alpha1.NotificationRuleTriggerEvent) bool {
	specStrings := notifyOnToStrings(spec)
	if len(remote) != len(specStrings) {
		return true
	}
	for i := range remote {
		if remote[i] != specStrings[i] {
			return true
		}
	}
	return false
}

// driftFilterExpression returns true when the spec and remote filterExpression disagree.
func driftFilterExpression(remote *string, spec string) bool {
	if remote == nil && spec == "" {
		return false
	}
	if remote == nil || spec == "" {
		return true
	}
	return *remote != spec
}

// validatePublisherConfig reads the Secret referenced by the rule,
// fetches the publisher config schema from DependencyTrack, and validates
// the JSON document against it. Returns the validated JSON string to be
// sent as publisherConfig, or an error with a stable reason.
// Returns ("", nil) when the rule has no SecretRef (no config to converge).
func (r *NotificationRuleReconciler) validatePublisherConfig(ctx context.Context, rule *dependencytrackv1alpha1.NotificationRule, apiClient *dtapi.APIClient, publisherUUID string) (string, error) {
	ref := rule.Spec.PublisherConfigSecretRef
	if ref == nil {
		return "", nil // no reference means no config to converge
	}

	validator := r.PublisherConfigValidator
	if validator == nil {
		return "", &PublisherConfigValidationError{
			Reason:  reasonSchemaRetrieval,
			Message: "publisher config validator not configured",
		}
	}

	// Step 1: fetch the publisher config schema from DT.
	schema, _, err := apiClient.NotificationAPI.GetNotificationPublisherConfigSchema(ctx, publisherUUID).Execute()
	if err != nil {
		return "", &PublisherConfigValidationError{
			Reason:  reasonSchemaRetrieval,
			Message: "failed to retrieve publisher config schema: " + sanitizeJSONError(err),
		}
	}
	if schema == nil {
		return "", &PublisherConfigValidationError{
			Reason:  reasonNoSchemaReturned,
			Message: "publisher config schema is empty",
		}
	}

	// Step 2: read the Secret data.
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: rule.Namespace}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return "", &PublisherConfigValidationError{
				Reason:  reasonSecretNotFound,
				Message: fmt.Sprintf("secret %q not found", ref.Name),
			}
		}
		return "", &PublisherConfigValidationError{
			Reason:  reasonSecretNotFound,
			Message: fmt.Sprintf("failed to read secret %q: %s", ref.Name, sanitizeJSONError(err)),
		}
	}

	data, ok := secret.Data[ref.Key]
	if !ok {
		return "", &PublisherConfigValidationError{
			Reason:  reasonSecretKeyNotFound,
			Message: fmt.Sprintf("key %q not found in secret %q", ref.Key, ref.Name),
		}
	}

	// Step 3: validate against schema.
	if err := validator.ValidateSecretData(data, schema); err != nil {
		return "", err
	}

	return string(data), nil
}

// indexFieldNotificationRuleSecretRef is the field index name used to map
// NotificationRules to the Secret they reference via publisherConfigSecretRef.
const indexFieldNotificationRuleSecretRef = "spec.publisherConfigSecretRef.name"

// indexFuncNotificationRuleBySecretRef is the indexer function used to create
// a namespace-scoped index from Secret-name → NotificationRule namespaced names.
func indexFuncNotificationRuleBySecretRef(obj client.Object) []string {
	rule, ok := obj.(*dependencytrackv1alpha1.NotificationRule)
	if !ok {
		return nil
	}
	ref := rule.Spec.PublisherConfigSecretRef
	if ref == nil {
		return nil
	}
	return []string{ref.Name}
}

// enqueueRequestsForSecret is the mapper function used for the Secret watch.
// It looks up NotificationRules in the same namespace that reference the
// Secret's name (via the index created by indexFuncNotificationRuleBySecretRef).
func enqueueRequestsForSecret(ctx context.Context, c client.Client, obj client.Object) []ctrl.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil
	}
	var rules dependencytrackv1alpha1.NotificationRuleList
	if err := c.List(ctx, &rules, client.MatchingFields{
		indexFieldNotificationRuleSecretRef: secret.Name,
	}, client.InNamespace(secret.Namespace)); err != nil {
		return nil
	}
	reqs := make([]ctrl.Request, 0, len(rules.Items))
	for _, rule := range rules.Items {
		reqs = append(reqs, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      rule.Name,
				Namespace: rule.Namespace,
			},
		})
	}
	return reqs
}

// SetupWithManager sets up the controller with the Manager.
func (r *NotificationRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&dependencytrackv1alpha1.NotificationRule{},
		indexFieldNotificationRuleSecretRef,
		indexFuncNotificationRuleBySecretRef,
	); err != nil {
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&dependencytrackv1alpha1.NotificationRule{},
		indexFieldNotificationRuleTeamName,
		indexFuncNotificationRuleByTeamName,
	); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.NotificationRule{}).
		Owns(&dependencytrackv1alpha1.NotificationPublisher{}).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(
				func(ctx context.Context, obj client.Object) []ctrl.Request {
					return enqueueRequestsForSecret(ctx, mgr.GetClient(), obj)
				},
			),
		).
		Watches(
			&dependencytrackv1alpha1.Team{},
			handler.EnqueueRequestsFromMapFunc(
				func(ctx context.Context, obj client.Object) []ctrl.Request {
					return enqueueRequestsForTeam(ctx, mgr.GetClient(), obj)
				},
			),
		).
		Named("notificationrule").
		Complete(r)
}
