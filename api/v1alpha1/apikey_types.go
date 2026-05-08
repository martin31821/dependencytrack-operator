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

// APIKeySpec defines the desired state of APIKey.
type APIKeySpec struct {
	// TeamRef is the name of the Team CR (in the same namespace) to generate the key for.
	// +kubebuilder:validation:Required
	TeamRef string `json:"teamRef"`

	// SecretName is the name of the Kubernetes Secret where the generated API key will be stored.
	// The Secret is created (or updated) in the same namespace as the APIKey resource.
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`

	// Comment is an optional human-readable label attached to the API key in DependencyTrack.
	// +kubebuilder:validation:Optional
	Comment string `json:"comment,omitempty"`
}

// APIKeyStatus defines the observed state of APIKey.
type APIKeyStatus struct {
	// PublicID is DependencyTrack's stable identifier for the API key.
	// It is used for all subsequent update and delete operations.
	PublicID string `json:"publicId,omitempty"`

	// Conditions reflect the current reconciliation state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Team",type=string,JSONPath=`.spec.teamRef`
// +kubebuilder:printcolumn:name="Secret",type=string,JSONPath=`.spec.secretName`
// +kubebuilder:printcolumn:name="PublicID",type=string,JSONPath=`.status.publicId`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Reconciled")].status`

// APIKey is the Schema for the apikeys API.
type APIKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIKeySpec   `json:"spec,omitempty"`
	Status APIKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// APIKeyList contains a list of APIKey.
type APIKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIKey{}, &APIKeyList{})
}
