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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

func TestLookupTeamByUUIDToleratesMissingLastPasswordChange(t *testing.T) {
	var authorization string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization = r.Header.Get("Authorization")
		if r.URL.Path != "/api/v1/team/team-uuid" {
			t.Fatalf("unexpected request path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"operators",
			"uuid":"team-uuid",
			"managedUsers":[{"username":"alice"}],
			"permissions":[{"name":"PORTFOLIO_MANAGEMENT","managedUsers":[{"username":"alice"}]}]
		}`))
	}))
	defer server.Close()

	apiClient := newTeamTestAPIClient(server)
	authCtx := context.WithValue(context.Background(), dtapi.ContextAccessToken, "session-token")
	team, status, err := lookupTeamByUUID(authCtx, apiClient, "team-uuid")
	if err != nil {
		t.Fatalf("lookupTeamByUUID returned an error: %v", err)
	}
	if status != http.StatusOK || team.UUID != "team-uuid" || team.Name != "operators" {
		t.Fatalf("unexpected lookup result: status=%d team=%+v", status, team)
	}
	if len(team.Permissions) != 1 || team.Permissions[0].Name != "PORTFOLIO_MANAGEMENT" {
		t.Fatalf("unexpected permissions: %+v", team.Permissions)
	}
	if authorization != "Bearer session-token" {
		t.Fatalf("unexpected Authorization header %q", authorization)
	}
}

func TestLookupTeamUUIDByNameRequiresExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("name"); got != "target team" {
			t.Fatalf("unexpected name query %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"other team","uuid":"wrong"},
			{"name":"target team","uuid":"right","managedUsers":[{"username":"alice"}]}
		]`))
	}))
	defer server.Close()

	apiClient := newTeamTestAPIClient(server)
	authCtx := context.WithValue(context.Background(), dtapi.ContextAccessToken, "session-token")
	uuid, err := (&TeamReconciler{}).lookupTeamUUIDByName(authCtx, apiClient, "target team")
	if err != nil {
		t.Fatalf("lookupTeamUUIDByName returned an error: %v", err)
	}
	if uuid != "right" {
		t.Fatalf("expected exact-match UUID, got %q", uuid)
	}
}

func newTeamTestAPIClient(server *httptest.Server) *dtapi.APIClient {
	cfg := dtapi.NewConfiguration()
	cfg.Servers = dtapi.ServerConfigurations{{URL: server.URL + "/api"}}
	cfg.HTTPClient = server.Client()
	return dtapi.NewAPIClient(cfg)
}
