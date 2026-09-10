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

// ViolationState defines the severity assigned when a policy is violated.
// Values match the Dependency-Track Policy API and UI directly.
// +kubebuilder:validation:Enum=INFO;WARN;FAIL
type ViolationState string

const (
	ViolationStateInfo ViolationState = "INFO"
	ViolationStateWarn ViolationState = "WARN"
	ViolationStateFail ViolationState = "FAIL"
)

// PolicyOperator defines how a policy combines its conditions.
// Values match the Dependency-Track Policy API and UI directly.
// +kubebuilder:validation:Enum=ANY;ALL
type PolicyOperator string

const (
	PolicyOperatorAny PolicyOperator = "ANY"
	PolicyOperatorAll PolicyOperator = "ALL"
)

// ViolationType defines the class of violation recorded when a condition matches.
// Values match the Dependency-Track Policy API and UI directly.
// +kubebuilder:validation:Enum=LICENSE;OPERATIONAL;SECURITY
type ViolationType string

const (
	ViolationTypeLicense     ViolationType = "LICENSE"
	ViolationTypeOperational ViolationType = "OPERATIONAL"
	ViolationTypeSecurity    ViolationType = "SECURITY"
)

// PolicyConditionSubject defines the Dependency-Track subject evaluated by a condition.
// +kubebuilder:validation:Enum=AGE;COORDINATES;CPE;EXPRESSION;LICENSE;LICENSE_GROUP;PACKAGE_URL;SEVERITY;SWID_TAGID;VERSION;COMPONENT_HASH;CWE;VULNERABILITY_ID;VERSION_DISTANCE;EPSS
type PolicyConditionSubject string

const (
	PolicyConditionSubjectAge             PolicyConditionSubject = "AGE"
	PolicyConditionSubjectCoordinates     PolicyConditionSubject = "COORDINATES"
	PolicyConditionSubjectCPE             PolicyConditionSubject = "CPE"
	PolicyConditionSubjectExpression      PolicyConditionSubject = "EXPRESSION"
	PolicyConditionSubjectLicense         PolicyConditionSubject = "LICENSE"
	PolicyConditionSubjectLicenseGroup    PolicyConditionSubject = "LICENSE_GROUP"
	PolicyConditionSubjectPackageURL      PolicyConditionSubject = "PACKAGE_URL"
	PolicyConditionSubjectSeverity        PolicyConditionSubject = "SEVERITY"
	PolicyConditionSubjectSWIDTagID       PolicyConditionSubject = "SWID_TAGID"
	PolicyConditionSubjectVersion         PolicyConditionSubject = "VERSION"
	PolicyConditionSubjectComponentHash   PolicyConditionSubject = "COMPONENT_HASH"
	PolicyConditionSubjectCWE             PolicyConditionSubject = "CWE"
	PolicyConditionSubjectVulnerabilityID PolicyConditionSubject = "VULNERABILITY_ID"
	PolicyConditionSubjectVersionDistance PolicyConditionSubject = "VERSION_DISTANCE"
	PolicyConditionSubjectEPSS            PolicyConditionSubject = "EPSS"
)

// PolicyConditionOperator defines how a condition compares its subject to its value.
// Values match the Dependency-Track PolicyCondition API.
// +kubebuilder:validation:Enum=IS;IS_NOT;MATCHES;NO_MATCH;NUMERIC_GREATER_THAN;NUMERIC_LESS_THAN;NUMERIC_EQUAL;NUMERIC_NOT_EQUAL;NUMERIC_GREATER_THAN_OR_EQUAL;NUMERIC_LESSER_THAN_OR_EQUAL;CONTAINS_ALL;CONTAINS_ANY
type PolicyConditionOperator string

const (
	PolicyConditionOperatorIs                        PolicyConditionOperator = "IS"
	PolicyConditionOperatorIsNot                     PolicyConditionOperator = "IS_NOT"
	PolicyConditionOperatorMatches                   PolicyConditionOperator = "MATCHES"
	PolicyConditionOperatorNoMatch                   PolicyConditionOperator = "NO_MATCH"
	PolicyConditionOperatorNumericGreaterThan        PolicyConditionOperator = "NUMERIC_GREATER_THAN"
	PolicyConditionOperatorNumericLessThan           PolicyConditionOperator = "NUMERIC_LESS_THAN"
	PolicyConditionOperatorNumericEqual              PolicyConditionOperator = "NUMERIC_EQUAL"
	PolicyConditionOperatorNumericNotEqual           PolicyConditionOperator = "NUMERIC_NOT_EQUAL"
	PolicyConditionOperatorNumericGreaterThanOrEqual PolicyConditionOperator = "NUMERIC_GREATER_THAN_OR_EQUAL"
	PolicyConditionOperatorNumericLesserThanOrEqual  PolicyConditionOperator = "NUMERIC_LESSER_THAN_OR_EQUAL"
	PolicyConditionOperatorContainsAll               PolicyConditionOperator = "CONTAINS_ALL"
	PolicyConditionOperatorContainsAny               PolicyConditionOperator = "CONTAINS_ANY"
)

// PolicyCondition defines a single condition within a Policy.
//
// The EXPRESSION subject is special: its Value is a CEL expression evaluated
// for every component in scope, its Operator is not used (Dependency-Track
// always records the condition as MATCHES), and its ViolationType must be
// chosen explicitly because the subject does not imply one.
// +kubebuilder:validation:XValidation:rule="self.subject == 'EXPRESSION' || has(self.operator)",message="operator is required when subject is not EXPRESSION"
// +kubebuilder:validation:XValidation:rule="self.subject != 'EXPRESSION' || has(self.violationType)",message="violationType is required when subject is EXPRESSION"
type PolicyCondition struct {
	// Subject specifies the Dependency-Track property to evaluate.
	// +kubebuilder:validation:Required
	Subject PolicyConditionSubject `json:"subject"`

	// Operator specifies whether the subject is or is not the configured value.
	// It is required for every subject except EXPRESSION, where Dependency-Track
	// always evaluates the expression as a MATCHES comparison.
	// +optional
	Operator PolicyConditionOperator `json:"operator,omitempty"`

	// Value is the value compared against the subject.
	// For the EXPRESSION subject, Value is a CEL expression that is evaluated
	// for every component in scope; the condition matches when it is true.
	// +kubebuilder:validation:Required
	Value string `json:"value"`

	// ViolationType is the class of violation recorded when the condition
	// matches. It is required for the EXPRESSION subject, which does not imply
	// a violation type. For every other subject the violation type is derived
	// from the subject and this field is ignored.
	// +optional
	ViolationType *ViolationType `json:"violationType,omitempty"`
}

// PolicySpec defines the desired state of Policy.
type PolicySpec struct {
	// Name is the human-readable name of the policy.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`

	// Operator defines whether any or all policy conditions must match.
	// +kubebuilder:validation:Required
	Operator PolicyOperator `json:"operator"`

	// ViolationState defines the severity assigned when policy conditions are violated.
	// +kubebuilder:validation:Required
	ViolationState ViolationState `json:"violationState"`

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
