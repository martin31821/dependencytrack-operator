package dtapi

import (
	"encoding/json"
	"testing"
)

// Newer DependencyTrack releases add response fields (e.g. 5.2 added
// Team.serviceAccounts and Policy.invertTagMatch). The generated models must
// ignore them; re-check after regenerating the client, since openapi-generator
// emits decoder.DisallowUnknownFields() by default.
func TestModelsIgnoreUnknownFields(t *testing.T) {
	var team Team
	if err := json.Unmarshal([]byte(`{"name":"devops","uuid":"6f1c1c3e-9a43-4a0e-8a5c-1d3c2f6d0b11","serviceAccounts":[]}`), &team); err != nil {
		t.Fatalf("Team: %v", err)
	}
	if team.Name != "devops" {
		t.Fatalf("Team.Name = %q", team.Name)
	}

	var policy Policy
	if err := json.Unmarshal([]byte(`{"name":"vuln","operator":"ANY","uuid":"0c7a7f0e-3a7c-4c43-9a8e-6f4f8b1a2d33","violationState":"FAIL","invertTagMatch":false}`), &policy); err != nil {
		t.Fatalf("Policy: %v", err)
	}
}
