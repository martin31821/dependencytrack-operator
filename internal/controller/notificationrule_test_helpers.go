/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

// configSchemaPathRE matches /v1/notification/publisher/{uuid}/configSchema.
var configSchemaPathRE = regexp.MustCompile(`/v1/notification/publisher/[^/]+/configSchema$`)

// buildTestScheme creates a runtime.Scheme with all required types registered.
func buildTestScheme(t *testing.T) *runtime.Scheme {
	s := runtime.NewScheme()
	if err := scheme.AddToScheme(s); err != nil {
		t.Fatalf("failed to add base scheme types: %v", err)
	}
	if err := dependencytrackv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("failed to add operator types to scheme: %v", err)
	}
	return s
}

// buildFakeClientWithTeams creates a fake client pre-loaded with K8s Team CRDs.
// k8sTeams maps team name → team UUID.
func buildFakeClientWithTeams(t *testing.T, scheme *runtime.Scheme, k8sTeams map[string]string) client.Client {
	objects := []client.Object{}
	for name, uuid := range k8sTeams {
		team := &dependencytrackv1alpha1.Team{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "dependencytrack.mko.dev/v1alpha1",
				Kind:       "Team",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: dependencytrackv1alpha1.TeamSpec{
				Name: name,
			},
		}
		if uuid != "" {
			team.Status = dependencytrackv1alpha1.TeamStatus{
				UUID: uuid,
			}
		}
		objects = append(objects, team)
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build()
	return fakeClient
}

// buildMockAPIClient creates a mock DT API client backed by an httptest.Server.
// remoteTeams maps team name → team UUID (for GET /v1/team/{uuid} responses).
func buildMockAPIClient(remoteTeams map[string]string) *dtapi.APIClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stripped := r.URL.Path
		if len(stripped) > 4 && stripped[:4] == "/api" {
			stripped = stripped[4:]
		}

		if r.Method == http.MethodGet && strings.HasPrefix(stripped, "/v1/team/") {
			parts := strings.Split(strings.TrimPrefix(stripped, "/v1/team/"), "/")
			if len(parts) == 1 {
				teamUUID := parts[0]
				for name, uuid := range remoteTeams {
					if uuid == teamUUID {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]string{
							"uuid": uuid,
							"name": name,
						})
						return
					}
				}
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"title":"Not Found","status":404}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	cfg := dtapi.NewConfiguration()
	cfg.HTTPClient = server.Client()
	cfg.Servers = dtapi.ServerConfigurations{{URL: server.URL}}
	return dtapi.NewAPIClient(cfg)
}

// buildMockAPIClientWithTeamsAndRules creates a mock DT API client for team convergence tests.
func buildMockAPIClientWithTeamsAndRules(
	remoteTeams map[string]string,
	teamRules map[string][]dtapi.Team,
	updateCalled *bool,
) *dtapi.APIClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stripped := r.URL.Path
		if len(stripped) > 4 && stripped[:4] == "/api" {
			stripped = stripped[4:]
		}

		if r.Method == http.MethodGet && strings.HasPrefix(stripped, "/v1/team/") {
			parts := strings.Split(strings.TrimPrefix(stripped, "/v1/team/"), "/")
			if len(parts) == 1 {
				teamUUID := parts[0]
				for name, uuid := range remoteTeams {
					if uuid == teamUUID {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]string{
							"uuid": uuid,
							"name": name,
						})
						return
					}
				}
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"title":"Not Found","status":404}`))
			return
		}

		if r.Method == http.MethodGet && stripped == "/v1/notification/rule" {
			rules := make([]dtapi.NotificationRule, 0)
			for ruleUUID, teams := range teamRules {
				rules = append(rules, dtapi.NotificationRule{Uuid: ruleUUID, Teams: teams})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rules)
			return
		}

		if r.Method == http.MethodPost && stripped == "/v1/notification/rule" {
			var req dtapi.UpdateNotificationRuleRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			*updateCalled = true
			if req.Teams != nil {
				teamRules[req.Uuid] = req.Teams
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dtapi.NotificationRule{Uuid: req.Uuid, Teams: req.Teams})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))

	cfg := dtapi.NewConfiguration()
	cfg.HTTPClient = server.Client()
	cfg.Servers = dtapi.ServerConfigurations{{URL: server.URL}}
	return dtapi.NewAPIClient(cfg)
}

// buildMockAPIClientWithProjectsAndRules creates a mock DT API client for project convergence tests.
func buildMockAPIClientWithProjectsAndRules(
	projectRules map[string][]dtapi.Project,
	updateCalled *bool,
) *dtapi.APIClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stripped := r.URL.Path
		if len(stripped) > 4 && stripped[:4] == "/api" {
			stripped = stripped[4:]
		}

		if r.Method == http.MethodGet && stripped == "/v1/notification/rule" {
			rules := make([]dtapi.NotificationRule, 0)
			for ruleUUID, projects := range projectRules {
				rules = append(rules, dtapi.NotificationRule{Uuid: ruleUUID, Projects: projects})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rules)
			return
		}

		if r.Method == http.MethodPost && stripped == "/v1/notification/rule" {
			var req dtapi.UpdateNotificationRuleRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			*updateCalled = true
			if req.Projects != nil {
				projectRules[req.Uuid] = req.Projects
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dtapi.NotificationRule{Uuid: req.Uuid, Projects: req.Projects})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))

	cfg := dtapi.NewConfiguration()
	cfg.HTTPClient = server.Client()
	cfg.Servers = dtapi.ServerConfigurations{{URL: server.URL}}
	return dtapi.NewAPIClient(cfg)
}
