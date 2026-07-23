package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

func TestDebugTeamAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/v1/team/") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"title":"Not Found","status":404}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := dtapi.NewConfiguration()
	cfg.Servers = dtapi.ServerConfigurations{{URL: server.URL}}
	apiClient := dtapi.NewAPIClient(cfg)

	ctx := context.Background()
	team, resp, err := apiClient.TeamAPI.GetTeam(ctx, "ghost-team").Execute()
	fmt.Printf("Result: team=%v resp=%v err=%v\n", team, resp, err)
	if err != nil {
		fmt.Printf("Error contains 404: %v\n", strings.Contains(err.Error(), "404"))
	}
}
