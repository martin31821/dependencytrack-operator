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
	// Valid permission names include: PORTFOLIO_MANAGEMENT,
	// VIEW_VULNERABILITY, VULNERABILITY_ANALYSIS, BOM_UPLOAD, PROJECT_CREATION,
	// PROJECT_CREATION_UPLOAD, SYSTEM_CONFIGURATION, ACCESS_MANAGEMENT,
	// VIEW_PORTFOLIO, PROJECT_READ, VULNERABILITY_ASSESSMENT, and others
	// specific to the DependencyTrack version in use.
	// Omit to leave existing permissions unchanged; pass an empty array to
	// clear all permissions.
	// +kubebuilder:validation:Optional
	Permissions []string `json:"permissions,omitempty"`

	// OIDC configures OpenID Connect group mapping for this team. It is a
	// pointer so that a nil value is distinguishable from an empty Groups
	// slice; the controller treats nil as "OIDC management disabled" and
	// therefore skips all OIDC API calls (see T04).
	// Minimum-cardinality / membership validation is deliberately deferred to
	// the admission webhook (see T02) rather than encoded as a CRD enum, so
	// that arbitrary upstream group-name casing is accepted verbatim.
	// +kubebuilder:validation:Optional
	OIDC *TeamOIDCConfig `json:"oidc,omitempty"`
}

// TeamOIDCConfig holds the desired OIDC group-mapping configuration for a Team.
type TeamOIDCConfig struct {
	// Groups is the list of OIDC group names to map to this team. Entries
	// are compared verbatim against upstream Identity Provider claims: casing
	// is intentionally preserved (no normalization happens in the API types)
	// so operators retain full control over matching semantics.
	// +kubebuilder:validation:Optional
	Groups []string `json:"groups,omitempty"`
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

	// OIDC reflects the observed OIDC-to-team mapping state. A nil pointer
	// means the controller has not yet reconciled OIDC for this team;
	// an empty OwnedMappings slice confirms that no mappings exist.
	// +kubebuilder:validation:Optional
	OIDC *TeamOIDCStatus `json:"oidc,omitempty"`
}

// OwnedOIDCMapping describes a single OIDC-to-Team mapping owned by a Team.
type OwnedOIDCMapping struct {
	// GroupName is the OIDC group name as delivered by the Identity Provider.
	// Stored verbatim, case-unaltered.
	GroupName string `json:"groupName,omitempty"`

	// GroupUUID is the DependencyTrack UUID of the OIDC group.
	GroupUUID string `json:"groupUuid,omitempty"`

	// TeamUUID is the DependencyTrack UUID of the owning team.
	TeamUUID string `json:"teamUuid,omitempty"`

	// MappingUUID is the DependencyTrack UUID of the OIDC group->team mapping.
	MappingUUID string `json:"mappingUuid,omitempty"`
}

// TeamOIDCStatus reflects the observed OIDC mapping state for a Team.
type TeamOIDCStatus struct {
	// OwnedMappings enumerates the OIDC-to-team mappings currently owned by
	// this team in DependencyTrack.
	// +kubebuilder:validation:Optional
	OwnedMappings []OwnedOIDCMapping `json:"ownedMappings,omitempty"`
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
