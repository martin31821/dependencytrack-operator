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

// TeamSpec defines the desired state of Team.
type TeamSpec struct {
	// Name of the team to create.
	Name string `json:"name,omitempty"`

	// Permissions is an ordered list of permission names to assign to this team.
	// Valid permission names include: PORTFOLIO_VIEW, PORTFOLIO_MANAGEMENT,
	// VIEW_VULNERABILITY, VULNERABILITY_ANALYSIS, BOM_UPLOAD, PROJECT_CREATION,
	// PROJECT_CREATION_UPLOAD, SYSTEM_CONFIGURATION, ACCESS_MANAGEMENT,
	// VIEW_PORTFOLIO, PROJECT_READ, VULNERABILITY_ASSESSMENT, and others
	// specific to the DependencyTrack version in use.
	// Omit to leave existing permissions unchanged; pass an empty array to
	// clear all permissions.
	// +kubebuilder:validation:Optional
	Permissions []string `json:"permissions,omitempty"`
}

// TeamStatus defines the observed state of Team.
type TeamStatus struct {
	// UUID is the actual UUID of the team in dependencytrack.
	UUID string `json:"uuid,omitempty"`

	// Permissions tracks the permissions last synced to DependencyTrack.
	// This is used for status-only observability; the controller reconciles
	// the actual permission set each reconciliation cycle.
	Permissions string `json:"permissions,omitempty"`

	// Conditions reflect the current reconciliation state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Name records the name that was last synced to DependencyTrack.
	Name string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="UUID",type=string,JSONPath=`.status.uuid`

// Team is the Schema for the teams API.
type Team struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpec   `json:"spec,omitempty"`
	Status TeamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TeamList contains a list of Team.
type TeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Team `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Team{}, &TeamList{})
}
