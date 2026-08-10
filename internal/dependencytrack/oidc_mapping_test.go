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
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

// Compile-time proof that the expanded OIDCAClient interface (now including
// DeleteMappingByUuid) is satisfied by both the production adapter and the
// test double.
var (
	_ OIDCAClient = (*dtapiAdapter)(nil)
	_ OIDCAClient = (*FakeOIDCAPI)(nil)
)

// ---------------------------------------------------------------------------
// FakeOIDCAPI
// ---------------------------------------------------------------------------

// FakeMappingRequest is the payload handed to AddMapping2, exposed for
// assertions in tests.
type FakeMappingRequest struct {
	GroupUUID string
	TeamUUID  string
}

// FakeOIDCAPI is an in-process, mutex-guarded stand-in for [OIDCAClient]. It
// records every mutating call and hands back sequentially-allocated UUIDs so
// tests can assert both the shape of interactions and the provenance of every
// identifier flowing into the returned [OwnedMapping] slice.
type FakeOIDCAPI struct {
	mu             sync.Mutex
	Available      bool
	ExistingGroups []dtapi.OidcGroup
	nextGroupID    int
	nextMappingID  int

	CreateGroupCalls    int
	AddMappingCalls     int
	CreateGroupRequests []string
	AddMappingRequests  []FakeMappingRequest

	DeleteMappingByUuidCalls int
	DeletedMappingUUIDs      []string

	ErrIsAvailable    error
	ErrRetrieveGroups error
	ErrCreateGroup    error
	ErrAddMapping2    error
	ErrDeleteMapping  error

	// StrictConflicts mirrors Dependency-Track's implicit UNIQUE constraint on
	// OIDC group names: when set, CreateGroup returns [ErrGroupAlreadyExists]
	// if a group bearing the requested name already lives in ExistingGroups.
	// Left false (lenient) for backward compatibility with earlier tests.
	StrictConflicts bool

	// RacingCreators simulates a competing reconciler that wins the create
	// race between our RetrieveGroups and CreateGroup. Each group is injected
	// into ExistingGroups at the START of the next CreateGroup call (then
	// cleared), so a StrictConflicts violation is observed exactly once --
	// mirroring a transient HTTP 409 from the server.
	RacingCreators []dtapi.OidcGroup
}

func (f *FakeOIDCAPI) IsAvailable(_ context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.ErrIsAvailable != nil {
		return false, f.ErrIsAvailable
	}
	return f.Available, nil
}

func (f *FakeOIDCAPI) RetrieveGroups(_ context.Context) ([]dtapi.OidcGroup, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.ErrRetrieveGroups != nil {
		return nil, f.ErrRetrieveGroups
	}
	cp := make([]dtapi.OidcGroup, len(f.ExistingGroups))
	copy(cp, f.ExistingGroups)
	return cp, nil
}

func (f *FakeOIDCAPI) CreateGroup(_ context.Context, name string) (dtapi.OidcGroup, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.CreateGroupCalls++
	f.CreateGroupRequests = append(f.CreateGroupRequests, name)

	// Consume any queued racing-creator insertions before the existence check,
	// simulating a competing writer that won the create race between our
	// RetrieveGroups and this CreateGroup. Cleared afterwards so the effect
	// is observed exactly once.
	if len(f.RacingCreators) > 0 {
		f.ExistingGroups = append(f.ExistingGroups, f.RacingCreators...)
		f.RacingCreators = nil
	}

	if f.StrictConflicts {
		for _, g := range f.ExistingGroups {
			if g.Name == name {
				return dtapi.OidcGroup{}, ErrGroupAlreadyExists
			}
		}
	}

	if f.ErrCreateGroup != nil {
		return dtapi.OidcGroup{}, f.ErrCreateGroup
	}
	g := dtapi.OidcGroup{Uuid: f.nextUUID(&f.nextGroupID), Name: name}
	// Mirror a stateful server: the freshly created group becomes discoverable
	// for any subsequent RetrieveGroups within the same fake lifetime.
	f.ExistingGroups = append(f.ExistingGroups, g)
	return g, nil
}

func (f *FakeOIDCAPI) AddMapping2(_ context.Context, groupUUID, teamUUID string) (dtapi.MappedOidcGroup, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.AddMappingCalls++
	f.AddMappingRequests = append(f.AddMappingRequests, FakeMappingRequest{GroupUUID: groupUUID, TeamUUID: teamUUID})
	if f.ErrAddMapping2 != nil {
		return dtapi.MappedOidcGroup{}, f.ErrAddMapping2
	}
	return dtapi.MappedOidcGroup{Uuid: f.nextUUID(&f.nextMappingID)}, nil
}

// DeleteMappingByUuid records the prune and returns ErrDeleteMapping when wired.
// Like the production adapter, the fake is 404-tolerant by construction: it
// never raises "not found" for a UUID it hasn't seen, modelling the
// idempotent-converge contract exercised by ReconcileOIDCMappings.
func (f *FakeOIDCAPI) DeleteMappingByUuid(_ context.Context, mappingUUID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.DeleteMappingByUuidCalls++
	f.DeletedMappingUUIDs = append(f.DeletedMappingUUIDs, mappingUUID)
	if f.ErrDeleteMapping != nil {
		return f.ErrDeleteMapping
	}
	return nil
}

// nextUUID emits monotonically increasing RFC-4122-shaped identifiers from a
// zero-initialised counter. Format: 00000000-0000-0000-0000-<12 hex>.
func (f *FakeOIDCAPI) nextUUID(counter *int) string {
	n := *counter
	if n == 0 {
		n = 1
	}
	*counter = n + 1
	return fmt.Sprintf("00000000-0000-0000-0000-%012x", n)
}

// ---------------------------------------------------------------------------
// recordingLogger
// ---------------------------------------------------------------------------

type recordingLogger struct {
	lines []string
}

func (l *recordingLogger) Printf(format string, v ...any) {
	l.lines = append(l.lines, fmt.Sprintf(format, v...))
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newFake() *FakeOIDCAPI {
	return &FakeOIDCAPI{Available: true}
}

func mustSingle(t *testing.T, fn func(context.Context, Logger, OIDCAClient, string, []string, []OwnedMapping) ([]OwnedMapping, error),
	log *recordingLogger, api OIDCAClient, team string, desired []string, prevOwned []OwnedMapping) []OwnedMapping {
	t.Helper()
	res, err := fn(context.Background(), log, api, team, desired, prevOwned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return res
}

// ---------------------------------------------------------------------------
// Positive path: two brand-new groups => 2 CreateGroup + 2 AddMapping
// ---------------------------------------------------------------------------

func TestReconcileOIDC_CreatesMissingGroupsAndMaps(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"alpha", "beta"}, nil)

	if api.CreateGroupCalls != 2 {
		t.Errorf("CreateGroup calls = %d, want 2", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 2 {
		t.Errorf("AddMapping calls = %d, want 2", api.AddMappingCalls)
	}
	if len(api.CreateGroupRequests) != 2 || api.CreateGroupRequests[0] != "alpha" || api.CreateGroupRequests[1] != "beta" {
		t.Errorf("CreateGroupRequests = %v", api.CreateGroupRequests)
	}
	if len(result) != 2 {
		t.Fatalf("OwnedList len = %d, want 2", len(result))
	}

	wantAlpha := OwnedMapping{GroupName: "alpha", TeamUUID: team, MappingUUID: "00000000-0000-0000-0000-000000000001"}
	gotAlpha := result[0]
	if gotAlpha.GroupName != wantAlpha.GroupName || gotAlpha.TeamUUID != wantAlpha.TeamUUID || gotAlpha.MappingUUID != wantAlpha.MappingUUID {
		t.Errorf("alpha mapping = %+v, want GroupName=%q TeamUUID=%q MappingUUID=%q",
			gotAlpha, wantAlpha.GroupName, wantAlpha.TeamUUID, wantAlpha.MappingUUID)
	}
	if gotAlpha.GroupUUID == "" || gotAlpha.MappingUUID == "" {
		t.Errorf("alpha mapping has empty UUID(s): %+v", gotAlpha)
	}
	if result[0].MappingUUID == result[1].MappingUUID {
		t.Errorf("mapping UUIDs collided: alpha=%q beta=%q", result[0].MappingUUID, result[1].MappingUUID)
	}
}

// ---------------------------------------------------------------------------
// Empty input => zero API traffic, nil slice
// ---------------------------------------------------------------------------

func TestReconcileOIDC_EmptyInputMakesNoCalls(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
	if result != nil {
		t.Errorf("expected nil slice for empty input, got %v", result)
	}
}

// Blank-string entries collapse identically to empty input.
func TestReconcileOIDC_BlankEntriesMakeNoCalls(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, []string{"", "", ""}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
	if result != nil {
		t.Errorf("expected nil slice for blank-only input, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Idempotent reuse: existing groups are looked up, never re-created
// ---------------------------------------------------------------------------

func TestReconcileOIDC_ReusesExistingGroups(t *testing.T) {
	existing := []dtapi.OidcGroup{
		{Uuid: "grp-alpha-uuid", Name: "alpha"},
	}
	api := &FakeOIDCAPI{Available: true, ExistingGroups: existing}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"alpha", "beta"}, nil)

	if api.CreateGroupCalls != 1 {
		t.Errorf("CreateGroup calls = %d, want 1 (only beta should be created)", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 2 {
		t.Errorf("AddMapping calls = %d, want 2", api.AddMappingCalls)
	}
	if len(result) != 2 {
		t.Fatalf("OwnedList len = %d, want 2", len(result))
	}
	if result[0].GroupUUID != "grp-alpha-uuid" {
		t.Errorf("reused alpha GroupUUID = %q, want grp-alpha-uuid", result[0].GroupUUID)
	}
	if result[0].GroupName != "alpha" {
		t.Errorf("alpha GroupName = %q, want alpha", result[0].GroupName)
	}
}

// Case sensitivity: an existing group titled differently cased must NOT satisfy
// a desired name that differs in case.
func TestReconcileOIDC_LookupIsCaseSensitive(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:      true,
		ExistingGroups: []dtapi.OidcGroup{{Uuid: "u-exact", Name: "Admin"}},
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"ADMIN"}, nil)

	if api.CreateGroupCalls != 1 {
		t.Errorf("CreateGroup calls = %d, want 1 (case mismatch should force creation)", api.CreateGroupCalls)
	}
	if result[0].GroupName != "ADMIN" {
		t.Errorf("GroupName = %q, want ADMIN", result[0].GroupName)
	}
}

// race_lost_relinks_via_exact_case_pass: a racing peer creates the group between
// our RetrieveGroups and CreateGroup. Our CreateGroup observes the STRICT
// conflict (HTTP 409 analogue), the reconcile re-lists, adopts the winner by
// EXACT-CASE name, and resumes to AddMapping2 -- converging idempotently in-pass
// instead of aborting for a blind full-retry.
func TestReconcileOIDC_RaceLostRelinksViaExactCasePass(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:       true,
		StrictConflicts: true,
		// The rival reconciler slipped this group in after our initial
		// RetrieveGroups returned empty.
		RacingCreators: []dtapi.OidcGroup{
			{Uuid: "racer-winner-uuid", Name: "admin"},
		},
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"admin"}, nil)

	// Exactly one CreateGroup attempted; it lost the race. The reconcile then
	// re-listed, adopted the winner, and bound the mapping.
	if api.CreateGroupCalls != 1 {
		t.Errorf("CreateGroup calls = %d, want 1", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 1 {
		t.Errorf("AddMapping calls = %d, want 1", api.AddMappingCalls)
	}
	if len(result) != 1 {
		t.Fatalf("result len = %d, want 1", len(result))
	}
	if result[0].GroupName != "admin" {
		t.Errorf("GroupName = %q, want admin", result[0].GroupName)
	}
	if result[0].GroupUUID != "racer-winner-uuid" {
		t.Errorf("GroupUUID = %q, want racer-winner-uuid (lost-create must adopt the winner)", result[0].GroupUUID)
	}
	if result[0].TeamUUID != team {
		t.Errorf("TeamUUID = %q, want %q", result[0].TeamUUID, team)
	}
}

// exact_case_relink_does_not_alias_twin: an existing group titled "Admin" must
// NOT satisfy a desired name "ADMIN" -- the two are distinct claims. The
// latter is materialised as its OWN group, never aliased to the twin.
func TestReconcileOIDC_ExactCaseRelinkDoesNotAliasTwin(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:      true,
		ExistingGroups: []dtapi.OidcGroup{{Uuid: "twin-admin-uuid", Name: "Admin"}},
		// StrictConflicts left false (lenient default): even so, the exact-case
		// index guarantees "ADMIN" never aliases the "Admin" twin.
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"ADMIN"}, nil)

	// "Admin" exists but "ADMIN" != "Admin" (exact-case), so "ADMIN" is created
	// as a NEW group, never aliased to the twin.
	if api.CreateGroupCalls != 1 {
		t.Errorf("CreateGroup calls = %d, want 1 (ADMIN != Admin, must create)", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 1 {
		t.Errorf("AddMapping calls = %d, want 1", api.AddMappingCalls)
	}
	if len(result) != 1 {
		t.Fatalf("result len = %d, want 1", len(result))
	}
	if result[0].GroupName != "ADMIN" {
		t.Errorf("GroupName = %q, want ADMIN", result[0].GroupName)
	}
	if result[0].GroupUUID == "twin-admin-uuid" {
		t.Errorf("must not alias the Admin twin; GroupUUID = %q", result[0].GroupUUID)
	}
}

// ---------------------------------------------------------------------------
// Deduplication: repeated desired names never multiply calls or entries
// ---------------------------------------------------------------------------

func TestReconcileOIDC_DeduplicatesRepeatedDesiredNames(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"dup", "dup", "uniq"}, nil)

	if api.CreateGroupCalls != 2 {
		t.Errorf("CreateGroup calls = %d, want 2 (dedup)", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 2 {
		t.Errorf("AddMapping calls = %d, want 2 (dedup)", api.AddMappingCalls)
	}
	if len(result) != 2 {
		t.Fatalf("OwnedList len = %d, want 2 after dedup", len(result))
	}
}

// ---------------------------------------------------------------------------
// Diff-aware behaviour driven by prevOwned
// ---------------------------------------------------------------------------

// stale_pruned_exactly_once: mappings whose group fell out of the desired set
// are pruned by UUID, exactly once each; the live edge is untouched.
func TestReconcileOIDC_StalePrunedExactlyOnce(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "alpha", GroupUUID: "ga", TeamUUID: team, MappingUUID: "ma-keep"},
		{GroupName: "beta", GroupUUID: "gb", TeamUUID: team, MappingUUID: "mb-prune"},
		{GroupName: "gamma", GroupUUID: "gc", TeamUUID: team, MappingUUID: "mc-prune"},
	}

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"alpha"}, prevOwned)

	// Only the two stale edges are pruned, each exactly once; the live alpha
	// edge is left untouched.
	if api.DeleteMappingByUuidCalls != 2 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 2", api.DeleteMappingByUuidCalls)
	}
	calls := make(map[string]int, len(api.DeletedMappingUUIDs))
	for _, u := range api.DeletedMappingUUIDs {
		calls[u]++
	}
	for _, want := range []string{"mb-prune", "mc-prune"} {
		if calls[want] != 1 {
			t.Errorf("DeleteMappingByUuid(%q) call count = %d, want 1", want, calls[want])
		}
	}
	if calls["ma-keep"] != 0 {
		t.Errorf("live alpha mapping ma-keep must not be pruned (calls=%d)", calls["ma-keep"])
	}
	if len(calls) != len(api.DeletedMappingUUIDs) {
		t.Errorf("duplicate delete calls detected: %v", api.DeletedMappingUUIDs)
	}

	// No group-level mutators fire: beta/gamma are severed from the team but
	// their underlying GROUP objects are never destroyed.
	assertCounts(t, api, 0, 0)

	if len(result) != 1 {
		t.Fatalf("result len = %d, want 1", len(result))
	}
	if result[0].GroupName != "alpha" || result[0].GroupUUID != "ga" || result[0].MappingUUID != "ma-keep" {
		t.Errorf("surviving alpha = %+v, want alpha/ga/ma-keep", result[0])
	}
}

// survivors_invoke_zero_mutators: a desired set fully covered by prevOwned is a
// pure no-op — nothing to establish, nothing to prune.
func TestReconcileOIDC_SurvivorsInvokeZeroMutators(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "alpha", GroupUUID: "ga", TeamUUID: team, MappingUUID: "ma"},
		{GroupName: "beta", GroupUUID: "gb", TeamUUID: team, MappingUUID: "mb"},
	}

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"alpha", "beta"}, prevOwned)

	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
	if len(result) != 2 {
		t.Fatalf("result len = %d, want 2", len(result))
	}
	if result[0].GroupName != "alpha" || result[1].GroupName != "beta" {
		t.Errorf("result order/groupnames = %+v", result)
	}
}

// clear_all_removes_owned_edges_not_groups: clearing the desired set prunes
// every owned EDGE but never destroys the underlying GROUP objects.
func TestReconcileOIDC_ClearAllRemovesEdgesNotGroups(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "alpha", GroupUUID: "ga", TeamUUID: team, MappingUUID: "ma"},
		{GroupName: "beta", GroupUUID: "gb", TeamUUID: team, MappingUUID: "mb"},
		{GroupName: "gamma", GroupUUID: "gc", TeamUUID: team, MappingUUID: "mc"},
	}

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, nil, prevOwned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Every owned EDGE is pruned...
	if api.DeleteMappingByUuidCalls != 3 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 3", api.DeleteMappingByUuidCalls)
	}
	pruned := make(map[string]bool, len(api.DeletedMappingUUIDs))
	for _, u := range api.DeletedMappingUUIDs {
		pruned[u] = true
	}
	for _, want := range []string{"ma", "mb", "mc"} {
		if !pruned[want] {
			t.Errorf("expected %q to be pruned", want)
		}
	}
	// ...but GROUP objects are never deleted: no CreateGroup, no AddMapping.
	assertCounts(t, api, 0, 0)

	if len(result) != 0 {
		t.Errorf("result len = %d, want 0 (everything cleared)", len(result))
	}
}

// ---------------------------------------------------------------------------
// Wipe/restore invariants: status-loss safety (prevOwned nil/empty)
// ---------------------------------------------------------------------------
//
// INVARIANT: when prevOwned is nil/empty (status wiped by upgrade/status-
// subresource reset, or never persisted) the reconcile enters RESTORE-ONLY
// MODE: it establishes EVERY desired binding and emits ZERO DeleteMappingByUuid.
// Pruning is sourced SOLELY from acknowledged prior ownership, so a wipe can
// never cascade into deleting bindings authored by peers/the platform under
// concurrency. Restore reuses the SAME exact-case relink path, so a racing
// peer's group is recovered (never clobbered) even under StrictConflicts/HTTP
// 409.
//
// Covered:
//   - WipeRestoresAllAndDeletesNothing_NothingServerState
//   - WipeRestoresAllAndDeletesNothing_PartialServerState
//   - WipePlusConflictStillSafe
func TestReconcileOIDC_WipeRestoresAllAndDeletesNothing_NothingServerState(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"
	desired := []string{"alpha", "beta", "gamma"}

	// Status wiped: prevOwned nil. Server empty. Expect a full restore.
	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, desired, nil)

	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0 (wipe => restore-only, no pruning)", api.DeleteMappingByUuidCalls)
	}
	if api.CreateGroupCalls != len(desired) {
		t.Errorf("CreateGroup calls = %d, want %d", api.CreateGroupCalls, len(desired))
	}
	if api.AddMappingCalls != len(desired) {
		t.Errorf("AddMapping calls = %d, want %d", api.AddMappingCalls, len(desired))
	}
	if len(result) != len(desired) {
		t.Fatalf("result len = %d, want %d", len(result), len(desired))
	}
	for i, want := range []string{"alpha", "beta", "gamma"} {
		if result[i].GroupName != want {
			t.Errorf("result[%d].GroupName = %q, want %q", i, result[i].GroupName, want)
		}
		if result[i].GroupUUID == "" || result[i].MappingUUID == "" {
			t.Errorf("result[%d] has empty UUID(s): %+v", i, result[i])
		}
		if result[i].TeamUUID != team {
			t.Errorf("result[%d].TeamUUID = %q, want %q", i, result[i].TeamUUID, team)
		}
	}
	// Freshly allocated UUIDs must be pairwise distinct (create-bind).
	if result[0].GroupUUID == result[1].GroupUUID || result[1].GroupUUID == result[2].GroupUUID {
		t.Errorf("group UUIDs collided across restores: %+v", result)
	}
}

func TestReconcileOIDC_WipeRestoresAllAndDeletesNothing_PartialServerState(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:      true,
		ExistingGroups: []dtapi.OidcGroup{{Uuid: "grp-preexisting", Name: "alpha"}},
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"
	desired := []string{"alpha", "beta", "gamma"}

	// Status wiped (prevOwned nil) but the SERVER already holds one of the
	// desired groups. Restore reuses the existing group via list-lookup and
	// creates the remainder -- still ZERO deletes.
	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, desired, nil)

	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0 (restore never prunes server state)", api.DeleteMappingByUuidCalls)
	}
	// alpha reused (no CreateGroup); beta+gamma created.
	if api.CreateGroupCalls != 2 {
		t.Errorf("CreateGroup calls = %d, want 2 (reuse alpha, create beta+gamma)", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != len(desired) {
		t.Errorf("AddMapping calls = %d, want %d", api.AddMappingCalls, len(desired))
	}
	if len(result) != len(desired) {
		t.Fatalf("result len = %d, want %d", len(result), len(desired))
	}
	idx := groupIndexByName(result)
	reused, ok := idx["alpha"]
	if !ok || reused.GroupUUID != "grp-preexisting" {
		t.Errorf("alpha = %+v, want reuse of grp-preexisting", reused)
	}
	for _, name := range []string{"beta", "gamma"} {
		b, ok := idx[name]
		if !ok {
			t.Errorf("missing restored binding for %q", name)
			continue
		}
		if b.GroupUUID == "" || b.MappingUUID == "" {
			t.Errorf("%q has empty UUID(s): %+v", name, b)
		}
		if b.GroupUUID == "grp-preexisting" {
			t.Errorf("%q wrongly aliased the preexisting alpha group", name)
		}
	}
}

func TestReconcileOIDC_WipePlusConflictStillSafe(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:       true,
		StrictConflicts: true,
		// A rival reconciler slips "alpha" in after our initial RetrieveGroups
		// returns empty, provoking an HTTP 409-equivalent on CreateGroup(alpha).
		RacingCreators: []dtapi.OidcGroup{{Uuid: "racer-alpha-uuid", Name: "alpha"}},
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"
	desired := []string{"alpha", "beta", "gamma"}

	// Status wiped + a conflicting creator on one name: the lost-create must
	// relink-and-recover the winner by EXACT-CASE name, binding ALL desired
	// names and emitting ZERO deletes.
	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, desired, nil)

	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0 (wipe => restore-only)", api.DeleteMappingByUuidCalls)
	}
	if api.CreateGroupCalls != len(desired) {
		t.Errorf("CreateGroup calls = %d, want %d (each name attempts create; alpha loses race)", api.CreateGroupCalls, len(desired))
	}
	if api.AddMappingCalls != len(desired) {
		t.Errorf("AddMapping calls = %d, want %d", api.AddMappingCalls, len(desired))
	}
	if len(result) != len(desired) {
		t.Fatalf("result len = %d, want %d", len(result), len(desired))
	}
	idx := groupIndexByName(result)
	alpha, ok := idx["alpha"]
	if !ok || alpha.GroupUUID != "racer-alpha-uuid" {
		t.Errorf("alpha = %+v, want adopt racer-alpha-uuid via exact-case relink", alpha)
	}
	// beta/gamma created afresh (distinct UUIDs, non-empty, never aliasing the
	// racing alpha group).
	for _, name := range []string{"beta", "gamma"} {
		b, ok := idx[name]
		if !ok {
			t.Errorf("missing restored binding for %q", name)
			continue
		}
		if b.GroupUUID == "" || b.MappingUUID == "" {
			t.Errorf("%q has empty UUID(s): %+v", name, b)
		}
		if b.GroupUUID == "racer-alpha-uuid" {
			t.Errorf("%q wrongly aliased the racing alpha group", name)
		}
	}
}

// groupIndexByName indexes an OwnedMapping slice by GroupName for easy lookup
// in assertions. Last-write-wins is intentional: a well-formed reconcile never
// returns duplicate names, so a collision here is itself a failure signal.
func groupIndexByName(ms []OwnedMapping) map[string]OwnedMapping {
	out := make(map[string]OwnedMapping, len(ms))
	for _, m := range ms {
		out[m.GroupName] = m
	}
	return out
}

// reuse_existing_group_skips_CreateGroup: a group already present on the
// server is reused rather than re-created.
func TestReconcileOIDC_ReuseExistingGroupSkipsCreateGroup(t *testing.T) {
	api := &FakeOIDCAPI{
		Available:      true,
		ExistingGroups: []dtapi.OidcGroup{{Uuid: "grp-alpha-uuid", Name: "alpha"}},
	}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"alpha", "beta"}, nil)

	// alpha is served from the existing-group index (no CreateGroup); only
	// beta is materialised.
	if api.CreateGroupCalls != 1 {
		t.Errorf("CreateGroup calls = %d, want 1", api.CreateGroupCalls)
	}
	if api.AddMappingCalls != 2 {
		t.Errorf("AddMapping calls = %d, want 2", api.AddMappingCalls)
	}
	if len(result) != 2 {
		t.Fatalf("result len = %d, want 2", len(result))
	}
	if result[0].GroupUUID != "grp-alpha-uuid" || result[0].GroupName != "alpha" {
		t.Errorf("reused alpha = %+v", result[0])
	}
}

// unavailable_preserves_prevOwned_no_calls: when OIDC is offline the reconcile
// aborts with the sentinel and performs NO mutation, leaving prevOwned intact
// for the caller to retain and retry.
func TestReconcileOIDC_UnavailablePreservesPrevOwnedNoCalls(t *testing.T) {
	api := &FakeOIDCAPI{Available: false}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "alpha", GroupUUID: "ga", TeamUUID: team, MappingUUID: "ma"},
	}

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, []string{"alpha"}, prevOwned)

	if !errors.Is(err, ErrOIDCUnavailable) {
		t.Fatalf("err = %v, want ErrOIDCUnavailable", err)
	}
	if result != nil {
		t.Errorf("expected nil result on unavailability, got %+v", result)
	}
	// No mutation whatsoever: prevOwned is retained by the caller for retry.
	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
}

// reorder_only_invokes_no_mutators: declaring the same set in a different order
// is a pure no-op — every previous owner survives, toEstablish is empty, and
// the server is never mutated. Output preserves prevOwned order.
func TestReconcileOIDC_ReorderInvokesNoMutators(t *testing.T) {
	api := newFake()
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "alpha", GroupUUID: "ga", TeamUUID: team, MappingUUID: "ma"},
		{GroupName: "beta", GroupUUID: "gb", TeamUUID: team, MappingUUID: "mb"},
	}

	// Declared in the opposite order from prevOwned.
	result := mustSingle(t, ReconcileOIDCMappings, log, api, team, []string{"beta", "alpha"}, prevOwned)

	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
	if len(result) != 2 {
		t.Fatalf("result len = %d, want 2", len(result))
	}
	// Output preserves PREVIOUS-ownership order, proving the reorder did not
	// drive reconstruction.
	if result[0].GroupName != "alpha" || result[1].GroupName != "beta" {
		t.Errorf("expected survivor order [alpha,beta], got [%s,%s]", result[0].GroupName, result[1].GroupName)
	}
}

// ---------------------------------------------------------------------------
// Availability gating
// ---------------------------------------------------------------------------

func TestReconcileOIDC_UnavailableReturnsSentinelWithoutCallingMutators(t *testing.T) {
	api := &FakeOIDCAPI{Available: false}
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, []string{"alpha"}, nil)

	if !errors.Is(err, ErrOIDCUnavailable) {
		t.Fatalf("err = %v, want ErrOIDCUnavailable", err)
	}
	if result != nil {
		t.Errorf("expected nil slice on unavailability, got %v", result)
	}
	assertCounts(t, api, 0, 0)
	if api.DeleteMappingByUuidCalls != 0 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 0", api.DeleteMappingByUuidCalls)
	}
}

// ---------------------------------------------------------------------------
// Error propagation: each upstream failure aborts and leaves later stages cold
// ---------------------------------------------------------------------------

func TestReconcileOIDC_ErrorPropagation(t *testing.T) {
	baseErr := errors.New("simulated upstream fault")

	type stage struct {
		name   string
		wire   func(*FakeOIDCAPI)
		expect func(*testing.T, *FakeOIDCAPI, error)
	}
	stages := []stage{
		{
			name: "IsAvailable faults",
			wire: func(f *FakeOIDCAPI) { f.ErrIsAvailable = baseErr },
			expect: func(t *testing.T, f *FakeOIDCAPI, err error) {
				t.Helper()
				if err == nil || !errors.Is(err, baseErr) {
					t.Fatalf("err = %v, want to wrap %v", err, baseErr)
				}
				assertCounts(t, f, 0, 0)
				if f.DeleteMappingByUuidCalls != 0 {
					t.Errorf("DeleteMappingByUuid calls = %d, want 0", f.DeleteMappingByUuidCalls)
				}
			},
		},
		{
			name: "RetrieveGroups faults",
			wire: func(f *FakeOIDCAPI) { f.ErrRetrieveGroups = baseErr },
			expect: func(t *testing.T, f *FakeOIDCAPI, err error) {
				t.Helper()
				if err == nil || !errors.Is(err, baseErr) {
					t.Fatalf("err = %v, want to wrap %v", err, baseErr)
				}
				assertCounts(t, f, 0, 0)
				if f.DeleteMappingByUuidCalls != 0 {
					t.Errorf("DeleteMappingByUuid calls = %d, want 0", f.DeleteMappingByUuidCalls)
				}
			},
		},
		{
			name: "CreateGroup faults",
			wire: func(f *FakeOIDCAPI) { f.ErrCreateGroup = baseErr },
			expect: func(t *testing.T, f *FakeOIDCAPI, err error) {
				t.Helper()
				if err == nil || !errors.Is(err, baseErr) {
					t.Fatalf("err = %v, want to wrap %v", err, baseErr)
				}
				// Exactly one CreateGroup attempted before the abort; no mapping.
				if f.CreateGroupCalls != 1 {
					t.Errorf("CreateGroup calls = %d, want 1", f.CreateGroupCalls)
				}
				assertCounts(t, f, 1, 0)
			},
		},
		{
			name: "AddMapping2 faults",
			wire: func(f *FakeOIDCAPI) { f.ErrAddMapping2 = baseErr },
			expect: func(t *testing.T, f *FakeOIDCAPI, err error) {
				t.Helper()
				if err == nil || !errors.Is(err, baseErr) {
					t.Fatalf("err = %v, want to wrap %v", err, baseErr)
				}
				assertCounts(t, f, 1, 1)
			},
		},
	}

	const team = "deadbeef-0000-0000-0000-000000000001"
	for _, st := range stages {
		t.Run(st.name, func(t *testing.T) {
			f := newFake()
			st.wire(f)
			_, err := ReconcileOIDCMappings(context.Background(), &recordingLogger{}, f, team, []string{"alpha", "beta"}, nil)
			st.expect(t, f, err)
		})
	}
}

// delete_fault_aborts_before_establishing: a failing prune aborts the reconcile
// BEFORE any CreateGroup/AddMapping runs, preventing half-applied convergence.
// Ordering matters: pruning precedes establishment, so a transient delete error
// must roll back the whole operation rather than partially rebuild the set.
func TestReconcileOIDC_DeleteFaultAbortsBeforeEstablishing(t *testing.T) {
	api := newFake()
	api.ErrDeleteMapping = errors.New("transient delete fault")
	log := &recordingLogger{}
	const team = "deadbeef-0000-0000-0000-000000000001"

	prevOwned := []OwnedMapping{
		{GroupName: "gone", GroupUUID: "gg", TeamUUID: team, MappingUUID: "mg-delete-me"},
	}

	result, err := ReconcileOIDCMappings(context.Background(), log, api, team, []string{"keep"}, prevOwned)

	if err == nil || !strings.Contains(err.Error(), "pruning stale OIDC mapping") {
		t.Fatalf("err = %v, want a pruning-wrap error", err)
	}
	if result != nil {
		t.Errorf("expected nil result on prune fault, got %+v", result)
	}
	// The prune was attempted exactly once; establishment never began.
	if api.DeleteMappingByUuidCalls != 1 {
		t.Errorf("DeleteMappingByUuid calls = %d, want 1", api.DeleteMappingByUuidCalls)
	}
	assertCounts(t, api, 0, 0)
}

// ---------------------------------------------------------------------------
// Invariants: group objects are never deleted through the seam
// ---------------------------------------------------------------------------
//
// INVARIANT: Dependency-Track OIDC GROUP objects are never destroyed by this
// operator. A team's OIDC bindings are reconciled by adjusting the membership
// EDGES only -- established with AddMapping2 and torn down with
// DeleteMappingByUuid -- whereas the underlying Group objects are CREATE-only,
// never DELETE. Erasing a shared group object would silently revoke memberships
// owed to OTHER tenants, breaking tenant isolation.
//
// Guaranteed BY CONSTRUCTION: the [OIDCAClient] seam intentionally exposes no
// DeleteGroup method -- only IsAvailable, RetrieveGroups, CreateGroup,
// AddMapping2 and DeleteMappingByUuid. Because the compiler-typed interface has
// no such verb, no controller code path can ever reach it.
//
// Enforcement is dual: an executable guard in
// TestReconcileOIDC_GroupObjectsAreNeverDeletableThroughSeam reflects over the
// [OIDCAClient] interface and fails if any method carries the "DeleteGroup"
// prefix; and, as a static human-readable smoke check, the identical claim is
// asserted with:
//
//	! grep -q DeleteGroup internal/dependencytrack/oidc_mapping.go
//
// (successful exit => the production seam file defines no such symbol).
func TestReconcileOIDC_GroupObjectsAreNeverDeletableThroughSeam(t *testing.T) {
	it := reflect.TypeOf((*OIDCAClient)(nil)).Elem()
	for i := 0; i < it.NumMethod(); i++ {
		m := it.Method(i)
		if strings.HasPrefix(m.Name, "DeleteGroup") {
			t.Errorf("OIDCAClient seam exposes destructive method %q; "+
				"group objects must never be deletable through the seam", m.Name)
		}
	}
}

// ---------------------------------------------------------------------------
// Adapters wiring sanity: dtapiAdapter delegates faithfully
// ---------------------------------------------------------------------------

func TestDtapiAdapter_SatisfiesInterface(t *testing.T) {
	var _ OIDCAClient = (*dtapiAdapter)(nil)
	var _ OIDCAClient = NewDtapiOIDCClient(nil)
}

// ---------------------------------------------------------------------------
// Assertions
// ---------------------------------------------------------------------------

func assertCounts(t *testing.T, f *FakeOIDCAPI, wantCreates, wantMappings int) {
	t.Helper()
	if f.CreateGroupCalls != wantCreates {
		t.Errorf("CreateGroup calls = %d, want %d", f.CreateGroupCalls, wantCreates)
	}
	if f.AddMappingCalls != wantMappings {
		t.Errorf("AddMapping calls = %d, want %d", f.AddMappingCalls, wantMappings)
	}
}
