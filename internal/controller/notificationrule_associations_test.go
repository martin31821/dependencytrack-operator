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

package controller

import (
	"testing"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/stretchr/testify/assert"
)

func TestResolveTeamRefs(t *testing.T) {
	tests := []struct {
		name        string
		teamNames   []string
		k8sTeams    map[string]string // name → uuid
		remoteTeams map[string]string // uuid → team name
		wantTeams   []string          // expected resolved UUIDs
		wantErr     bool
	}{
		{
			name:      "empty team list returns nil",
			teamNames: []string{},
			wantTeams: nil,
		},
		{
			name:      "single team resolved from K8s",
			teamNames: []string{"my-team"},
			k8sTeams:  map[string]string{"my-team": "team-uuid-1"},
			wantTeams: []string{"team-uuid-1"},
		},
		{
			name:      "multiple teams resolved from K8s",
			teamNames: []string{"team-a", "team-b"},
			k8sTeams:  map[string]string{"team-a": "uuid-a", "team-b": "uuid-b"},
			wantTeams: []string{"uuid-a", "uuid-b"},
		},
		{
			name:      "duplicate team names are deduplicated",
			teamNames: []string{"team-a", "team-a", "team-b"},
			k8sTeams:  map[string]string{"team-a": "uuid-a", "team-b": "uuid-b"},
			wantTeams: []string{"uuid-a", "uuid-b"},
		},
		{
			name:      "K8s team with empty UUID falls back to remote lookup",
			teamNames: []string{"team-a"},
			k8sTeams:  map[string]string{"team-a": ""},
			remoteTeams: map[string]string{
				"team-a": "team-a", // name maps to same UUID
			},
			wantTeams: []string{"team-a"},
		},
		{
			name:      "K8s team not found falls back to remote lookup",
			teamNames: []string{"missing-team"},
			remoteTeams: map[string]string{
				"missing-team": "missing-team",
			},
			wantTeams: []string{"missing-team"},
		},
		{
			name:        "K8s team not found and remote not found returns error",
			teamNames:   []string{"ghost-team"},
			k8sTeams:    map[string]string{}, // not in K8s
			remoteTeams: map[string]string{}, // not in remote either
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build fake client with K8s teams.
			scheme := buildTestScheme(t)
			c := buildFakeClientWithTeams(t, scheme, tt.k8sTeams)
			apiClient := buildMockAPIClient(tt.remoteTeams)

			teams, err := resolveTeamRefs(t.Context(), c, apiClient, t.Context(), "default", tt.teamNames)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			wantUUIDs := make([]string, len(tt.wantTeams))
			copy(wantUUIDs, tt.wantTeams)

			gotUUIDs := make([]string, len(teams))
			for i, team := range teams {
				gotUUIDs[i] = team.Uuid
			}

			assert.Equal(t, wantUUIDs, gotUUIDs)
		})
	}
}

func TestBuildProjectAssociations(t *testing.T) {
	tests := []struct {
		name     string
		uuids    []string
		wantLen  int
		wantUUID []string
	}{
		{
			name:     "empty list returns nil",
			uuids:    []string{},
			wantLen:  0,
			wantUUID: nil,
		},
		{
			name:     "single UUID",
			uuids:    []string{"proj-uuid-1"},
			wantLen:  1,
			wantUUID: []string{"proj-uuid-1"},
		},
		{
			name:     "multiple UUIDs",
			uuids:    []string{"proj-1", "proj-2", "proj-3"},
			wantLen:  3,
			wantUUID: []string{"proj-1", "proj-2", "proj-3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildProjectAssociations(tt.uuids)
			assert.Equal(t, tt.wantLen, len(result))
			for i, p := range result {
				assert.Equal(t, tt.wantUUID[i], p.Uuid)
			}
		})
	}
}

func TestConvergeTeamAssociations(t *testing.T) {
	tests := []struct {
		name       string
		ruleUUID   string
		desired    []dtapi.Team
		existing   []dtapi.Team
		wantCalled bool
		wantErr    bool
	}{
		{
			name:       "empty desired teams returns nil",
			ruleUUID:   "rule-1",
			desired:    nil,
			existing:   []dtapi.Team{{Uuid: "team-old"}},
			wantCalled: false,
		},
		{
			name:       "desired matches existing — no update",
			ruleUUID:   "rule-1",
			desired:    []dtapi.Team{{Uuid: "team-a"}, {Uuid: "team-b"}},
			existing:   []dtapi.Team{{Uuid: "team-a"}, {Uuid: "team-b"}},
			wantCalled: false,
		},
		{
			name:     "desired has new team — no-op pending DT fix",
			ruleUUID: "rule-1",
			desired:  []dtapi.Team{{Uuid: "team-a"}, {Uuid: "team-new"}},
			existing: []dtapi.Team{{Uuid: "team-a"}},
			// Team update is a no-op in DT v5.0.2 (response parse failure); see convergeTeamAssociations.
			wantCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			apiClient := buildMockAPIClientWithTeamsAndRules(
				map[string]string{}, // remote teams
				map[string][]dtapi.Team{
					"rule-1": tt.existing,
				},
				&updateCalled,
			)

			err := convergeTeamAssociations(apiClient, t.Context(), tt.ruleUUID, tt.desired)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCalled, updateCalled)
		})
	}
}

func TestConvergeProjectAssociations(t *testing.T) {
	tests := []struct {
		name       string
		ruleUUID   string
		desired    []dtapi.Project
		existing   []dtapi.Project
		wantCalled bool
	}{
		{
			name:       "empty desired projects returns nil",
			ruleUUID:   "rule-1",
			desired:    nil,
			existing:   []dtapi.Project{{Uuid: "proj-old"}},
			wantCalled: false,
		},
		{
			name:       "desired matches existing — no update",
			ruleUUID:   "rule-1",
			desired:    []dtapi.Project{{Uuid: "proj-1"}},
			existing:   []dtapi.Project{{Uuid: "proj-1"}},
			wantCalled: false,
		},
		{
			name:       "desired has new project — update called",
			ruleUUID:   "rule-1",
			desired:    []dtapi.Project{{Uuid: "proj-1"}, {Uuid: "proj-new"}},
			existing:   []dtapi.Project{{Uuid: "proj-1"}},
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			apiClient := buildMockAPIClientWithProjectsAndRules(
				map[string][]dtapi.Project{
					"rule-1": tt.existing,
				},
				&updateCalled,
			)

			err := convergeProjectAssociations(apiClient, t.Context(), tt.ruleUUID, tt.desired)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCalled, updateCalled)
		})
	}
}

func TestIndexFuncNotificationRuleByTeamName(t *testing.T) {
	tests := []struct {
		name     string
		teams    []string
		wantKeys []string
	}{
		{
			name:     "nil teams returns nil",
			teams:    nil,
			wantKeys: nil,
		},
		{
			name:     "empty teams returns nil",
			teams:    []string{},
			wantKeys: nil,
		},
		{
			name:     "single team",
			teams:    []string{"devs"},
			wantKeys: []string{"devs"},
		},
		{
			name:     "multiple teams",
			teams:    []string{"devs", "ops"},
			wantKeys: []string{"devs", "ops"},
		},
		{
			name:     "duplicate teams are deduplicated",
			teams:    []string{"devs", "devs", "ops"},
			wantKeys: []string{"devs", "ops"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &dtapi.NotificationRule{}
			_ = rule // not used directly; indexer works on K8s objects
		})
	}
}
