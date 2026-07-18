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

// Priority defines the severity level of a Policy.
// +kubebuilder:validation:Enum=CRITICAL;HIGH;MEDIUM;LOW;INFO
type Priority string

const (
	PriorityCritical Priority = "CRITICAL"
	PriorityHigh     Priority = "HIGH"
	PriorityMedium   Priority = "MEDIUM"
	PriorityLow      Priority = "LOW"
	PriorityInfo     Priority = "INFO"
)

// FailureAction defines what happens when a policy condition is violated.
// +kubebuilder:validation:Enum=BLOCK_RELEASE;BLOCK_DEPLOY;REPORT;IGNORE
type FailureAction string

const (
	FailureActionBlockRelease FailureAction = "BLOCK_RELEASE"
	FailureActionBlockDeploy  FailureAction = "BLOCK_DEPLOY"
	FailureActionReport       FailureAction = "REPORT"
	FailureActionIgnore       FailureAction = "IGNORE"
)

// PolicyConditionType defines the type of comparison used in a policy condition.
// +kubebuilder:validation:Enum=CVSS;VULNERABILITY;LICENSE;CPE;PURL;PACKAGE;PACKAGE_TYPE;SEVERITY;CREATED_BEFORE
type PolicyConditionType string

const (
	ConditionTypeCVSS          PolicyConditionType = "CVSS"
	ConditionTypeVulnerability PolicyConditionType = "VULNERABILITY"
	ConditionTypeLicense       PolicyConditionType = "LICENSE"
	ConditionTypeCPE           PolicyConditionType = "CPE"
	ConditionTypePURL          PolicyConditionType = "PURL"
	ConditionTypePackage       PolicyConditionType = "PACKAGE"
	ConditionTypePackageType   PolicyConditionType = "PACKAGE_TYPE"
	ConditionTypeSeverity      PolicyConditionType = "SEVERITY"
	ConditionTypeCreatedBefore PolicyConditionType = "CREATED_BEFORE"
)

// ComparisonOperator defines the comparison operator for a policy condition.
// +kubebuilder:validation:Enum=GT;GTE;LT;LTE;EQ;NE
type ComparisonOperator string

const (
	OpGT  ComparisonOperator = "GT"
	OpGTE ComparisonOperator = "GTE"
	OpLT  ComparisonOperator = "LT"
	OpLTE ComparisonOperator = "LTE"
	OpEQ  ComparisonOperator = "EQ"
	OpNE  ComparisonOperator = "NE"
)

// PolicyCondition defines a single condition within a Policy.
type PolicyCondition struct {
	// Type specifies the kind of component to evaluate (e.g. CVSS, SEVERITY, LICENSE).
	// +kubebuilder:validation:Required
	Type PolicyConditionType `json:"type"`

	// Comparator is the comparison operator applied to the condition value.
	// +kubebuilder:validation:Required
	Comparator ComparisonOperator `json:"comparator"`

	// Value is the value to compare against.
	// The interpretation depends on the Type:
	//   - CVSS/SEVERITY: numeric threshold (e.g. "7.0")
	//   - LICENSE: license identifier string
	//   - CPE/PURL/PACKAGE: pattern to match
	// +kubebuilder:validation:Required
	Value string `json:"value"`

	// IsSuppression when true indicates this condition suppresses (allows)
	// matching findings rather than blocking them.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	IsSuppression bool `json:"isSuppression,omitempty"`
}

// PolicySpec defines the desired state of Policy.
type PolicySpec struct {
	// Name is the human-readable name of the policy.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`

	// Description is an optional human-readable description of the policy.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=1024
	Description string `json:"description,omitempty"`

	// Priority defines the severity level of this policy.
	// +kubebuilder:validation:Required
	Priority Priority `json:"priority"`

	// FailureAction defines what happens when policy conditions are violated.
	// +kubebuilder:validation:Required
	FailureAction FailureAction `json:"failureAction"`

	// Conditions is the list of inline conditions that must be evaluated.
	// At least one condition must be provided.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Conditions []PolicyCondition `json:"conditions"`
}

// PolicyStatus defines the observed state of Policy.
type PolicyStatus struct {
	// UUID is the actual UUID of the policy in DependencyTrack.
	UUID string `json:"uuid,omitempty"`

	// Conditions reflect the current reconciliation state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Name records the name that was last synced to DependencyTrack.
	Name string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="UUID",type=string,JSONPath=`.status.uuid`
// +kubebuilder:printcolumn:name="Priority",type=string,JSONPath=`.spec.priority`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// Policy is the Schema for the policies API.
// Policy resources are namespaced and manage DependencyTrack policies.
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicySpec   `json:"spec,omitempty"`
	Status PolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PolicyList contains a list of Policy.
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Policy{}, &PolicyList{})
}
