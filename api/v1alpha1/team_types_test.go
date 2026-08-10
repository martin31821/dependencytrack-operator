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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestTeamGVKFromScheme mirrors the Policy/NotificationPublisher scheme
// coverage so a broken init()/SchemeBuilder registration fails loudly.
func TestTeamGVKFromScheme(t *testing.T) {
	testScheme := runtime.NewScheme()
	err := SchemeBuilder.AddToScheme(testScheme)
	assert.NoError(t, err)

	gvks, _, err := testScheme.ObjectKinds(&Team{})
	assert.NoError(t, err)
	assert.Len(t, gvks, 1)
	assert.Equal(t, "dependencytrack.mko.dev", gvks[0].Group)
	assert.Equal(t, "v1alpha1", gvks[0].Version)
	assert.Equal(t, "Team", gvks[0].Kind)

	listGvks, _, err := testScheme.ObjectKinds(&TeamList{})
	assert.NoError(t, err)
	assert.Len(t, listGvks, 1)
	assert.Equal(t, "TeamList", listGvks[0].Kind)
}

// TestTeamOIDCNilVsEmptyPointer locks the sentinel semantics the controller
// (T04) and webhook (T02) rely on: a nil *TeamOIDCConfig disables OIDC
// entirely, whereas a non-nil config with an empty Groups slice is an
// intentional "mapped to nothing" state.
func TestTeamOIDCNilVsEmptyPointer(t *testing.T) {
	disabled := TeamSpec{OIDC: nil}
	enabledNoGroups := TeamSpec{OIDC: &TeamOIDCConfig{Groups: []string{}}}

	assert.True(t, disabled.OIDC == nil, "nil pointer marks OIDC-disabled")
	assert.False(t, enabledNoGroups.OIDC == nil, "non-nil pointer marks OIDC-enabled")
	assert.Empty(t, enabledNoGroups.OIDC.Groups, "empty slice is not nil")

	bDisabled, err := json.Marshal(disabled)
	assert.NoError(t, err)
	assert.NotContains(t, string(bDisabled), "\"oidc\"",
		"a nil OIDC pointer must be omitted from the serialized spec")

	bEnabled, err := json.Marshal(enabledNoGroups)
	assert.NoError(t, err)
	assert.Contains(t, string(bEnabled), "\"oidc\":{}",
		"a non-nil OIDC config must survive as an empty object")
}

// TestTeamOIDCGroupsCasePreservation is the central negative test for this
// task: group names must round-trip through JSON byte-for-byte. Any
// implicit TrimSpace/ToLower/UpperFold normalization would break IdP claim
// matching, so we assert several adversarial casings are preserved verbatim.
func TestTeamOIDCGroupsCasePreservation(t *testing.T) {
	adversarial := []string{
		"Admins",            // mixed case
		"admins",            // lower
		"ADMINS",            // upper
		" Admins ",          // surrounding whitespace must be retained
		"Département-Vente", // unicode letters
		"group.with.dots",
		"", // empty string is a legal group name; must persist
	}

	team := Team{
		Spec: TeamSpec{
			Name: "svc-team",
			OIDC: &TeamOIDCConfig{Groups: adversarial},
		},
	}

	data, err := json.Marshal(team)
	assert.NoError(t, err)

	var decoded Team
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	decodedGroups := decoded.Spec.OIDC.Groups
	if decodedGroups == nil {
		t.Fatal("decoded Spec.OIDC must be non-nil")
	}
	assert.Equal(t, adversarial, decodedGroups,
		"group names must round-trip case/whitespace-unaltered")

	// Positive: spaced/original spellings co-exist with their lowercased twins.
	assert.Contains(t, decodedGroups, " Admins ")
	assert.Contains(t, decodedGroups, "Département-Vente")

	// Negative: prove we did NOT normalize. A naive normalizer
	// (TrimSpace + ToLower) collapses this multiset to a different slice,
	// so the decoded value must diverge from it.
	normalized := make([]string, len(adversarial))
	for i, g := range adversarial {
		normalized[i] = strings.ToLower(strings.TrimSpace(g))
	}
	assert.NotEqual(t, normalized, decodedGroups,
		"group names must not be trimmed/lowercased")
}

// TestTeamStatusOwnedMappingsRoundTrip covers the status-side schema: the
// controller reports discovered IdP-to-Team mappings via
// status.oidc.ownedMappings. Each field is opaque from the DT API and must
// survive unmarshalling intact.
func TestTeamStatusOwnedMappingsRoundTrip(t *testing.T) {
	team := Team{
		ObjectMeta: metav1.ObjectMeta{Name: "rt-team", Namespace: "deptrack"},
		Status: TeamStatus{
			OIDC: &TeamOIDCStatus{
				OwnedMappings: []OwnedOIDCMapping{
					{
						GroupName:   "Ops-Engineers",
						GroupUUID:   "grp-aaa-111",
						TeamUUID:    "team-bbb-222",
						MappingUUID: "map-ccc-333",
					},
				},
			},
		},
	}

	data, err := json.Marshal(team)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "ownedMappings"))

	var decoded Team
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.NotNil(t, decoded.Status.OIDC)
	assert.Len(t, decoded.Status.OIDC.OwnedMappings, 1)
	assert.Equal(t, "Ops-Engineers", decoded.Status.OIDC.OwnedMappings[0].GroupName)
	assert.Equal(t, "grp-aaa-111", decoded.Status.OIDC.OwnedMappings[0].GroupUUID)
	assert.Equal(t, "team-bbb-222", decoded.Status.OIDC.OwnedMappings[0].TeamUUID)
	assert.Equal(t, "map-ccc-333", decoded.Status.OIDC.OwnedMappings[0].MappingUUID)
}

// TestTeamStatusOIDCNilMeansUnset encodes the status sentinel: a nil
// status.OIDC means "never reconciled", an empty OwnedMappings slice means
// "reconciled, no mappings exist".
func TestTeamStatusOIDCNilMeansUnset(t *testing.T) {
	unsettled := TeamStatus{OIDC: nil}
	settled := TeamStatus{OIDC: &TeamOIDCStatus{OwnedMappings: nil}}

	assert.True(t, unsettled.OIDC == nil)
	assert.False(t, settled.OIDC == nil)
	assert.Nil(t, settled.OIDC.OwnedMappings)
}

// TestTeamDeepCopyIndependence proves the controller-gen-generated Clone
// path copies the new nested slice types (Groups and OwnedMappings) rather
// than aliasing them. Mutation of the copy must never leak into the original.
func TestTeamDeepCopyIndependence(t *testing.T) {
	original := &Team{
		ObjectMeta: metav1.ObjectMeta{Name: "dc-team", Namespace: "default"},
		Spec: TeamSpec{
			Name: "Deep Copy Team",
			OIDC: &TeamOIDCConfig{Groups: []string{"Admins", "Auditors"}},
		},
		Status: TeamStatus{
			UUID: "team-uuid-original",
			OIDC: &TeamOIDCStatus{
				OwnedMappings: []OwnedOIDCMapping{
					{GroupName: "Admins", GroupUUID: "g1", TeamUUID: "t1", MappingUUID: "m1"},
				},
			},
		},
	}

	cp := original.DeepCopy()
	assert.Equal(t, original, cp)

	// Mutate the copy's spec groups.
	cp.Spec.OIDC.Groups[0] = "HACKED"
	assert.NotEqual(t, "HACKED", original.Spec.OIDC.Groups[0],
		"spec.OIDC.Groups must be independently cloned")

	// Append to the copy's owned mappings.
	cp.Status.OIDC.OwnedMappings = append(cp.Status.OIDC.OwnedMappings,
		OwnedOIDCMapping{GroupName: "Extra"})
	assert.Len(t, original.Status.OIDC.OwnedMappings, 1,
		"status.OIDC.OwnedMappings must be independently cloned")

	// Point-swap the whole OIDC config to prove pointer isolation.
	cp.Spec.OIDC = &TeamOIDCConfig{Groups: []string{"other"}}
	assert.NotSame(t, original.Spec.OIDC, cp.Spec.OIDC)
	assert.Equal(t, "Admins", original.Spec.OIDC.Groups[0])
}

// TestTeamDeepCopyObject_NilSafe exercises the runtime.Object interface on a
// freshly-allocated Team to catch regressions in the generated boilerplate.
func TestTeamDeepCopyObject_NilSafe(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		obj := (&Team{}).DeepCopyObject()
		assert.NotNil(t, obj)
		assert.IsType(t, &Team{}, obj)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var nilTeam *Team
		assert.Nil(t, nilTeam.DeepCopyObject())
	})
}
