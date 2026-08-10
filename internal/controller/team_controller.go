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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	logr "github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/dependencytrack"
)

const (
	teamFinalizer  = "dependencytrack.mko.dev/finalizer"
	reasonAPIError = "APIError"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	DTProvider *dependencytrack.ClientProvider
}

// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dependencytrack.mko.dev,resources=teams/finalizers,verbs=update

func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	team := &dependencytrackv1alpha1.Team{}
	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure our finalizer is registered before doing any external work.
	if !controllerutil.ContainsFinalizer(team, teamFinalizer) {
		controllerutil.AddFinalizer(team, teamFinalizer)
		if err := r.Update(ctx, team); err != nil {
			return ctrl.Result{}, err
		}
		// Requeue; the update will trigger a new reconcile with the finalizer in place.
		return ctrl.Result{}, nil
	}

	if !team.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, team)
	}

	return r.reconcileUpsert(ctx, team)
}

// lookupTeamUUIDByName returns the UUID of the DependencyTrack team whose name
// matches the given value, or "" if no such team exists. It makes team creation
// idempotent across reconcile interruptions (e.g. a rolling restart that occurs
// after CreateTeam succeeds but before the UUID is persisted to status), which
// would otherwise wedge the reconcile loop with a 409 Conflict.
//
// It deliberately issues a raw GET /v1/team?name= request instead of using the
// generated dtapi TeamAPI.GetTeams list endpoint: the generated Team model marks
// lastPasswordChange as a required field, and DependencyTrack omits it for some
// teams (e.g. LDAP-mapped), which makes the generated list decoder reject
// otherwise-valid responses.
func (r *TeamReconciler) lookupTeamUUIDByName(authCtx context.Context, apiClient *dtapi.APIClient, name string) (string, error) {
	token, _ := authCtx.Value(dtapi.ContextAccessToken).(string)
	if token == "" {
		return "", fmt.Errorf("no bearer token in auth context")
	}
	cfg := apiClient.GetConfig()
	baseURL, err := cfg.ServerURLWithContext(authCtx, "TeamAPIService.GetTeams")
	if err != nil {
		return "", fmt.Errorf("resolve DependencyTrack server URL: %w", err)
	}

	u := strings.TrimRight(baseURL, "/") + "/v1/team?name=" + url.QueryEscape(name)
	req, err := http.NewRequestWithContext(authCtx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("looking up team %q by name: %s", name, resp.Status)
	}
	var decoded interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	switch v := decoded.(type) {
	case map[string]interface{}:
		matchedName, _ := v["name"].(string)
		if uuid, ok := v["uuid"].(string); ok && uuid != "" && matchedName == name {
			return uuid, nil
		}
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				matchedName, _ := m["name"].(string)
				if uuid, ok := m["uuid"].(string); ok && uuid != "" && matchedName == name {
					return uuid, nil
				}
			}
		}
	}
	return "", nil
}

// teamSnapshot contains only fields used during reconciliation. Dependency-Track
// may omit lastPasswordChange from nested managed users, while the generated
// client requires it and rejects the otherwise-valid Team response.
type teamSnapshot struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	Permissions []struct {
		Name string `json:"name"`
	} `json:"permissions"`
}

func lookupTeamByUUID(authCtx context.Context, apiClient *dtapi.APIClient, uuid string) (*teamSnapshot, int, error) {
	token, _ := authCtx.Value(dtapi.ContextAccessToken).(string)
	if token == "" {
		return nil, 0, fmt.Errorf("no bearer token in auth context")
	}
	cfg := apiClient.GetConfig()
	baseURL, err := cfg.ServerURLWithContext(authCtx, "TeamAPIService.GetTeam")
	if err != nil {
		return nil, 0, fmt.Errorf("resolve DependencyTrack server URL: %w", err)
	}
	req, err := http.NewRequestWithContext(authCtx, http.MethodGet,
		strings.TrimRight(baseURL, "/")+"/v1/team/"+url.PathEscape(uuid), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("DependencyTrack team lookup: %s", resp.Status)
	}

	var team teamSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decode DependencyTrack team response: %w", err)
	}
	return &team, resp.StatusCode, nil
}

func (r *TeamReconciler) reconcileUpsert(ctx context.Context, team *dependencytrackv1alpha1.Team) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack")
		setCondition(team, metav1.ConditionFalse, "CredentialsError", err.Error())
		_ = r.Status().Update(ctx, team)
		return ctrl.Result{}, err
	}

	var dtUUID string

	if team.Status.UUID != "" {
		// Team may already exist — fetch it by the UUID we recorded last time.
		existing, statusCode, err := lookupTeamByUUID(authCtx, apiClient, team.Status.UUID)
		switch {
		case err == nil:
			// Team still exists; rename if the spec changed.
			if existing.Name != team.Spec.Name {
				log.Info("updating team name", "uuid", existing.UUID, "oldName", existing.Name, "newName", team.Spec.Name)
				updated, _, err := apiClient.TeamAPI.UpdateTeam(authCtx).Team(dtapi.Team{
					Name: team.Spec.Name,
					Uuid: existing.UUID,
				}).Execute()
				if err != nil {
					return r.failStatus(ctx, team, "failed to update team name: "+err.Error(), err)
				}
				dtUUID = updated.Uuid
			} else {
				dtUUID = existing.UUID
			}
		case statusCode == http.StatusNotFound:
			// Team was deleted from DependencyTrack out-of-band; recreate it.
			log.Info("team not found in DependencyTrack, recreating", "uuid", team.Status.UUID)
		default:
			return r.failStatus(ctx, team, "failed to get team from DependencyTrack: "+err.Error(), err)
		}
	}

	if dtUUID == "" {
		// Interrupt-tolerant creation: a prior reconcile may have created the team
		// in DependencyTrack but been restarted (e.g. during a rolling update) before
		// team.Status.UUID was persisted. Adopting an existing team by name avoids a
		// redundant CreateTeam that would 409 Conflict and wedge the reconcile loop.
		adopted, err := r.lookupTeamUUIDByName(authCtx, apiClient, team.Spec.Name)
		if err != nil {
			return r.failStatus(ctx, team, "failed to look up team by name: "+err.Error(), err)
		}
		if adopted != "" {
			log.Info("adopting pre-existing team by name", "name", team.Spec.Name, "uuid", adopted)
			dtUUID = adopted
		} else {
			createTeam := dtapi.Team{
				Name: team.Spec.Name,
			}
			created, _, err := apiClient.TeamAPI.CreateTeam(authCtx).Team(createTeam).Execute()
			if err != nil {
				return r.failStatus(ctx, team, "failed to create team: "+err.Error(), err)
			}
			dtUUID = created.Uuid
			log.Info("created team in DependencyTrack", "uuid", dtUUID)
			r.Recorder.Eventf(team, "Normal", "TeamCreated", "Created team %q in DependencyTrack (uuid=%s)", team.Spec.Name, dtUUID)
		}
	}

	team.Status.UUID = dtUUID

	// Sync permissions if the spec declares any.
	if err := r.syncPermissions(ctx, apiClient, team, dtUUID); err != nil {
		return ctrl.Result{}, err
	}

	// Persist name in status for observability.
	team.Status.Name = team.Spec.Name

	// Reconcile OIDC group mappings when the spec opts in. Teams without an
	// OIDC stanza issue NO DependencyTrack API traffic and stay on baseline
	// behaviour; the guard is enforced both here and inside reconcileOIDC.
	//
	// The generated *dtapi.APIClient is bridged through the T03 seam
	// (dependencytrack.NewDtapiOIDCClient) here at the composition root;
	// tests bypass the adapter by handing reconcileOIDC a minimal OIDCAClient
	// stub directly.
	if team.Spec.OIDC != nil {
		api := dependencytrack.NewDtapiOIDCClient(apiClient.OidcAPI)
		// reconcileOIDC delegates to ReconcileOIDCMappings and DeleteMappingByUuid,
		// which the generated dtapi client authenticates only when the request
		// context carries dtapi.ContextAccessToken. Pass authCtx (populated by
		// ClientProvider.Get) rather than the ambient reconcile ctx, lest the
		// OIDC calls fly without a bearer token and DT 401s them.
		if res, err := r.reconcileOIDC(authCtx, api, team, dtUUID); err != nil || !res.IsZero() {
			return res, err
		}
	}

	setCondition(team, metav1.ConditionTrue, "TeamSynced", "Team successfully reconciled in DependencyTrack")
	if err := r.Status().Update(ctx, team); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

/*
	scrubOwnedOIDCMappings reclaims the OIDC group->team mapping edges this

// operator authorised for `team` (recorded in team.Status.OIDC.OwnedMappings)
// before the underlying Dependency-Track team object is torn down. Scoped to
// OWNED edges only: shared/group objects are never cascaded, mirroring the
// upsert-time guarantee that only binding edges are ever mutated.
//
// Matrix:
//   - team.Spec.OIDC == nil        -> nil, zero API traffic (never engaged).
//   - no owned mappings            -> nil (fast path; skips avail probe).
//   - OIDC not provisioned (!avail)-> ErrOIDCUnavailable, NO mutations;
//                                    caller defers+requeues w/ finalizer held.
//   - transport error on probe     -> wrapped and bubbled.
//   - available                    -> DELETE-MAPPING-BY-UUID per owned edge,
//                                    fail-STOP on first error so a partial
//                                    scrub is retried wholesale.
//
// ctx MUST carry the DependencyTrack bearer token (the auth context from
// DTProvider.Get) so probe/deletes propagate authentication.
*/
func (r *TeamReconciler) scrubOwnedOIDCMappings(ctx context.Context, api dependencytrack.OIDCAClient, team *dependencytrackv1alpha1.Team) error {
	// Vanilla teams opted out of OIDC entirely: zero API traffic.
	if team.Spec.OIDC == nil {
		return nil
	}

	owned := ownedStatusOf(team)
	if len(owned) == 0 {
		return nil
	}

	log := logf.FromContext(ctx)

	available, err := api.IsAvailable(ctx)
	if err != nil {
		return fmt.Errorf("probing OIDC availability while scrubbing owned mappings for team %q: %w", team.Spec.Name, err)
	}
	if !available {
		return dependencytrack.ErrOIDCUnavailable
	}

	for _, m := range owned {
		if err := api.DeleteMappingByUuid(ctx, m.MappingUUID); err != nil {
			log.Error(err, "failed to delete owned OIDC mapping", "mappingUuid", m.MappingUUID, "group", m.GroupName)
			return fmt.Errorf("deleting owned OIDC mapping %q (group %q) during team deletion: %w", m.MappingUUID, m.GroupName, err)
		}
		log.Info("scrubbed owned OIDC mapping", "mappingUuid", m.MappingUUID, "group", m.GroupName)
	}
	return nil
}

func (r *TeamReconciler) reconcileDelete(ctx context.Context, team *dependencytrackv1alpha1.Team) error {
	log := logf.FromContext(ctx)

	owned := ownedStatusOf(team)
	needScrub := team.Spec.OIDC != nil && len(owned) > 0
	needTeardown := team.Status.UUID != ""

	// Fast path: nothing to scrub and no team object to tear down. Strip the
	// finalizer directly WITHOUT acquiring a DependencyTrack credential, so a
	// vanilla/no-op deletion never forces an unnecessary auth handshake.
	if !needScrub && !needTeardown {
		controllerutil.RemoveFinalizer(team, teamFinalizer)
		return r.Update(ctx, team)
	}

	// Shared acquisition: scrub or team teardown (often both) is required, so
	// we authenticate ONCE and feed the same auth context + API client to
	// both blocks.
	authCtx, apiClient, err := r.DTProvider.Get(ctx)
	if err != nil {
		log.Error(err, "failed to authenticate with DependencyTrack during deletion")
		return err
	}

	// Scrub owned OIDC mapping edges BEFORE destroying the team so the
	// dependent bindings vanish ahead of the principal object. Engaged only
	// when owned mappings exist AND the team opts into OIDC.
	if needScrub {
		api := dependencytrack.NewDtapiOIDCClient(apiClient.OidcAPI)
		if err := r.scrubOwnedOIDCMappings(authCtx, api, team); err != nil {
			if errors.Is(err, dependencytrack.ErrOIDCUnavailable) {
				// OIDC not provisioned: retain last-seen owned mappings (do
				// NOT clear status.OIDC), flush the advisory condition/event,
				// and requeue with the finalizer held so the scrub retries
				// once OIDC becomes available.
				setCondition(team, metav1.ConditionFalse, "OIDCUnavailable",
					fmt.Sprintf("OIDC is not provisioned; owned mappings for team %q could not be scrubbed during deletion", team.Spec.Name))
				r.Recorder.Eventf(team, "Warning", "OIDCUnavailable",
					"OIDC not provisioned; owned mappings for team %q left intact during deletion", team.Spec.Name)
				_ = r.Status().Update(ctx, team)
				return err
			}
			// Non-availability scrub error: structured logs already emitted
			// inside scrubOwnedOIDCMappings. Retain the finalizer and requeue.
			log.Error(err, "failed to scrub owned OIDC mappings during deletion", "team", team.Spec.Name)
			return err
		}
		r.Recorder.Eventf(team, "Normal", "OIDCMappingRemoved",
			"Scrubbed %d owned OIDC mapping(s) for team %q during deletion", len(owned), team.Spec.Name)
	}

	if needTeardown {
		httpResp, err := apiClient.TeamAPI.DeleteTeam(authCtx).Team(dtapi.Team{
			Name: team.Spec.Name,
			Uuid: team.Status.UUID,
		}).Execute()
		if err != nil && (httpResp == nil || httpResp.StatusCode != http.StatusNotFound) {
			log.Error(err, "failed to delete team from DependencyTrack", "uuid", team.Status.UUID)
			return err
		}
		log.Info("deleted team from DependencyTrack", "uuid", team.Status.UUID)
		r.Recorder.Eventf(team, "Normal", "TeamDeleted", "Deleted team %q from DependencyTrack (uuid=%s)", team.Spec.Name, team.Status.UUID)
	}

	controllerutil.RemoveFinalizer(team, teamFinalizer)
	return r.Update(ctx, team)
}

// syncPermissions applies the desired permission set to the team in DependencyTrack.
// If team.Spec.Permissions is nil/empty it does nothing (leave existing permissions alone).
// If it is an empty slice it clears all permissions. If it has values it replaces
// the existing set atomically via SetTeamPermissions.
func (r *TeamReconciler) syncPermissions(
	ctx context.Context,
	apiClient *dtapi.APIClient,
	team *dependencytrackv1alpha1.Team,
	uuid string,
) error {
	log := logf.FromContext(ctx)

	// Nil or absent: leave permissions alone.
	if team.Spec.Permissions == nil {
		return nil
	}

	authCtx, _, err := r.DTProvider.Get(ctx)
	if err != nil {
		return err
	}

	// DependencyTrack 5.0 exposes per-permission POST/DELETE operations. The
	// generated client also contains a newer bulk PUT operation, but that is not
	// supported by all server versions. Compute and apply a delta instead.
	existing, _, err := lookupTeamByUUID(authCtx, apiClient, uuid)
	if err != nil {
		log.Error(err, "failed to read team permissions", "uuid", uuid)
		setCondition(team, metav1.ConditionFalse, "PermissionSyncError", "failed to read current permissions: "+err.Error())
		_ = r.Status().Update(ctx, team)
		return err
	}

	existingPermissions := make([]dtapi.Permission, 0, len(existing.Permissions))
	for _, permission := range existing.Permissions {
		existingPermissions = append(existingPermissions, dtapi.Permission{Name: permission.Name})
	}
	toAdd, toRemove, desired := permissionDelta(existingPermissions, team.Spec.Permissions)
	for _, permission := range toAdd {
		if _, _, err := apiClient.PermissionAPI.AddPermissionToTeam(authCtx, uuid, permission).Execute(); err != nil {
			log.Error(err, "failed to add team permission", "uuid", uuid, "permission", permission)
			setCondition(team, metav1.ConditionFalse, "PermissionSyncError", "failed to add permission "+permission+": "+err.Error())
			_ = r.Status().Update(ctx, team)
			return err
		}
	}
	for _, permission := range toRemove {
		if _, _, err := apiClient.PermissionAPI.RemovePermissionFromTeam(authCtx, uuid, permission).Execute(); err != nil {
			log.Error(err, "failed to remove team permission", "uuid", uuid, "permission", permission)
			setCondition(team, metav1.ConditionFalse, "PermissionSyncError", "failed to remove permission "+permission+": "+err.Error())
			_ = r.Status().Update(ctx, team)
			return err
		}
	}

	// Record the synced permission set (sorted, deduplicated) for observability.
	team.Status.Permissions = joinString(desired, ",")
	log.Info("synced team permissions", "uuid", uuid, "added", toAdd, "removed", toRemove)

	if len(team.Spec.Permissions) == 0 {
		r.Recorder.Eventf(team, "Normal", "PermissionsCleared", "Cleared all permissions for team %q", team.Spec.Name)
	} else {
		r.Recorder.Eventf(team, "Normal", "PermissionsSet", "Set %d permission(s) for team %q", len(team.Spec.Permissions), team.Spec.Name)
	}

	return nil
}

// reconcileOIDC converges the OIDC group mappings backing a team onto
// team.Spec.OIDC, anchoring the diff on team.Status.OIDC (prevOwned) so that
// drops, reordering, and edits all converge rather than regressing to
// additive-only provisioning. It is invoked from reconcileUpsert only when
// team.Spec.OIDC != nil, and additionally self-guards against a nil pointer so
// a forgotten caller guard degrades to a no-op rather than a panic.
//
// Adaptation to the T03 seam ([dependencytrack.dtapiAdapter]) lives at the call
// site: the generated *dtapi.APIClient is bridged via NewDtapiOIDCClient and
// handed to ReconcileOIDCMappings. Tests bypass the adapter entirely by
// supplying a minimal OIDCAClient stub directly, sidestepping the dtapiAdapter
// boundary.
//
// Error handling:
//   - ErrOIDCUnavailable (OIDC not provisioned) is treated as TRANSIENT:
//     team.Status.OIDC is left UNTOUCHED (stale owned mappings are retained
//     conservatively) and the reconcile requeues instead of mutating the
//     bindings the controller does not own.
//   - Any other error faults the Ready condition and bubbles up to the caller.
func (r *TeamReconciler) reconcileOIDC(ctx context.Context, api dependencytrack.OIDCAClient, team *dependencytrackv1alpha1.Team, teamUUID string) (ctrl.Result, error) {
	if team.Spec.OIDC == nil {
		// Baseline teams opt out of OIDC entirely: zero API traffic.
		return ctrl.Result{}, nil
	}

	log := logf.FromContext(ctx)
	prev := ownedStatusToSeam(ownedStatusOf(team))
	mappings, err := dependencytrack.ReconcileOIDCMappings(ctx, logrLogger{l: log}, api, teamUUID, team.Spec.OIDC.Groups, prev)
	if err != nil {
		switch {
		case errors.Is(err, dependencytrack.ErrOIDCUnavailable):
			// Feature not provisioned: do NOT reassign team.Status.OIDC so the
			// last-seen owned mappings are retained conservatively. Flush the
			// advisory condition/event and requeue for retry.
			setCondition(team, metav1.ConditionFalse, "OIDCUnavailable",
				fmt.Sprintf("OIDC is not provisioned on the DependencyTrack instance; mappings for team %q left untouched", team.Spec.Name))
			r.Recorder.Eventf(team, "Warning", "OIDCUnavailable",
				"OIDC not provisioned; existing mappings left untouched for team %q", team.Spec.Name)
			_ = r.Status().Update(ctx, team)
			return ctrl.Result{Requeue: true}, err
		default:
			setCondition(team, metav1.ConditionFalse, "OIDCError", "failed to reconcile OIDC mappings: "+err.Error())
			r.Recorder.Eventf(team, "Warning", "OIDCError",
				"failed to reconcile OIDC mappings for team %q: %s", team.Spec.Name, err.Error())
			_ = r.Status().Update(ctx, team)
			return ctrl.Result{}, err
		}
	}

	team.Status.OIDC = &dependencytrackv1alpha1.TeamOIDCStatus{
		OwnedMappings: ownedMappingsToStatus(mappings),
	}

	if len(mappings) > 0 {
		r.Recorder.Eventf(team, "Normal", "OIDCMappingCreated",
			"Ensured %d OIDC group mapping(s) for team %q", len(mappings), team.Spec.Name)
	} else {
		r.Recorder.Event(team, "Normal", "OIDCMappingSkipped",
			fmt.Sprintf("No OIDC mappings to create for team %q", team.Spec.Name))
	}

	return ctrl.Result{}, nil
}

// ownedMappingsToStatus translates seam-owned bindings ([dependencytrack.OwnedMapping])
// into the API-status DTOs exposed on Team.status.oidc. Kept as a pure helper so
// the field-map contract is unit-tested independently of the generated client.
func ownedMappingsToStatus(seam []dependencytrack.OwnedMapping) []dependencytrackv1alpha1.OwnedOIDCMapping {
	out := make([]dependencytrackv1alpha1.OwnedOIDCMapping, 0, len(seam))
	for _, m := range seam {
		out = append(out, dependencytrackv1alpha1.OwnedOIDCMapping{
			GroupName:   m.GroupName,
			GroupUUID:   m.GroupUUID,
			TeamUUID:    m.TeamUUID,
			MappingUUID: m.MappingUUID,
		})
	}
	return out
}

// ownedStatusOf plucks the owned-mapping rows currently held in a Team's status
// for feeding back into the diff engine. Returns nil when status.OIDC is unset
// so a never-reconciled team behaves identically to "previously owned nothing".
func ownedStatusOf(team *dependencytrackv1alpha1.Team) []dependencytrackv1alpha1.OwnedOIDCMapping {
	if team == nil || team.Status.OIDC == nil {
		return nil
	}
	return team.Status.OIDC.OwnedMappings
}

// ownedStatusToSeam is the inverse of [ownedMappingsToStatus]: it materialises
// the status-recorded bindings into the seam-owned tuples that
// [dependencytrack.ReconcileOIDCMappings] consumes as its prevOwned anchor.
// Empty/nil input yields nil so the diff engine observes "nothing previously
// owned" rather than an empty-but-present ownership set.
func ownedStatusToSeam(status []dependencytrackv1alpha1.OwnedOIDCMapping) []dependencytrack.OwnedMapping {
	if len(status) == 0 {
		return nil
	}
	out := make([]dependencytrack.OwnedMapping, 0, len(status))
	for _, m := range status {
		out = append(out, dependencytrack.OwnedMapping{
			GroupName:   m.GroupName,
			GroupUUID:   m.GroupUUID,
			TeamUUID:    m.TeamUUID,
			MappingUUID: m.MappingUUID,
		})
	}
	return out
}

// logrLogger adapts a [logr.Logger] to the [dependencytrack.Logger] interface so
// the OIDC seam participates in the controller's structured log stream.

type logrLogger struct{ l logr.Logger }

func (lg logrLogger) Printf(format string, v ...any) {
	lg.l.Info(fmt.Sprintf(format, v...))
}

// failStatus sets a failed condition, persists the status, and returns the error so the reconcile loop requeues.
func (r *TeamReconciler) failStatus(ctx context.Context, team *dependencytrackv1alpha1.Team, msg string, cause error) (ctrl.Result, error) {
	logf.FromContext(ctx).Error(cause, msg)
	setCondition(team, metav1.ConditionFalse, reasonAPIError, msg)
	_ = r.Status().Update(ctx, team)
	return ctrl.Result{}, cause
}

// setCondition upserts the Ready condition on the team's status.
func setCondition(team *dependencytrackv1alpha1.Team, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&team.Status.Conditions, metav1.Condition{
		Type:               conditionReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: team.Generation,
	})
}

// permissionDelta returns the sorted permissions to add and remove, plus a
// canonical sorted and deduplicated snapshot of the desired permissions.
func permissionDelta(current []dtapi.Permission, desired []string) (toAdd, toRemove, canonical []string) {
	currentSet := make(map[string]struct{}, len(current))
	for _, permission := range current {
		currentSet[permission.Name] = struct{}{}
	}

	desiredSet := make(map[string]struct{}, len(desired))
	for _, permission := range desired {
		desiredSet[permission] = struct{}{}
	}

	for permission := range desiredSet {
		canonical = append(canonical, permission)
		if _, exists := currentSet[permission]; !exists {
			toAdd = append(toAdd, permission)
		}
	}
	for permission := range currentSet {
		if _, exists := desiredSet[permission]; !exists {
			toRemove = append(toRemove, permission)
		}
	}

	slices.Sort(toAdd)
	slices.Sort(toRemove)
	slices.Sort(canonical)
	return toAdd, toRemove, canonical
}

// joinString joins a sorted string slice with "," separators.
func joinString(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for _, s := range ss[1:] {
		result += sep + s
	}
	return result
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dependencytrackv1alpha1.Team{}).
		Named("team").
		Complete(r)
}
