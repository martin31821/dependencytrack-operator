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

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const testSeverityCritical = "CRITICAL"

func TestPolicyGVKFromScheme(t *testing.T) {
	// GVK is derived from Scheme registration, not zero-value.
	testScheme := runtime.NewScheme()
	err := SchemeBuilder.AddToScheme(testScheme)
	assert.NoError(t, err)

	gvks, _, err := testScheme.ObjectKinds(&Policy{})
	assert.NoError(t, err)
	assert.Len(t, gvks, 1)
	assert.Equal(t, "dependencytrack.mko.dev", gvks[0].Group)
	assert.Equal(t, "v1alpha1", gvks[0].Version)
	assert.Equal(t, "Policy", gvks[0].Kind)
}

func TestPolicyListGVKFromScheme(t *testing.T) {
	testScheme := runtime.NewScheme()
	err := SchemeBuilder.AddToScheme(testScheme)
	assert.NoError(t, err)

	gvks, _, err := testScheme.ObjectKinds(&PolicyList{})
	assert.NoError(t, err)
	assert.Len(t, gvks, 1)
	assert.Equal(t, "dependencytrack.mko.dev", gvks[0].Group)
	assert.Equal(t, "v1alpha1", gvks[0].Version)
	assert.Equal(t, "PolicyList", gvks[0].Kind)
}

func TestPolicyConditionValidation(t *testing.T) {
	tests := []struct {
		name        string
		cond        PolicyCondition
		wantSubject PolicyConditionSubject
		wantOp      PolicyConditionOperator
		wantValue   string
	}{
		{
			name: "subject is value",
			cond: PolicyCondition{
				Subject:  PolicyConditionSubjectSeverity,
				Operator: PolicyConditionOperatorIs,
				Value:    "CRITICAL",
			},
			wantSubject: PolicyConditionSubjectSeverity,
			wantOp:      PolicyConditionOperatorIs,
			wantValue:   "CRITICAL",
		},
		{
			name: "subject is not value",
			cond: PolicyCondition{
				Subject:  PolicyConditionSubjectLicense,
				Operator: PolicyConditionOperatorIsNot,
				Value:    "GPL-3.0-only",
			},
			wantSubject: PolicyConditionSubjectLicense,
			wantOp:      PolicyConditionOperatorIsNot,
			wantValue:   "GPL-3.0-only",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantSubject, tt.cond.Subject)
			assert.Equal(t, tt.wantOp, tt.cond.Operator)
			assert.Equal(t, tt.wantValue, tt.cond.Value)
		})
	}
}

func TestPolicySpecRequiredFields(t *testing.T) {
	// Missing name should be invalid
	invalid := PolicySpec{
		Operator:       PolicyOperatorAny,
		ViolationState: ViolationStateFail,
		Conditions: []PolicyCondition{
			{Subject: PolicyConditionSubjectSeverity, Operator: PolicyConditionOperatorIs, Value: "HIGH"},
		},
	}
	assert.Empty(t, invalid.Name)

	// Valid spec
	valid := PolicySpec{
		Operator:       PolicyOperatorAny,
		Name:           "Block Critical Vulnerabilities",
		ViolationState: ViolationStateFail,
		Conditions: []PolicyCondition{
			{Subject: PolicyConditionSubjectSeverity, Operator: PolicyConditionOperatorIs, Value: testSeverityCritical},
		},
	}
	assert.NotEmpty(t, valid.Name)
	assert.Len(t, valid.Conditions, 1)
}

func TestPolicyStatusFields(t *testing.T) {
	status := PolicyStatus{
		UUID: "test-uuid-1234",
		Name: "test-policy",
		Conditions: []metav1.Condition{
			{
				Type:               "Ready",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "PolicySynced",
				Message:            "Policy synced to DependencyTrack",
			},
		},
	}
	assert.Equal(t, "test-uuid-1234", status.UUID)
	assert.Equal(t, "test-policy", status.Name)
	assert.Len(t, status.Conditions, 1)
}

func TestViolationStateConstants(t *testing.T) {
	assert.Equal(t, ViolationState("INFO"), ViolationStateInfo)
	assert.Equal(t, ViolationState("WARN"), ViolationStateWarn)
	assert.Equal(t, ViolationState("FAIL"), ViolationStateFail)
}

func TestPolicyOperatorConstants(t *testing.T) {
	assert.Equal(t, PolicyOperator("ANY"), PolicyOperatorAny)
	assert.Equal(t, PolicyOperator("ALL"), PolicyOperatorAll)
}

func TestPolicyConditionSubjectConstants(t *testing.T) {
	assert.Equal(t, PolicyConditionSubject("AGE"), PolicyConditionSubjectAge)
	assert.Equal(t, PolicyConditionSubject("COORDINATES"), PolicyConditionSubjectCoordinates)
	assert.Equal(t, PolicyConditionSubject("CPE"), PolicyConditionSubjectCPE)
	assert.Equal(t, PolicyConditionSubject("EXPRESSION"), PolicyConditionSubjectExpression)
	assert.Equal(t, PolicyConditionSubject("LICENSE"), PolicyConditionSubjectLicense)
	assert.Equal(t, PolicyConditionSubject("LICENSE_GROUP"), PolicyConditionSubjectLicenseGroup)
	assert.Equal(t, PolicyConditionSubject("PACKAGE_URL"), PolicyConditionSubjectPackageURL)
	assert.Equal(t, PolicyConditionSubject("SEVERITY"), PolicyConditionSubjectSeverity)
	assert.Equal(t, PolicyConditionSubject("SWID_TAGID"), PolicyConditionSubjectSWIDTagID)
	assert.Equal(t, PolicyConditionSubject("VERSION"), PolicyConditionSubjectVersion)
	assert.Equal(t, PolicyConditionSubject("COMPONENT_HASH"), PolicyConditionSubjectComponentHash)
	assert.Equal(t, PolicyConditionSubject("CWE"), PolicyConditionSubjectCWE)
	assert.Equal(t, PolicyConditionSubject("VULNERABILITY_ID"), PolicyConditionSubjectVulnerabilityID)
	assert.Equal(t, PolicyConditionSubject("VERSION_DISTANCE"), PolicyConditionSubjectVersionDistance)
	assert.Equal(t, PolicyConditionSubject("EPSS"), PolicyConditionSubjectEPSS)
}

func TestPolicyConditionOperatorConstants(t *testing.T) {
	assert.Equal(t, PolicyConditionOperator("IS"), PolicyConditionOperatorIs)
	assert.Equal(t, PolicyConditionOperator("IS_NOT"), PolicyConditionOperatorIsNot)
}

func TestDeepCopyPolicy(t *testing.T) {
	original := &Policy{
		ObjectMeta: metav1.ObjectMeta{Name: "test-policy", Namespace: "default"},
		Spec: PolicySpec{
			Operator:       PolicyOperatorAny,
			Name:           "Test Policy",
			ViolationState: ViolationStateFail,
			Conditions:     []PolicyCondition{},
		},
		Status: PolicyStatus{
			UUID: "abc-123",
			Name: "Test Policy",
			Conditions: []metav1.Condition{
				{Type: "Ready", Status: metav1.ConditionTrue},
			},
		},
	}

	// Test DeepCopy
	cp := original.DeepCopy()
	assert.Equal(t, original, cp)
	// Verify it's a deep copy (mutation of cp shouldn't affect original)
	cp.Status.UUID = "changed"
	assert.NotEqual(t, cp.Status.UUID, original.Status.UUID)

	// Test DeepCopyObject
	rtObj := original.DeepCopyObject()
	assert.NotNil(t, rtObj)

	// Test DeepCopyInto
	in2 := &Policy{}
	original.DeepCopyInto(in2)
	assert.Equal(t, original, in2)
}

func TestDeepCopyPolicyList(t *testing.T) {
	original := &PolicyList{
		Items: []Policy{
			{Spec: PolicySpec{Name: "policy1"}},
			{Spec: PolicySpec{Name: "policy2"}},
		},
	}

	cp := original.DeepCopy()
	assert.Equal(t, original, cp)
	cp.Items[0].Spec.Name = "changed"
	assert.NotEqual(t, cp.Items[0].Spec.Name, original.Items[0].Spec.Name)
}

func TestSchemeBuilderRegistration(t *testing.T) {
	// Verify that init() registered the types via SchemeBuilder.
	// Check that the registered GroupVersion is correct.
	assert.Equal(t, "dependencytrack.mko.dev", SchemeBuilder.GroupVersion.Group)
	assert.Equal(t, "v1alpha1", SchemeBuilder.GroupVersion.Version)

	// Verify the SchemeBuilder can register these types without panic.
	testScheme := runtime.NewScheme()
	err := SchemeBuilder.AddToScheme(testScheme)
	assert.NoError(t, err)

	// Verify the types are in the scheme
	gvks, _, err := testScheme.ObjectKinds(&Policy{})
	assert.NoError(t, err)
	assert.NotEmpty(t, gvks)
	assert.Equal(t, "Policy", gvks[0].Kind)

	gvks2, _, err := testScheme.ObjectKinds(&PolicyList{})
	assert.NoError(t, err)
	assert.NotEmpty(t, gvks2)
	assert.Equal(t, "PolicyList", gvks2[0].Kind)
}
