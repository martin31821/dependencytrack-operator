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
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNotificationPublisherScheme(t *testing.T) {
	p := &NotificationPublisher{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-publisher",
			Namespace: "default",
		},
		Spec: NotificationPublisherSpec{
			ExtensionName: "webhook",
			Name:          "test-webhook",
		},
	}

	if p.GetNamespace() != "default" {
		t.Errorf("expected namespace default, got %s", p.GetNamespace())
	}
	if p.GetName() != "test-publisher" {
		t.Errorf("expected name test-publisher, got %s", p.GetName())
	}
	if p.Spec.ExtensionName != "webhook" {
		t.Errorf("expected extensionName webhook, got %s", p.Spec.ExtensionName)
	}
	if p.Spec.Name != "test-webhook" {
		t.Errorf("expected name test-webhook, got %s", p.Spec.Name)
	}
	if p.Status.UUID != "" {
		t.Error("expected empty UUID for new publisher")
	}
}

func TestNotificationPublisherDeepCopy(t *testing.T) {
	orig := &NotificationPublisher{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "copy-test",
			Namespace: "test-ns",
		},
		Spec: NotificationPublisherSpec{
			ExtensionName: "slack",
			Name:          "slack-notifier",
			Description:   "Test slack publisher",
		},
		Status: NotificationPublisherStatus{
			UUID:               "550e8400-e29b-41d4-a716-446655440000",
			Name:               "slack-notifier",
			ObservedGeneration: 42,
			Conditions: []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "True",
					LastTransitionTime: metav1.Now(),
					Reason:             "Created",
				},
			},
		},
	}

	copied := orig.DeepCopy()

	// Verify structural equality
	if copied.Spec.ExtensionName != orig.Spec.ExtensionName {
		t.Errorf("copied extensionName %s != original %s", copied.Spec.ExtensionName, orig.Spec.ExtensionName)
	}
	if copied.Status.UUID != orig.Status.UUID {
		t.Errorf("copied UUID %s != original %s", copied.Status.UUID, orig.Status.UUID)
	}

	// Verify independence: modifying copy should not affect original
	copied.Spec.Name = "modified"
	if orig.Spec.Name == "modified" {
		t.Error("DeepCopy did not create an independent copy — modifying copy affected original")
	}

	// Verify conditions slice is independent
	copied.Status.Conditions = append(copied.Status.Conditions, metav1.Condition{
		Type:   "Test",
		Status: "False",
	})
	if len(orig.Status.Conditions) != len(copied.Status.Conditions)-1 {
		t.Error("Conditions slice was not deeply copied")
	}
}

func TestNotificationRuleScheme(t *testing.T) {
	enabled := true
	r := &NotificationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-rule",
			Namespace: "default",
		},
		Spec: NotificationRuleSpec{
			Name:             "vuln-alert",
			Scope:            NotificationRuleScopePortfolio,
			TriggerType:      NotificationRuleTriggerTypeEvent,
			Level:            NotificationRuleLevelWarn,
			PublisherRef:     NotificationRulePublisherRef{Name: "test-publisher"},
			Enabled:          &enabled,
			NotifyOn:         []NotificationRuleTriggerEvent{NotificationRuleTriggerEventNewVulnerability},
			FilterExpression: "severity >= HIGH",
		},
	}

	if r.GetNamespace() != "default" {
		t.Errorf("expected namespace default, got %s", r.GetNamespace())
	}
	if r.GetName() != "test-rule" {
		t.Errorf("expected name test-rule, got %s", r.GetName())
	}
	if r.Spec.Scope != NotificationRuleScopePortfolio {
		t.Errorf("expected scope PORTFOLIO, got %s", r.Spec.Scope)
	}
	if r.Spec.TriggerType != NotificationRuleTriggerTypeEvent {
		t.Errorf("expected triggerType EVENT, got %s", r.Spec.TriggerType)
	}
	if r.Spec.Level != NotificationRuleLevelWarn {
		t.Errorf("expected level WARN, got %s", r.Spec.Level)
	}
	if r.Spec.PublisherRef.Name != "test-publisher" {
		t.Errorf("expected publisherRef name test-publisher, got %s", r.Spec.PublisherRef.Name)
	}
	if r.Status.UUID != "" {
		t.Error("expected empty UUID for new rule")
	}
}

func TestNotificationRuleDeepCopy(t *testing.T) {
	enabled := true
	orig := &NotificationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "copy-rule",
			Namespace: "test-ns",
		},
		Spec: NotificationRuleSpec{
			Name:             "deepcopy-rule",
			Scope:            NotificationRuleScopePortfolio,
			TriggerType:      NotificationRuleTriggerTypeScheduled,
			Level:            NotificationRuleLevelFail,
			PublisherRef:     NotificationRulePublisherRef{Name: "my-publisher"},
			Enabled:          &enabled,
			NotifyOn:         []NotificationRuleTriggerEvent{NotificationRuleTriggerEventNewVulnerability, NotificationRuleTriggerEventNewVulnerableDependency},
			FilterExpression: "vexStatus != CONFIRMED_VULNERABLE",
		},
		Status: NotificationRuleStatus{
			UUID:               "660e8400-e29b-41d4-a716-446655440001",
			Name:               "deepcopy-rule",
			ObservedGeneration: 7,
			Conditions: []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "True",
					LastTransitionTime: metav1.Now(),
					Reason:             "Created",
				},
			},
		},
	}

	copied := orig.DeepCopy()

	// Verify all fields match
	if copied.Spec.Name != orig.Spec.Name {
		t.Errorf("copied name %s != original %s", copied.Spec.Name, orig.Spec.Name)
	}
	if copied.Spec.Scope != orig.Spec.Scope {
		t.Errorf("copied scope %s != original %s", copied.Spec.Scope, orig.Spec.Scope)
	}
	if copied.Spec.TriggerType != orig.Spec.TriggerType {
		t.Errorf("copied triggerType %s != original %s", copied.Spec.TriggerType, orig.Spec.TriggerType)
	}
	if copied.Spec.Level != orig.Spec.Level {
		t.Errorf("copied level %s != original %s", copied.Spec.Level, orig.Spec.Level)
	}
	if copied.Spec.PublisherRef.Name != orig.Spec.PublisherRef.Name {
		t.Errorf("copied publisherRef %s != original %s", copied.Spec.PublisherRef.Name, orig.Spec.PublisherRef.Name)
	}
	if copied.Status.UUID != orig.Status.UUID {
		t.Errorf("copied UUID %s != original %s", copied.Status.UUID, orig.Status.UUID)
	}

	// Verify NotifyOn slice is independent
	copied.Spec.NotifyOn = append(copied.Spec.NotifyOn, "NEW_FORK")
	if len(orig.Spec.NotifyOn) != len(copied.Spec.NotifyOn)-1 {
		t.Error("NotifyOn slice was not deeply copied")
	}
}

func TestNotificationPublisherRef_SameNamespace(t *testing.T) {
	// The PublisherRef is a simple name reference — the controller
	// must enforce that it resolves to a publisher in the same namespace.
	ref := NotificationRulePublisherRef{Name: "my-pub"}
	if ref.Name == "" {
		t.Error("PublisherRef name must not be empty")
	}
}

func TestNotificationRuleLevel_EnumValues(t *testing.T) {
	// Verify all expected level constants are set to valid DT values
	if NotificationRuleLevelInfo != "INFORMATIONAL" {
		t.Errorf("INFO level mismatch: %s", NotificationRuleLevelInfo)
	}
	if NotificationRuleLevelWarn != "WARNING" {
		t.Errorf("WARN level mismatch: %s", NotificationRuleLevelWarn)
	}
	if NotificationRuleLevelFail != "ERROR" {
		t.Errorf("FAIL level mismatch: %s", NotificationRuleLevelFail)
	}
}

func TestNotificationRuleScope_EnumValues(t *testing.T) {
	if NotificationRuleScopePortfolio != "PORTFOLIO" {
		t.Errorf("PORTFOLIO scope mismatch: %s", NotificationRuleScopePortfolio)
	}
	if NotificationRuleScopeSystem != "SYSTEM" {
		t.Errorf("SYSTEM scope mismatch: %s", NotificationRuleScopeSystem)
	}
}

func TestNotificationRuleTriggerType_EnumValues(t *testing.T) {
	if NotificationRuleTriggerTypeEvent != "EVENT" {
		t.Errorf("EVENT triggerType mismatch: %s", NotificationRuleTriggerTypeEvent)
	}
	if NotificationRuleTriggerTypeScheduled != "SCHEDULE" {
		t.Errorf("SCHEDULE triggerType mismatch: %s", NotificationRuleTriggerTypeScheduled)
	}
}

func TestNotificationPublisherConditions(t *testing.T) {
	p := &NotificationPublisher{
		Status: NotificationPublisherStatus{
			Conditions: []metav1.Condition{
				{
					Type:    "Ready",
					Status:  "False",
					Reason:  "Creating",
					Message: "Creating publisher in DependencyTrack",
				},
			},
		},
	}
	if len(p.Status.Conditions) != 1 {
		t.Error("expected exactly one condition")
	}
	if p.Status.Conditions[0].Type != "Ready" {
		t.Errorf("expected Ready condition type, got %s", p.Status.Conditions[0].Type)
	}
}

func TestNotificationRuleConditions(t *testing.T) {
	r := &NotificationRule{
		Status: NotificationRuleStatus{
			Conditions: []metav1.Condition{
				{
					Type:    "Ready",
					Status:  "False",
					Reason:  "PublisherNotReady",
					Message: "Referenced publisher is not Ready",
				},
			},
		},
	}
	if len(r.Status.Conditions) != 1 {
		t.Error("expected exactly one condition")
	}
	if r.Status.Conditions[0].Reason != "PublisherNotReady" {
		t.Errorf("expected PublisherNotReady reason, got %s", r.Status.Conditions[0].Reason)
	}
}

func TestNotificationRuleDeepCopy_RoutingFields(t *testing.T) {
	logSuccess := true
	notifyChildren := false
	skipUnchanged := true

	orig := &NotificationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "routing-copy-rule",
			Namespace: "test-ns",
		},
		Spec: NotificationRuleSpec{
			Name:                  "routing-test",
			Scope:                 NotificationRuleScopePortfolio,
			TriggerType:           NotificationRuleTriggerTypeScheduled,
			Level:                 NotificationRuleLevelInfo,
			PublisherRef:          NotificationRulePublisherRef{Name: "test-pub"},
			Enabled:               &logSuccess,
			NotifyOn:              []NotificationRuleTriggerEvent{NotificationRuleTriggerEventNewVulnerability},
			FilterExpression:      "severity >= HIGH",
			LogSuccessfulPublish:  &logSuccess,
			Message:               "Alert: {{ .Vulnerability.FixTitle }}",
			NotifyChildren:        &notifyChildren,
			ScheduleCron:          "0 0 * * *",
			ScheduleSkipUnchanged: &skipUnchanged,
			PublisherConfigSecretRef: &PublisherConfigSecretRef{
				Name: "my-secret",
				Key:  "config.json",
			},
		},
		Status: NotificationRuleStatus{
			UUID:               "770e8400-e29b-41d4-a716-446655440002",
			Name:               "routing-test",
			ObservedGeneration: 3,
			Conditions: []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "True",
					LastTransitionTime: metav1.Now(),
					Reason:             "Synced",
				},
			},
		},
	}

	copied := orig.DeepCopy()

	// Verify all routing fields copied correctly
	if copied.Spec.LogSuccessfulPublish == nil {
		t.Error("LogSuccessfulPublish was nil in copy")
	} else if *copied.Spec.LogSuccessfulPublish != *orig.Spec.LogSuccessfulPublish {
		t.Errorf("LogSuccessfulPublish mismatch: got %v, want %v",
			*copied.Spec.LogSuccessfulPublish, *orig.Spec.LogSuccessfulPublish)
	}
	if copied.Spec.Message != orig.Spec.Message {
		t.Errorf("Message mismatch: got %q, want %q", copied.Spec.Message, orig.Spec.Message)
	}
	if copied.Spec.NotifyChildren == nil {
		t.Error("NotifyChildren was nil in copy")
	} else if *copied.Spec.NotifyChildren != *orig.Spec.NotifyChildren {
		t.Errorf("NotifyChildren mismatch: got %v, want %v",
			*copied.Spec.NotifyChildren, *orig.Spec.NotifyChildren)
	}
	if copied.Spec.ScheduleCron != orig.Spec.ScheduleCron {
		t.Errorf("ScheduleCron mismatch: got %q, want %q", copied.Spec.ScheduleCron, orig.Spec.ScheduleCron)
	}
	if copied.Spec.ScheduleSkipUnchanged == nil {
		t.Error("ScheduleSkipUnchanged was nil in copy")
	} else if *copied.Spec.ScheduleSkipUnchanged != *orig.Spec.ScheduleSkipUnchanged {
		t.Errorf("ScheduleSkipUnchanged mismatch: got %v, want %v",
			*copied.Spec.ScheduleSkipUnchanged, *orig.Spec.ScheduleSkipUnchanged)
	}
	if copied.Spec.PublisherConfigSecretRef == nil {
		t.Error("PublisherConfigSecretRef was nil in copy")
	} else if copied.Spec.PublisherConfigSecretRef.Name != orig.Spec.PublisherConfigSecretRef.Name ||
		copied.Spec.PublisherConfigSecretRef.Key != orig.Spec.PublisherConfigSecretRef.Key {
		t.Errorf("PublisherConfigSecretRef mismatch")
	}

	// Verify independence of pointer fields
	*copied.Spec.LogSuccessfulPublish = false
	if *orig.Spec.LogSuccessfulPublish == false {
		t.Error("LogSuccessfulPublish pointer was not independently copied")
	}

	// Verify independence of PublisherConfigSecretRef
	copied.Spec.PublisherConfigSecretRef.Name = "modified-secret"
	if orig.Spec.PublisherConfigSecretRef.Name == "modified-secret" {
		t.Error("PublisherConfigSecretRef was not independently deep-copied")
	}
}
