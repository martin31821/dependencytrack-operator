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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NotificationRuleTriggerEvent defines valid notification event types
// (NotificationGroup values) accepted by Dependency-Track v5.0.2.
// +kubebuilder:validation:Enum=NEW_VULNERABILITIES_SUMMARY;NEW_POLICY_VIOLATIONS_SUMMARY;BOM_VALIDATION_FAILED;USER_DELETED;ANALYZER;BOM_CONSUMED;DATASOURCE_MIRRORING;BOM_PROCESSING_FAILED;BOM_PROCESSED;FILE_SYSTEM;INTEGRATION;CONFIGURATION;VEX_PROCESSED;PROJECT_AUDIT_CHANGE;VULNERABILITY_RETRACTED;VEX_CONSUMED;NEW_VULNERABLE_DEPENDENCY;REPOSITORY;USER_CREATED;POLICY_VIOLATION;NEW_VULNERABILITY;PROJECT_CREATED
type NotificationRuleTriggerEvent string

const (
	NotificationRuleTriggerEventNewVulnerabilitiesSummary  NotificationRuleTriggerEvent = "NEW_VULNERABILITIES_SUMMARY"
	NotificationRuleTriggerEventNewPolicyViolationsSummary NotificationRuleTriggerEvent = "NEW_POLICY_VIOLATIONS_SUMMARY"
	NotificationRuleTriggerEventBomValidationFailed        NotificationRuleTriggerEvent = "BOM_VALIDATION_FAILED"
	NotificationRuleTriggerEventUserDeleted                NotificationRuleTriggerEvent = "USER_DELETED"
	NotificationRuleTriggerEventAnalyzer                   NotificationRuleTriggerEvent = "ANALYZER"
	NotificationRuleTriggerEventBomConsumed                NotificationRuleTriggerEvent = "BOM_CONSUMED"
	NotificationRuleTriggerEventDatasourceMirroring        NotificationRuleTriggerEvent = "DATASOURCE_MIRRORING"
	NotificationRuleTriggerEventBomProcessingFailed        NotificationRuleTriggerEvent = "BOM_PROCESSING_FAILED"
	NotificationRuleTriggerEventBomProcessed               NotificationRuleTriggerEvent = "BOM_PROCESSED"
	NotificationRuleTriggerEventFileSystem                 NotificationRuleTriggerEvent = "FILE_SYSTEM"
	NotificationRuleTriggerEventIntegration                NotificationRuleTriggerEvent = "INTEGRATION"
	NotificationRuleTriggerEventConfiguration              NotificationRuleTriggerEvent = "CONFIGURATION"
	NotificationRuleTriggerEventVexProcessed               NotificationRuleTriggerEvent = "VEX_PROCESSED"
	NotificationRuleTriggerEventProjectAuditChange         NotificationRuleTriggerEvent = "PROJECT_AUDIT_CHANGE"
	NotificationRuleTriggerEventVulnerabilityRetracted     NotificationRuleTriggerEvent = "VULNERABILITY_RETRACTED"
	NotificationRuleTriggerEventVexConsumed                NotificationRuleTriggerEvent = "VEX_CONSUMED"
	NotificationRuleTriggerEventNewVulnerableDependency    NotificationRuleTriggerEvent = "NEW_VULNERABLE_DEPENDENCY"
	NotificationRuleTriggerEventRepository                 NotificationRuleTriggerEvent = "REPOSITORY"
	NotificationRuleTriggerEventUserCreated                NotificationRuleTriggerEvent = "USER_CREATED"
	NotificationRuleTriggerEventPolicyViolation            NotificationRuleTriggerEvent = "POLICY_VIOLATION"
	NotificationRuleTriggerEventNewVulnerability           NotificationRuleTriggerEvent = "NEW_VULNERABILITY"
	NotificationRuleTriggerEventProjectCreated             NotificationRuleTriggerEvent = "PROJECT_CREATED"
)

// NotificationRuleScope defines the scope of a notification rule.
// Matches Dependency-Track v5.0.2 API values: SYSTEM, PORTFOLIO.
// +kubebuilder:validation:Enum=SYSTEM;PORTFOLIO
type NotificationRuleScope string

const (
	NotificationRuleScopeSystem    NotificationRuleScope = "SYSTEM"
	NotificationRuleScopePortfolio NotificationRuleScope = "PORTFOLIO"
)

// NotificationRuleTriggerType defines when a rule fires.
// Matches Dependency-Track v5.0.2 API values: EVENT, SCHEDULE.
// +kubebuilder:validation:Enum=EVENT;SCHEDULE
type NotificationRuleTriggerType string

const (
	NotificationRuleTriggerTypeEvent     NotificationRuleTriggerType = "EVENT"
	NotificationRuleTriggerTypeScheduled NotificationRuleTriggerType = "SCHEDULE"
)

// NotificationRuleLevel defines the notification level.
// Matches Dependency-Track v5.0.2 API values: INFORMATIONAL, WARNING, ERROR.
// +kubebuilder:validation:Enum=INFORMATIONAL;WARNING;ERROR
type NotificationRuleLevel string

const (
	NotificationRuleLevelInfo NotificationRuleLevel = "INFORMATIONAL"
	NotificationRuleLevelWarn NotificationRuleLevel = "WARNING"
	NotificationRuleLevelFail NotificationRuleLevel = "ERROR"
)

// PublisherConfigSecretRef references a Secret key containing publisherConfig JSON.
// The Secret must be in the same namespace as the NotificationRule.
type PublisherConfigSecretRef struct {
	// Name of the Secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Key within the Secret that holds the publisherConfig JSON document.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`
}

// NotificationRulePublisherRef references a managed publisher in the same namespace.
type NotificationRulePublisherRef struct {
	// Name is the name of the NotificationPublisher resource in the same namespace.
	Name string `json:"name"`
}

// NotificationRuleSpec defines the desired state of NotificationRule.
type NotificationRuleSpec struct {
	// Name is the display name for this rule in Dependency-Track.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`

	// Scope controls which projects the rule applies to.
	// +kubebuilder:validation:Required
	Scope NotificationRuleScope `json:"scope"`

	// TriggerType determines whether the rule fires on events or on a schedule.
	// +kubebuilder:validation:Required
	TriggerType NotificationRuleTriggerType `json:"triggerType"`

	// Level filters notifications by severity.
	// +kubebuilder:validation:Required
	Level NotificationRuleLevel `json:"level"`

	// PublisherRef references a managed NotificationPublisher in the same namespace.
	// +kubebuilder:validation:Required
	PublisherRef NotificationRulePublisherRef `json:"publisherRef"`

	// Enabled indicates whether the rule is active.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// NotifyOn lists event types that trigger the rule (e.g. NEW_VULNERABILITY).
	// +kubebuilder:validation:Optional
	NotifyOn []NotificationRuleTriggerEvent `json:"notifyOn,omitempty"`

	// FilterExpression is an optional QL filter string for the rule.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=1024
	FilterExpression string `json:"filterExpression,omitempty"`

	// LogSuccessfulPublish, when true, enables logging for successful
	// notification publishes. When false or omitted, only failures are logged.
	// +kubebuilder:validation:Optional
	LogSuccessfulPublish *bool `json:"logSuccessfulPublish,omitempty"`

	// Message is an optional custom notification message template.
	// Allowed characters: whitespace, Unicode letters, numbers, and punctuation.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=4096
	Message string `json:"message,omitempty"`

	// NotifyChildren controls whether the rule also applies to child projects.
	// Only meaningful when scope is PORTFOLIO or SYSTEM.
	// +kubebuilder:validation:Optional
	NotifyChildren *bool `json:"notifyChildren,omitempty"`

	// ScheduleCron is a cron expression that defines when the scheduled rule
	// fires. Only applicable when triggerType is SCHEDULE. Must not be set
	// for rules with triggerType EVENT.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=256
	ScheduleCron string `json:"scheduleCron,omitempty"`

	// ScheduleSkipUnchanged controls whether the scheduled rule skips emitting
	// a notification if the result is unchanged since the last emission.
	// Only applicable when triggerType is SCHEDULE.
	// +kubebuilder:validation:Optional
	ScheduleSkipUnchanged *bool `json:"scheduleSkipUnchanged,omitempty"`

	// PublisherConfigSecretRef references a Secret whose data key contains
	// the publisherConfig JSON document for this rule.
	// The Secret must reside in the same namespace as the NotificationRule.
	// +kubebuilder:validation:Optional
	PublisherConfigSecretRef *PublisherConfigSecretRef `json:"publisherConfigSecretRef,omitempty"`

	// Teams lists NotificationRuleTeam CRD names in the same namespace
	// whose remote DependencyTrack UUID will be associated with this rule.
	// Only the first matching Team CRD is used; duplicates are deduplicated.
	// +kubebuilder:validation:Optional
	Teams []string `json:"teams,omitempty"`

	// Projects lists project UUIDs to associate with this rule when
	// scope is PORTFOLIO. For PORTFOLIO or SYSTEM scopes,
	// the controller ignores this field.
	// +kubebuilder:validation:Optional
	Projects []string `json:"projects,omitempty"`
}

// NotificationRuleStatus defines the observed state of NotificationRule.
type NotificationRuleStatus struct {
	// UUID is the actual UUID of the rule in DependencyTrack.
	UUID string `json:"uuid,omitempty"`

	// Conditions reflect the current reconciliation state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Name records the name that was last synced to DependencyTrack.
	Name string `json:"name,omitempty"`

	// ObservedGeneration reflects the most recent metadata.Generation
	// observed by the reconciler.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="UUID",type=string,JSONPath=`.status.uuid`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// NotificationRule is the Schema for the notificationrules API.
type NotificationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NotificationRuleSpec   `json:"spec,omitempty"`
	Status NotificationRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NotificationRuleList contains a list of NotificationRule.
type NotificationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NotificationRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NotificationRule{}, &NotificationRuleList{})
}
