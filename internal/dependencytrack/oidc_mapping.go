// Copyright 2026.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dependencytrack

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

// Logger is the minimal subset of the standard library *log.Logger that the
// OIDC reconciliation helpers depend on. Keeping it to Printf allows callers
// to supply a *log.Logger, a logr adapter wrapper, or a test sink.
type Logger interface {
	Printf(format string, v ...any)
}

// OwnedMapping captures the binding established between an OIDC group and a
// Dependency-Track team during a reconcile. GroupName is the human-readable
// claim value requested by the operator; the UUID fields are the server-issued
// identifiers returned by Dependency-Track.
type OwnedMapping struct {
	GroupName   string
	GroupUUID   string
	TeamUUID    string
	MappingUUID string
}

// ErrOIDCUnavailable is returned by ReconcileOIDCMappings when the Dependency-Track
// instance reports that OpenID Connect is not provisioned. It is a sentinel so
// callers can distinguish "feature off" from transient transport failures.
var ErrOIDCUnavailable = errors.New("openid connect is not available on the DependencyTrack instance")

// ErrGroupAlreadyExists is returned by [OIDCAClient.CreateGroup] (surfaced
// through [dtapiAdapter.CreateGroup] on HTTP 409 Conflict) when a group bearing
// the requested name already exists on the Dependency-Track instance. Rather
// than aborting the reconcile, callers use it to drive a conservative
// relist-and-relink: re-fetching the group roster and adopting the existing
// group by EXACT-CASE name so a racing peer's victory converges idempotently.
var ErrGroupAlreadyExists = errors.New("oidc group already exists")

// OIDCAClient is the seams-facing abstraction over the generated dtapi
// OidcAPIService. Signatures deliberately mirror the generated client so the
// production adapter ([dtapiAdapter]) is mechanical and tests can substitute
// [FakeOIDCAPI] without standing up a real Dependency-Track instance.
type OIDCAClient interface {
	// IsAvailable reports whether OIDC is provisioned on the instance.
	IsAvailable(ctx context.Context) (bool, error)
	// RetrieveGroups lists every OIDC group known to the instance.
	RetrieveGroups(ctx context.Context) ([]dtapi.OidcGroup, error)
	// CreateGroup creates a group identified solely by name and returns the
	// server-populated group (including its issued UUID). Returns
	// [ErrGroupAlreadyExists] when the instance already holds a group with
	// that exact name (HTTP 409).
	CreateGroup(ctx context.Context, name string) (dtapi.OidcGroup, error)
	// AddMapping2 binds an existing group to a team and returns the resulting
	// mapping (including its issued UUID). Corresponds to PUT /v1/oidc/mapping.
	AddMapping2(ctx context.Context, groupUUID, teamUUID string) (dtapi.MappedOidcGroup, error)
	// DeleteMappingByUuid removes the OIDC mapping identified by its UUID.
	// Implementations MUST treat an already-absent mapping (HTTP 404) as success
	// so that pruning is idempotent across concurrent/repeated reconciles.
	DeleteMappingByUuid(ctx context.Context, mappingUUID string) error
}

// dtapiAdapter adapts []*dtapi.OidcAPIService to [OIDCAClient]. It inspects the
// *http.Response emitted by the generated client where status codes matter
// (notably 409 on CreateGroup and 404 on DeleteMappingByUuid) and unwraps the
// pointer-valued results into value-typed DTOs so consumers deal exclusively
// with owned copies.
type dtapiAdapter struct {
	svc *dtapi.OidcAPIService
}

// NewDtapiOIDCClient wraps a generated OIDC service handle so it can be fed to
// [ReconcileOIDCMappings] in place of a fake.
func NewDtapiOIDCClient(svc *dtapi.OidcAPIService) OIDCAClient {
	return &dtapiAdapter{svc: svc}
}

// Compile-time assertion that dtapiAdapter satisfies OIDCAClient.
var _ OIDCAClient = (*dtapiAdapter)(nil)

func (a *dtapiAdapter) IsAvailable(ctx context.Context) (bool, error) {
	ok, httpResp, err := a.svc.IsAvailable(ctx).Execute()
	if err == nil {
		return ok, nil
	}
	// The generated dtapi client (built against an older Dependency-Track
	// OpenAPI spec) decodes /v1/oidc/available as text/plain into a bool, but
	// Dependency-Track 5.x replies with Content-Type
	// "text/plain;charset=iso-8859-1", which the generated decoder rejects as
	// "undefined response type" -- even on HTTP 200. The generator nonetheless
	// rescues the response body onto the returned *http.Response, so leniently
	// parse the already-fetched body instead of propagating the spurious decode
	// error. Falls back to the wrapped error only if the body cannot be read
	// or does not match the expected boolean text.
	if httpResp != nil {
		if body, readErr := io.ReadAll(httpResp.Body); readErr == nil {
			switch strings.ToLower(strings.TrimSpace(string(body))) {
			case "true":
				return true, nil
			case "false":
				return false, nil
			}
		}
	}
	return false, fmt.Errorf("dtapi IsAvailable: %w", err)
}

func (a *dtapiAdapter) RetrieveGroups(ctx context.Context) ([]dtapi.OidcGroup, error) {
	groups, _, err := a.svc.RetrieveGroups(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("dtapi RetrieveGroups: %w", err)
	}
	return groups, nil
}

func (a *dtapiAdapter) CreateGroup(ctx context.Context, name string) (dtapi.OidcGroup, error) {
	resp, httpResp, err := a.svc.CreateGroup(ctx).OidcGroup(dtapi.OidcGroup{Name: name}).Execute()
	if err != nil {
		// The generated client populates *http.Response even on >=300
		// statuses, so a 409 (Dependency-Track rejecting a duplicate name)
		// is translated to the sentinel that drives relist-and-relink.
		if httpResp != nil && httpResp.StatusCode == http.StatusConflict {
			return dtapi.OidcGroup{}, fmt.Errorf("dtapi CreateGroup(%q): %w", name, ErrGroupAlreadyExists)
		}
		return dtapi.OidcGroup{}, fmt.Errorf("dtapi CreateGroup(%q): %w", name, err)
	}
	if resp == nil {
		return dtapi.OidcGroup{}, nil
	}
	return *resp, nil
}

func (a *dtapiAdapter) AddMapping2(ctx context.Context, groupUUID, teamUUID string) (dtapi.MappedOidcGroup, error) {
	resp, _, err := a.svc.AddMapping2(ctx).
		MappedOidcGroupRequest(dtapi.MappedOidcGroupRequest{Group: groupUUID, Team: teamUUID}).
		Execute()
	if err != nil {
		return dtapi.MappedOidcGroup{}, fmt.Errorf("dtapi AddMapping2(group=%q,team=%q): %w", groupUUID, teamUUID, err)
	}
	if resp == nil {
		return dtapi.MappedOidcGroup{}, nil
	}
	return *resp, nil
}

// DeleteMappingByUuid prunes an individual OIDC mapping. The generated
// DeleteMappingByUuid op returns the *http.Response even on a >=300 status, so
// a 404 (mapping already reclaimed by a racing/previous reconcile, or removed
// manually) is detected and swallowed to keep the prune convergent and
// idempotent. Any other error is propagated wrapped.
func (a *dtapiAdapter) DeleteMappingByUuid(ctx context.Context, mappingUUID string) error {
	resp, err := a.svc.DeleteMappingByUuid(ctx, mappingUUID).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("dtapi DeleteMappingByUuid(%q): %w", mappingUUID, err)
	}
	return nil
}

// ReconcileOIDCMappings reconciles the OIDC mappings backing a team against a
// declared desired set, anchoring convergence on the bindings previously
// recorded in status (prevOwned) so that removals, edits, and reorderings all
// converge safely.
//
// Algorithm:
//  1. De-duplicate the desired list (order-preserving, blanks dropped) and fold
//     prevOwned by GroupName (first occurrence wins) so duplicate status rows
//     cannot fan out into redundant mutators.
//  2. Bail out with zero API traffic when NEITHER desired NOR prevOwned carry
//     anything to reconcile.
//  3. Probe [OIDCAClient.IsAvailable]; abort with [ErrOIDCUnavailable] and
//     perform NO mutation when OIDC is not provisioned. The caller retains
//     prevOwned and simply retries later.
//  4. Partition prevOwned into survivors (GroupName still desired) and stale
//     (GroupName no longer desired). Each stale edge is pruned exactly once via
//     [OIDCAClient.DeleteMappingByUuid], which treats an already-absent mapping
//     (404) as success.
//  5. Compute toEstablish = desired names not among the survivors, preserving
//     the desired declaration order. Issue a SINGLE batched
//     [OIDCAClient.RetrieveGroups] (only when something remains to establish)
//     and index it by EXACT-CASE name.
//  6. For each to-establish name: lookup-or-[OIDCAClient.CreateGroup] then
//     [OIDCAClient.AddMapping2]. If CreateGroup yields [ErrGroupAlreadyExists]
//     (a racing peer won the create), re-issue [OIDCAClient.RetrieveGroups],
//     adopt the winner by EXACT-CASE name, and resume to AddMapping2; any other
//     CreateGroup error still aborts. Append the claimed binding to the result.
//  7. Return survivors concatenated with newly claimed bindings.
//
// Reordering the desired list never touches the server: a reordered-but-equal
// set leaves every previous owner surviving, empties toEstablish, and therefore
// invokes no mutators. Group OBJECTS are never deleted — only mapping EDGES
// (by UUID) are pruned, so shared/existing groups survive a team shrink.
func ReconcileOIDCMappings(ctx context.Context, logger Logger, api OIDCAClient, teamUUID string, desiredGroups []string, prevOwned []OwnedMapping) ([]OwnedMapping, error) {
	desired := uniqueNonEmpty(desiredGroups)
	prevClean := collapsePrevOwned(prevOwned)

	if len(desired) == 0 && len(prevClean) == 0 {
		logger.Printf("reconcile OIDC mappings: no desired groups and no previous ownerships for team %q; skipping", teamUUID)
		return nil, nil
	}

	available, err := api.IsAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("probing OIDC availability for team %q: %w", teamUUID, err)
	}
	if !available {
		logger.Printf("reconcile OIDC mappings: OIDC not available; aborting for team %q (previous ownerships left untouched)", teamUUID)
		return nil, ErrOIDCUnavailable
	}

	desiredSet := make(map[string]struct{}, len(desired))
	for _, d := range desired {
		desiredSet[d] = struct{}{}
	}

	// =========================================================================
	// LOCK-COMMENT -- Conservative-Restore Invariant (status-loss safety)
	// -------------------------------------------------------------------------
	// When prevOwned is nil/empty, status has been wiped (operator/upgrade
	// status-subresource reset) or was never persisted. Under wipe the reconcile
	// MUST enter RESTORE-ONLY MODE:
	//
	//   * Ensure EVERY desired binding exists -- reuse-via-list or create-bind.
	//   * Emit ZERO DeleteMappingByUuid -- there is no trusted prior ownership
	//     to difference against, so pruning would erase bindings authored by
	//     peers/the platform itself, violating tenant isolation under concurrency.
	//
	// STRUCTURAL ENFORCEMENT: `stale` is derived EXCLUSIVELY from `prevClean`
	// (collapsed prevOwned). When prevOwned is nil/empty, prevClean is empty,
	// hence `stale` is empty, hence the prune loop below is PROVABLY DEAD and
	// DeleteMappingByUuid is UNREACHABLE. Deletion can never originate from
	// server state or the desired set -- only from acknowledged prior ownership.
	// Diverting `stale` to consult the server/eDesired would VIOLATE this
	// invariant and is forbidden.
	//
	// Restore-only therefore rides the SAME exact-case relink path as normal
	// reconciliation: a racing peer's group (even under StrictConflicts / HTTP
	// 409) is recovered via relist-and-adopt-by-EXACT name and rebound, never
	// clobbered. Pinned by TestReconcileOIDC_Wipe* in oidc_mapping_test.go.
	// =========================================================================

	// Step 4: partition prevOwned into survivors (still desired) and stale.
	survivors := make([]OwnedMapping, 0, len(prevClean))
	var stale []OwnedMapping
	survivorSeen := make(map[string]struct{}, len(prevClean))
	for _, om := range prevClean {
		if _, wanted := desiredSet[om.GroupName]; !wanted {
			stale = append(stale, om)
			continue
		}
		// collapsePrevOwned already deduplicated, but defend against future
		// regressions: never emit a duplicate survivor slot.
		if _, dup := survivorSeen[om.GroupName]; dup {
			continue
		}
		survivorSeen[om.GroupName] = struct{}{}
		survivors = append(survivors, om)
	}

	// Prune stale edges idempotently.
	for _, om := range stale {
		delErr := api.DeleteMappingByUuid(ctx, om.MappingUUID)
		if delErr != nil {
			return nil, fmt.Errorf("pruning stale OIDC mapping %q for group %q: %w", om.MappingUUID, om.GroupName, delErr)
		}
		logger.Printf("reconcile OIDC mappings: pruned stale mapping %q (group %q)", om.MappingUUID, om.GroupName)
	}

	// Step 5: names that still need establishing = desired minus survivors.
	toEstablish := make([]string, 0, len(desired))
	for _, d := range desired {
		if _, surv := survivorSeen[d]; surv {
			continue
		}
		toEstablish = append(toEstablish, d)
	}

	newlyClaimed := make([]OwnedMapping, 0, len(toEstablish))
	if len(toEstablish) > 0 {
		existing, err := api.RetrieveGroups(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing existing OIDC groups for team %q: %w", teamUUID, err)
		}
		// Indexed by exact-case name: deliberate NOT lower-cased so a casing
		// typo surfaces as a distinct group rather than silently colliding.
		indexed := make(map[string]dtapi.OidcGroup, len(existing))
		for i := range existing {
			indexed[existing[i].Name] = existing[i]
		}

		for _, name := range toEstablish {
			group, existed := indexed[name]
			if !existed {
				created, cerr := api.CreateGroup(ctx, name)
				if cerr != nil {
					if errors.Is(cerr, ErrGroupAlreadyExists) {
						// Lost the create race. Authoritatively re-fetch the
						// roster and adopt the winning group by EXACT-CASE
						// name so a case-distinct sibling (e.g. "Admin" vs
						// "ADMIN") is never mistaken for ours.
						refreshed, rerr := api.RetrieveGroups(ctx)
						if rerr != nil {
							return nil, fmt.Errorf("re-listing OIDC groups after conflict for %q: %w", name, rerr)
						}
						adopted := false
						for ri := range refreshed {
							if refreshed[ri].Name == name {
								group = refreshed[ri]
								adopted = true
								break
							}
						}
						if !adopted {
							return nil, fmt.Errorf("OIDC group %q vanished between 409 and relist", name)
						}
						logger.Printf("reconcile OIDC mappings: create of %q lost to a race; adopted existing group %s via relist", name, group.Uuid)
					} else {
						return nil, fmt.Errorf("creating OIDC group %q: %w", name, cerr)
					}
				} else {
					group = created
					logger.Printf("reconcile OIDC mappings: created group %q (%s)", name, created.Uuid)
				}
			} else {
				logger.Printf("reconcile OIDC mappings: reusing existing group %q (%s)", name, group.Uuid)
			}

			mapped, merr := api.AddMapping2(ctx, group.Uuid, teamUUID)
			if merr != nil {
				return nil, fmt.Errorf("mapping OIDC group %q to team %q: %w", name, teamUUID, merr)
			}
			logger.Printf("reconcile OIDC mappings: mapped group %q (%s) to team %q as mapping %s", name, group.Uuid, teamUUID, mapped.Uuid)

			newlyClaimed = append(newlyClaimed, OwnedMapping{
				GroupName:   name,
				GroupUUID:   group.Uuid,
				TeamUUID:    teamUUID,
				MappingUUID: mapped.Uuid,
			})
		}
	}

	// Step 7: survivors (verbatim, in prevOwned order) ++ newly claimed.
	out := make([]OwnedMapping, 0, len(survivors)+len(newlyClaimed))
	out = append(out, survivors...)
	out = append(out, newlyClaimed...)
	return out, nil
}

// uniqueNonEmpty collapses consecutive/non-consecutive duplicates while
// dropping blank entries (a blank group name is never a legal Identifier
// Provider claim and would only confuse the server). Insertion order is kept.
func uniqueNonEmpty(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// collapsePrevOwned folds duplicate prevOwned entries by GroupName, retaining
// the FIRST occurrence and dropping blank-named ghosts. Duplicate status rows
// (possible during a crash-midway upgrade of the CRD status) must never fan
// out into multiple DeleteMapping or AddMapping invocations.
func collapsePrevOwned(in []OwnedMapping) []OwnedMapping {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]OwnedMapping, 0, len(in))
	for _, om := range in {
		if om.GroupName == "" {
			continue
		}
		if _, ok := seen[om.GroupName]; ok {
			continue
		}
		seen[om.GroupName] = struct{}{}
		out = append(out, om)
	}
	return out
}
