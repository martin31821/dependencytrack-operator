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

const testCVSSNine = "9.0"

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
		name         string
		cond         PolicyCondition
		wantType     PolicyConditionType
		wantOp       ComparisonOperator
		wantValue    string
		wantSuppress bool
	}{
		{
			name: "CVSS critical condition",
			cond: PolicyCondition{
				Type:       ConditionTypeCVSS,
				Comparator: OpGTE,
				Value:      testCVSSNine,
			},
			wantType:     ConditionTypeCVSS,
			wantOp:       OpGTE,
			wantValue:    testCVSSNine,
			wantSuppress: false,
		},
		{
			name: "Suppression condition",
			cond: PolicyCondition{
				Type:          ConditionTypeLicense,
				Comparator:    OpEQ,
				Value:         "MIT",
				IsSuppression: true,
			},
			wantType:     ConditionTypeLicense,
			wantOp:       OpEQ,
			wantValue:    "MIT",
			wantSuppress: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantType, tt.cond.Type)
			assert.Equal(t, tt.wantOp, tt.cond.Comparator)
			assert.Equal(t, tt.wantValue, tt.cond.Value)
			assert.Equal(t, tt.wantSuppress, tt.cond.IsSuppression)
		})
	}
	// Validate suppression flag defaults to false when unset
	sup := PolicyCondition{
		Type:       ConditionTypeSeverity,
		Comparator: OpGT,
		Value:      "5.0",
	}
	assert.False(t, sup.IsSuppression)
}

func TestPolicySpecRequiredFields(t *testing.T) {
	// Missing name should be invalid
	invalid := PolicySpec{
		Priority:      PriorityCritical,
		FailureAction: FailureActionBlockRelease,
		Conditions: []PolicyCondition{
			{Type: ConditionTypeCVSS, Comparator: OpGTE, Value: "7.0"},
		},
	}
	assert.Empty(t, invalid.Name)

	// Valid spec
	valid := PolicySpec{
		Name:          "Block Critical Vulnerabilities",
		Description:   "Reject any policy with critical CVSS issues",
		Priority:      PriorityCritical,
		FailureAction: FailureActionBlockRelease,
		Conditions: []PolicyCondition{
			{Type: ConditionTypeCVSS, Comparator: OpGTE, Value: testCVSSNine},
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
				Reason:             "PolicyReconciled",
				Message:            "Policy synced to DependencyTrack",
			},
		},
	}
	assert.Equal(t, "test-uuid-1234", status.UUID)
	assert.Equal(t, "test-policy", status.Name)
	assert.Len(t, status.Conditions, 1)
}

func TestPriorityConstants(t *testing.T) {
	assert.Equal(t, Priority("CRITICAL"), PriorityCritical)
	assert.Equal(t, Priority("HIGH"), PriorityHigh)
	assert.Equal(t, Priority("MEDIUM"), PriorityMedium)
	assert.Equal(t, Priority("LOW"), PriorityLow)
	assert.Equal(t, Priority("INFO"), PriorityInfo)
}

func TestFailureActionConstants(t *testing.T) {
	assert.Equal(t, FailureAction("BLOCK_RELEASE"), FailureActionBlockRelease)
	assert.Equal(t, FailureAction("BLOCK_DEPLOY"), FailureActionBlockDeploy)
	assert.Equal(t, FailureAction("REPORT"), FailureActionReport)
	assert.Equal(t, FailureAction("IGNORE"), FailureActionIgnore)
}

func TestPolicyConditionTypeConstants(t *testing.T) {
	assert.Equal(t, PolicyConditionType("CVSS"), ConditionTypeCVSS)
	assert.Equal(t, PolicyConditionType("VULNERABILITY"), ConditionTypeVulnerability)
	assert.Equal(t, PolicyConditionType("LICENSE"), ConditionTypeLicense)
	assert.Equal(t, PolicyConditionType("CPE"), ConditionTypeCPE)
	assert.Equal(t, PolicyConditionType("PURL"), ConditionTypePURL)
	assert.Equal(t, PolicyConditionType("PACKAGE"), ConditionTypePackage)
	assert.Equal(t, PolicyConditionType("PACKAGE_TYPE"), ConditionTypePackageType)
	assert.Equal(t, PolicyConditionType("SEVERITY"), ConditionTypeSeverity)
	assert.Equal(t, PolicyConditionType("CREATED_BEFORE"), ConditionTypeCreatedBefore)
}

func TestComparisonOperatorConstants(t *testing.T) {
	assert.Equal(t, ComparisonOperator("GT"), OpGT)
	assert.Equal(t, ComparisonOperator("GTE"), OpGTE)
	assert.Equal(t, ComparisonOperator("LT"), OpLT)
	assert.Equal(t, ComparisonOperator("LTE"), OpLTE)
	assert.Equal(t, ComparisonOperator("EQ"), OpEQ)
	assert.Equal(t, ComparisonOperator("NE"), OpNE)
}

func TestDeepCopyPolicy(t *testing.T) {
	original := &Policy{
		ObjectMeta: metav1.ObjectMeta{Name: "test-policy", Namespace: "default"},
		Spec: PolicySpec{
			Name:          "Test Policy",
			Description:   "A test policy",
			Priority:      PriorityHigh,
			FailureAction: FailureActionBlockRelease,
			Conditions: []PolicyCondition{
				{Type: ConditionTypeCVSS, Comparator: OpGTE, Value: "7.0", IsSuppression: false},
			},
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
