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

// NotificationPublisherSpec defines the desired state of NotificationPublisher.
type NotificationPublisherSpec struct {
	// ExtensionName identifies the Dependency-Track publisher extension.
	// Valid values include: email, slack, webhook, opsgenie, and others
	// supported by the installed DT plugins.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	ExtensionName string `json:"extensionName"`

	// Name is the display name for this publisher in Dependency-Track.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`

	// Description is an optional human-readable description of the publisher.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=1024
	Description string `json:"description,omitempty"`

	// Template is an optional custom notification message template that the
	// publisher applies when delivering notifications. When omitted,
	// Dependency-Track applies the publisher extension's default template.
	// This is the supported mechanism for supplying a bespoke notification
	// body; unlike the legacy per-rule "message" field, it is unconstrained
	// in length and survives beyond Dependency-Track v5.0.2.
	// +kubebuilder:validation:Optional
	Template *string `json:"template,omitempty"`

	// TemplateMimeType declares the media type of the Template body
	// (e.g. "text/plain", "text/html", "application/json"). When omitted,
	// Dependency-Track applies the publisher extension's default mime type.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	TemplateMimeType *string `json:"templateMimeType,omitempty"`
}

// NotificationPublisherStatus defines the observed state of NotificationPublisher.
type NotificationPublisherStatus struct {
	// UUID is the actual UUID of the publisher in DependencyTrack.
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

// NotificationPublisher is the Schema for the notificationpublishers API.
type NotificationPublisher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NotificationPublisherSpec   `json:"spec,omitempty"`
	Status NotificationPublisherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NotificationPublisherList contains a list of NotificationPublisher.
type NotificationPublisherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NotificationPublisher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NotificationPublisher{}, &NotificationPublisherList{})
}
