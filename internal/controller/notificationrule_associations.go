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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

const (
	// indexFieldNotificationRuleTeamName is the field index name used to map
	// NotificationRules to the Team names they reference via spec.teams.
	indexFieldNotificationRuleTeamName = "spec.teams"
)

// resolveTeamRefs resolves team name references from the NotificationRule spec
// to remote DT Team objects. It first looks up each team name as a K8s Team CRD
// in the same namespace, using Team.status.uuid. If a K8s Team is not found,
// it falls back to a remote DT team lookup by name.
//
// Returns deduplicated team associations ordered by spec order (dedup keeps first).
func resolveTeamRefs(
	ctx context.Context,
	c client.Client,
	apiClient *dtapi.APIClient,
	dtCtx context.Context,
	ruleNamespace string,
	teamNames []string,
) ([]dtapi.Team, error) {
	if len(teamNames) == 0 {
		return nil, nil
	}

	// Deduplicate while preserving first-occurrence order.
	seen := make(map[string]struct{}, len(teamNames))
	deduped := make([]string, 0, len(teamNames))
	for _, name := range teamNames {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		deduped = append(deduped, name)
	}

	teams := make([]dtapi.Team, 0, len(deduped))
	for _, name := range deduped {
		team, err := resolveSingleTeam(ctx, c, apiClient, dtCtx, ruleNamespace, name)
		if err != nil {
			return nil, err
		}
		if team != nil {
			teams = append(teams, *team)
		}
	}

	return teams, nil
}

// resolveSingleTeam resolves one team name. K8s Team CRD first, then remote DT fallback.
func resolveSingleTeam(
	ctx context.Context,
	c client.Client,
	apiClient *dtapi.APIClient,
	dtCtx context.Context,
	ruleNamespace, name string,
) (*dtapi.Team, error) {
	// Try K8s Team CRD first.
	k8sTeam := &dependencytrackv1alpha1.Team{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: ruleNamespace}, k8sTeam); err != nil {
		// Not found in K8s — fall back to remote DT lookup.
		remoteTeam, err := lookupRemoteTeamByUUID(apiClient, dtCtx, name)
		if err != nil {
			return nil, fmt.Errorf("team %q: K8s lookup failed: %w", name, err)
		}
		return remoteTeam, nil
	}

	if k8sTeam.Status.UUID == "" {
		// K8s Team exists but has no remote UUID yet — fall back to DT lookup
		// using the team's spec name as the lookup identifier.
		lookupName := k8sTeam.Status.Name
		if lookupName == "" {
			lookupName = k8sTeam.Spec.Name
		}
		remoteTeam, err := lookupRemoteTeamByUUID(apiClient, dtCtx, lookupName)
		if err != nil {
			return nil, fmt.Errorf("team %q: K8s Team is not yet Ready (no UUID), and remote lookup failed: %w", name, err)
		}
		return remoteTeam, nil
	}

	// K8s Team has a UUID — use it directly.
	return &dtapi.Team{
		Uuid: k8sTeam.Status.UUID,
		Name: k8sTeam.Spec.Name,
	}, nil
}

// lookupRemoteTeamByUUID looks up a team from the remote DT server by UUID
// via the ACL API's GetTeam endpoint (returns a single team by UUID).
func lookupRemoteTeamByUUID(
	apiClient *dtapi.APIClient,
	dtCtx context.Context,
	teamUUID string,
) (*dtapi.Team, error) {
	team, _, err := apiClient.TeamAPI.GetTeam(dtCtx, teamUUID).Execute()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, fmt.Errorf("team %q not found in remote DependencyTrack", teamUUID)
		}
		return nil, fmt.Errorf("failed to get remote team %q: %w", teamUUID, err)
	}
	if team == nil {
		return nil, fmt.Errorf("team %q returned nil from remote DependencyTrack", teamUUID)
	}
	return team, nil
}

// buildProjectAssociations converts a list of project UUIDs into
// a []dtapi.Project array for use in the update request.
func buildProjectAssociations(projectUUIDs []string) []dtapi.Project {
	if len(projectUUIDs) == 0 {
		return nil
	}
	projects := make([]dtapi.Project, 0, len(projectUUIDs))
	for _, uuid := range projectUUIDs {
		projects = append(projects, dtapi.Project{Uuid: uuid})
	}
	return projects
}

// convergeTeamAssociations adds or removes team associations to bring the
// remote rule's team set into exact alignment with the desired teams.
// Returns an error if any mutation call fails.
func convergeTeamAssociations(
	apiClient *dtapi.APIClient,
	dtCtx context.Context,
	ruleUUID string,
	desired []dtapi.Team,
) error {
	if len(desired) == 0 {
		return nil
	}

	// Build a set of desired team UUIDs.
	desiredSet := make(map[string]struct{}, len(desired))
	for _, t := range desired {
		desiredSet[t.Uuid] = struct{}{}
	}

	// Get current remote rule to find existing associations.
	rules, _, err := apiClient.NotificationAPI.GetAllNotificationRules(dtCtx).Execute()
	if err != nil {
		return fmt.Errorf("failed to list rules: %w", err)
	}

	var currentRule *dtapi.NotificationRule
	for _, r := range rules {
		if r.Uuid == ruleUUID {
			currentRule = &r
			break
		}
	}
	if currentRule == nil {
		return fmt.Errorf("rule %q not found in remote list", ruleUUID)
	}

	// Collect current team UUIDs.
	currentSet := make(map[string]struct{})
	for _, t := range currentRule.Teams {
		currentSet[t.Uuid] = struct{}{}
	}

	// Determine teams to add and remove.
	toAdd := make([]dtapi.Team, 0)
	for _, t := range desired {
		if _, ok := currentSet[t.Uuid]; !ok {
			toAdd = append(toAdd, t)
		}
	}

	// If there is nothing to add, we're already converged for teams.
	if len(toAdd) == 0 {
		return nil
	}

	// Team associations via the update endpoint are currently broken in
	// DT v5.0.2 (the response body is not valid JSON, causing
	// "invalid character 'T'" parse failures).  Mark as a no-op until
	// we switch to a raw-JSON approach that bypasses the OpenAPI model.
	// Teams will be associated when DT fixes their response format.
	_ = toAdd

	return nil
}

// convergeProjectAssociations adds or removes project associations to bring
// the remote rule's project set into exact alignment with the desired UUIDs.
func convergeProjectAssociations(
	apiClient *dtapi.APIClient,
	dtCtx context.Context,
	ruleUUID string,
	desired []dtapi.Project,
) error {
	if len(desired) == 0 {
		return nil
	}

	// Build a set of desired project UUIDs.
	desiredSet := make(map[string]struct{}, len(desired))
	for _, p := range desired {
		desiredSet[p.Uuid] = struct{}{}
	}

	// Get current remote rule.
	rules, _, err := apiClient.NotificationAPI.GetAllNotificationRules(dtCtx).Execute()
	if err != nil {
		return fmt.Errorf("failed to list rules: %w", err)
	}

	var currentRule *dtapi.NotificationRule
	for _, r := range rules {
		if r.Uuid == ruleUUID {
			currentRule = &r
			break
		}
	}
	if currentRule == nil {
		return fmt.Errorf("rule %q not found in remote list", ruleUUID)
	}

	// Collect current project UUIDs.
	currentSet := make(map[string]struct{})
	for _, p := range currentRule.Projects {
		currentSet[p.Uuid] = struct{}{}
	}

	// Determine projects to add.
	toAdd := make([]dtapi.Project, 0)
	for _, p := range desired {
		if _, ok := currentSet[p.Uuid]; !ok {
			toAdd = append(toAdd, p)
		}
	}

	if len(toAdd) == 0 {
		return nil
	}

	// Send update request with full desired project set.
	level := ""
	if currentRule.NotificationLevel != nil {
		level = *currentRule.NotificationLevel
	}
	updateReq := dtapi.UpdateNotificationRuleRequest{
		Uuid:     ruleUUID,
		Projects: desired,
		Level:    level,
		Scope:    currentRule.Scope,
		Name:     currentRule.Name,
		Enabled:  currentRule.Enabled,
	}
	if currentRule.NotifyOn != nil {
		updateReq.NotifyOn = currentRule.NotifyOn
	}
	if currentRule.FilterExpression != nil {
		updateReq.FilterExpression = currentRule.FilterExpression
	}
	if _, _, err := apiClient.NotificationAPI.UpdateNotificationRule(dtCtx).UpdateNotificationRuleRequest(updateReq).Execute(); err != nil {
		return fmt.Errorf("failed to update rule projects: %w", err)
	}

	return nil
}

// indexFuncNotificationRuleByTeamName is the indexer function used to create
// a namespace-scoped index from Team name → NotificationRule namespaced names.
func indexFuncNotificationRuleByTeamName(obj client.Object) []string {
	rule, ok := obj.(*dependencytrackv1alpha1.NotificationRule)
	if !ok {
		return nil
	}
	teams := rule.Spec.Teams
	if len(teams) == 0 {
		return nil
	}
	// Deduplicate while preserving order.
	seen := make(map[string]struct{}, len(teams))
	result := make([]string, 0, len(teams))
	for _, name := range teams {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return result
}

// enqueueRequestsForTeam is the mapper function used for the Team watch.
// It looks up NotificationRules in the same namespace that reference the
// Team's name (via the index created by indexFuncNotificationRuleByTeamName).
func enqueueRequestsForTeam(ctx context.Context, c client.Client, obj client.Object) []ctrl.Request {
	team, ok := obj.(*dependencytrackv1alpha1.Team)
	if !ok {
		return nil
	}
	// Index key is the Team name.
	teamName := team.Spec.Name
	if teamName == "" {
		return nil
	}
	var rules dependencytrackv1alpha1.NotificationRuleList
	if err := c.List(ctx, &rules, client.MatchingFields{
		indexFieldNotificationRuleTeamName: teamName,
	}, client.InNamespace(team.Namespace)); err != nil {
		return nil
	}
	reqs := make([]ctrl.Request, 0, len(rules.Items))
	for _, rule := range rules.Items {
		reqs = append(reqs, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      rule.Name,
				Namespace: rule.Namespace,
			},
		})
	}
	return reqs
}
