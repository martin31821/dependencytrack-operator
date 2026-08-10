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
	"errors"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/martin31821/dependencytrack-operator/gen/dtapi"
	"github.com/martin31821/dependencytrack-operator/internal/dependencytrack"
)

// TestReconcileOIDC_SkipsWhenSpecAbsent locks in the T04 invariant that a Team
// without an OIDC stanza incurs ZERO DependencyTrack API traffic and emits ZERO
// events. The apiClient argument is intentionally nil: the nil-pointer guard
// returns before the generated client is ever consulted, proving the zero-cost
// fast path is safe even when no real client is wired up.
func TestReconcileOIDC_SkipsWhenSpecAbsent(t *testing.T) {
	rec := record.NewFakeRecorder(16)
	r := &TeamReconciler{Recorder: rec}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "noop"},
		Spec:       dependencytrackv1alpha1.TeamSpec{}, // OIDC == nil
	}

	res, err := r.reconcileOIDC(context.Background(), nil, team, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsZero() {
		t.Fatalf("expected zero reconcile result, got %+v", res)
	}
	if team.Status.OIDC != nil {
		t.Fatalf("expected status.OIDC to remain nil, got %#v", team.Status.OIDC)
	}

	events := drainFakeEvents(rec)
	if len(events) != 0 {
		t.Fatalf("expected zero recorder events, got %d: %v", len(events), events)
	}
}

// TestOwnedMappingsToStatus_TableDriven pins the seam <-> status-field mapping
// contract introduced by reconcileOIDC. Boundary cases (nil/empty input yielding
// a non-nil empty slice so "reconciled, no mappings" is distinguishable from
// "never reconciled") are asserted here because they govern status semantics.
func TestOwnedMappingsToStatus_TableDriven(t *testing.T) {
	in := []dependencytrack.OwnedMapping{
		{GroupName: "admins", GroupUUID: "g-admin", TeamUUID: "t-1", MappingUUID: "m-1"},
		{GroupName: "developers", GroupUUID: "g-dev", TeamUUID: "t-1", MappingUUID: "m-2"},
	}

	out := ownedMappingsToStatus(in)
	if len(out) != len(in) {
		t.Fatalf("translation length mismatch: got %d want %d", len(out), len(in))
	}
	for i, m := range in {
		got := out[i]
		if got != (dependencytrackv1alpha1.OwnedOIDCMapping{
			GroupName:   m.GroupName,
			GroupUUID:   m.GroupUUID,
			TeamUUID:    m.TeamUUID,
			MappingUUID: m.MappingUUID,
		}) {
			t.Errorf("row %d: translation drift\n got={%+v}\nwant={%+v}", i, got, m)
		}
	}

	// Boundary: nil and empty inputs both collapse to a non-nil, empty slice.
	for _, src := range [][]dependencytrack.OwnedMapping{nil, {}} {
		got := ownedMappingsToStatus(src)
		if len(got) != 0 {
			t.Errorf("expected empty slice for nil/empty input, got len=%d", len(got))
		}
	}
}

// drainFakeEvents collects pending FakeRecorder events for failure reporting.
func drainFakeEvents(rec *record.FakeRecorder) []string {
	close(rec.Events)
	var ev []string
	for e := range rec.Events {
		ev = append(ev, e)
	}
	return ev
}

// ---------------------------------------------------------------------------
// Controller-level OIDC seam stub
// ---------------------------------------------------------------------------

// fakeMapReq captures the arguments handed to AddMapping2 for assertions.
type fakeMapReq struct {
	groupUUID string
	teamUUID  string
}

// fakeOIDCAClient is a minimal [dependencytrack.OIDCAClient] stub for
// controller-level tests. It bypasses the dtapiAdapter boundary entirely
// (no generated client, no HTTP) and exposes only the knobs the controller
// exercises so interaction counts and ordering can be asserted.
type fakeOIDCAClient struct {
	available      bool
	isAvailErr     error
	existingGroups []dtapi.OidcGroup
	retrieveErr    error
	retrieveCalls  int

	createErr error
	addErr    error
	deleteErr error

	creates []string
	maps    []fakeMapReq
	deletes []string
}

// Compile-time proof the stub satisfies the seam interface.
var _ dependencytrack.OIDCAClient = (*fakeOIDCAClient)(nil)

func (f *fakeOIDCAClient) IsAvailable(context.Context) (bool, error) {
	return f.available, f.isAvailErr
}

func (f *fakeOIDCAClient) RetrieveGroups(context.Context) ([]dtapi.OidcGroup, error) {
	f.retrieveCalls++
	if f.retrieveErr != nil {
		return nil, f.retrieveErr
	}
	return f.existingGroups, nil
}

func (f *fakeOIDCAClient) CreateGroup(_ context.Context, name string) (dtapi.OidcGroup, error) {
	f.creates = append(f.creates, name)
	if f.createErr != nil {
		return dtapi.OidcGroup{}, f.createErr
	}
	// Synthetic, predictable UUID so status assertions can pin provenance.
	return dtapi.OidcGroup{Uuid: "grp::" + name, Name: name}, nil
}

func (f *fakeOIDCAClient) AddMapping2(_ context.Context, groupUUID, teamUUID string) (dtapi.MappedOidcGroup, error) {
	f.maps = append(f.maps, fakeMapReq{groupUUID: groupUUID, teamUUID: teamUUID})
	if f.addErr != nil {
		return dtapi.MappedOidcGroup{}, f.addErr
	}
	return dtapi.MappedOidcGroup{Uuid: "map::" + groupUUID}, nil
}

func (f *fakeOIDCAClient) DeleteMappingByUuid(_ context.Context, mappingUUID string) error {
	f.deletes = append(f.deletes, mappingUUID)
	return f.deleteErr
}

// equalOwnedMappings compares two status-ownership slices by value.
// OwnedOIDCMapping is a strings-only struct, hence directly comparable.
func equalOwnedMappings(got, want []dependencytrackv1alpha1.OwnedOIDCMapping) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestReconcileOIDC_PassesPrevOwnedSurvivorship proves the diff engine is
// anchored on yesterday's status: a group dropped from spec while still owned
// in status must be pruned (DeleteMappingByUuid invoked exactly once), whereas
// a newly-desired group must be claimed. The prevOwned input is sourced
// exclusively from team.Status.OIDC — no hardcoded seeds leak into the seam.
func TestReconcileOIDC_PassesPrevOwnedSurvivorship(t *testing.T) {
	const teamUUID = "t-uuid"
	fake := &fakeOIDCAClient{available: true}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "survivors"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"admins"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{
				OwnedMappings: []dependencytrackv1alpha1.OwnedOIDCMapping{
					{GroupName: "devs", GroupUUID: "g-devs", TeamUUID: teamUUID, MappingUUID: "m-devs"},
				},
			},
		},
	}

	res, err := r.reconcileOIDC(context.Background(), fake, team, teamUUID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsZero() {
		t.Fatalf("expected zero reconcile result, got %+v", res)
	}

	// The stale "devs" edge must be pruned exactly once — survival hinges on
	// the prevOwned anchor reaching ReconcileOIDCMappings.
	if got, want := len(fake.deletes), 1; got != want {
		t.Fatalf("DeleteMappingByUuid calls: got %d want %d (%v)", got, want, fake.deletes)
	}
	if fake.deletes[0] != "m-devs" {
		t.Fatalf("pruned unexpected mapping uuid: got %q want %q", fake.deletes[0], "m-devs")
	}

	// "admins" was not previously owned => a new group + mapping must be claimed.
	if got, want := len(fake.creates), 1; got != want {
		t.Fatalf("CreateGroup calls: got %d want %d (%v)", got, want, fake.creates)
	}
	if fake.creates[0] != "admins" {
		t.Fatalf("expected CreateGroup(admins), got %q", fake.creates[0])
	}
	if got, want := len(fake.maps), 1; got != want {
		t.Fatalf("AddMapping2 calls: got %d want %d", got, want)
	}

	// Converged status records the net result: admins claimed, devs pruned.
	if team.Status.OIDC == nil {
		t.Fatal("expected status.OIDC to be repopulated")
	}
	wantMappings := []dependencytrackv1alpha1.OwnedOIDCMapping{
		{GroupName: "admins", GroupUUID: "grp::admins", TeamUUID: teamUUID, MappingUUID: "map::grp::admins"},
	}
	if !equalOwnedMappings(team.Status.OIDC.OwnedMappings, wantMappings) {
		t.Fatalf("converged status mismatch:\n got=%v\nwant=%v", team.Status.OIDC.OwnedMappings, wantMappings)
	}
}

// TestReconcileOIDC_StatusRoundTripPreservesUnchangedEdges locks in the
// idempotency invariant: when the desired set is identical to the owned set,
// the reconcile issues ZERO server mutations and the status is carried
// through byte-for-byte (survivors emitted verbatim, no rewrites).
func TestReconcileOIDC_StatusRoundTripPreservesUnchangedEdges(t *testing.T) {
	const teamUUID = "t-uuid"
	fake := &fakeOIDCAClient{available: true}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	prevOwned := []dependencytrackv1alpha1.OwnedOIDCMapping{
		{GroupName: "admins", GroupUUID: "g-admin", TeamUUID: teamUUID, MappingUUID: "m-admin"},
		{GroupName: "devs", GroupUUID: "g-dev", TeamUUID: teamUUID, MappingUUID: "m-dev"},
	}
	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "steady"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"admins", "devs"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{OwnedMappings: prevOwned},
		},
	}

	res, err := r.reconcileOIDC(context.Background(), fake, team, teamUUID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsZero() {
		t.Fatalf("expected zero reconcile result, got %+v", res)
	}

	// Steady state => no server mutations of any kind.
	if len(fake.deletes) != 0 {
		t.Fatalf("expected zero deletions, got %v", fake.deletes)
	}
	if len(fake.creates) != 0 {
		t.Fatalf("expected zero creations, got %v", fake.creates)
	}
	if len(fake.maps) != 0 {
		t.Fatalf("expected zero mappings, got %v", fake.maps)
	}
	if fake.retrieveCalls != 0 {
		t.Fatalf("expected zero group retrievals, got %d", fake.retrieveCalls)
	}

	// Status survives byte-for-byte: survivors emitted verbatim, no rewrites.
	if team.Status.OIDC == nil {
		t.Fatal("expected status.OIDC preserved")
	}
	if !equalOwnedMappings(team.Status.OIDC.OwnedMappings, prevOwned) {
		t.Fatalf("status drifted on steady-state reconcile:\n got=%v\nwant=%v", team.Status.OIDC.OwnedMappings, prevOwned)
	}
}

// TestReconcileOIDC_RestoreUnderWipedStatus drives the controller-level restore
// path end-to-end through reconcileOIDC: a team whose status.OIDC was wiped
// (nil) is rebuilt from scratch. This complements the seam-level
// TestReconcileOIDC_WipeRestoresAllAndDeletesNothing_* tests by pinning the
// controller CONTRACT -- namely that the rewoven bindings land in
// team.Status.OIDC as a 2-row set, zero deletions occur (restore-only mode),
// and a Normal "bindings ensured" event is emitted.
func TestReconcileOIDC_RestoreUnderWipedStatus(t *testing.T) {
	const teamUUID = "t-restore"
	fake := &fakeOIDCAClient{available: true /* existingGroups defaults to nil */}
	rec := record.NewFakeRecorder(16)
	r := &TeamReconciler{Recorder: rec}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "restored"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"alpha", "beta"}},
		},
		// Status.OIDC intentionally nil: simulates a wiped/reset status
		// subresource (upgrade, migration, never-reconciled).
	}

	res, err := r.reconcileOIDC(context.Background(), fake, team, teamUUID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsZero() {
		t.Fatalf("expected zero reconcile result, got %+v", res)
	}

	// Restore-only: no deletions ever fire from a wiped status.
	if len(fake.deletes) != 0 {
		t.Fatalf("expected zero deletions under wiped status, got %v", fake.deletes)
	}

	// Desired order preserved: alpha then beta.
	if got, want := fake.creates, []string{"alpha", "beta"}; len(got) != len(want) {
		t.Fatalf("CreateGroup calls = %d, want %d (%v)", len(got), len(want), got)
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("CreateGroup[%d] = %q, want %q (%v)", i, got[i], want[i], got)
			}
		}
	}
	if len(fake.maps) != 2 {
		t.Fatalf("expected 2 AddMapping2 calls, got %d (%v)", len(fake.maps), fake.maps)
	}

	// team.Status.OIDC must be rewired with exactly 2 rows carrying the
	// synthetic UUIDs minted by the fake seam.
	if team.Status.OIDC == nil {
		t.Fatal("expected status.OIDC to be repopulated after restore")
	}
	gotRows := team.Status.OIDC.OwnedMappings
	if len(gotRows) != 2 {
		t.Fatalf("expected status.OIDC to hold 2 rows, got %d: %+v", len(gotRows), gotRows)
	}
	wantRows := []dependencytrackv1alpha1.OwnedOIDCMapping{
		{GroupName: "alpha", GroupUUID: "grp::alpha", TeamUUID: teamUUID, MappingUUID: "map::grp::alpha"},
		{GroupName: "beta", GroupUUID: "grp::beta", TeamUUID: teamUUID, MappingUUID: "map::grp::beta"},
	}
	if !equalOwnedMappings(gotRows, wantRows) {
		t.Fatalf("rewired status mismatch:\n got=%+v\nwant=%+v", gotRows, wantRows)
	}

	// A Normal "bindings ensured" event must announce convergence.
	ev := drainFakeEvents(rec)
	if !containsSubstring(ev, "OIDCMappingCreated") {
		t.Fatalf("expected a Normal OIDCMappingCreated event, got %v", ev)
	}
}

// TestReconcileOIDC_UnavailabilityLeavesStatusUntouched closes the controller-
// level gap left by the seam's TestReconcileOIDC_Unavailable*: although the
// diff engine proves it mutates nothing on ErrOIDCUnavailable, only the
// controller can prove team.Status.OIDC is NOT NULLED mid-flight. This guards
// the status-loss hazard encoded in the conservative-restore invariant -- a
// reachable instance must not overwrite a previously-acknowledged ownership set
// with a blanked one, which would orphan the bindings the controller does not
// own. The unavailability branch wires a fake client so its Status().Update
// call (which flushes the advisory condition) is exercisable without envtest.
func TestReconcileOIDC_UnavailabilityLeavesStatusUntouched(t *testing.T) {
	const teamUUID = "t-locked"
	fake := &fakeOIDCAClient{available: false}
	cl := newFakeCtrlClient(t)
	rec := record.NewFakeRecorder(16)
	r := &TeamReconciler{Client: cl, Recorder: rec}

	// A team carrying LIVE owned mappings in status -- the very rows a wiper
	// would destroy if the controller reassigned status.OIDC on the
	// unavailability path.
	prevOwned := []dependencytrackv1alpha1.OwnedOIDCMapping{
		{GroupName: "legacy-alpha", GroupUUID: "ga", TeamUUID: teamUUID, MappingUUID: "ma"},
		{GroupName: "legacy-beta", GroupUUID: "gb", TeamUUID: teamUUID, MappingUUID: "mb"},
	}
	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "locked"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"brand-new"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{OwnedMappings: prevOwned},
		},
	}
	// Seed the cluster store so the Status().Update call in the unavailability
	// branch persists without panicking on a nil/absent client.
	if err := cl.Create(context.Background(), team); err != nil {
		t.Fatalf("seeding team: %v", err)
	}

	// Independent snapshot of the protected status-OIDC payload.
	before := team.Status.OIDC.DeepCopy()

	res, err := r.reconcileOIDC(context.Background(), fake, team, teamUUID)

	// Transient: requeue for retry.
	if !errors.Is(err, dependencytrack.ErrOIDCUnavailable) {
		t.Fatalf("err = %v, want ErrOIDCUnavailable", err)
	}
	if !res.Requeue {
		t.Fatalf("expected Requeue=true on unavailability, got %+v", res)
	}

	// Provably untouched: the controller must NOT have reassigned status.OIDC.
	after := team.Status.OIDC
	if after == nil {
		t.Fatal("status.OIDC was nulled during unavailability -- violates conservative-restore")
	}
	if !equalOwnedMappings(after.OwnedMappings, before.OwnedMappings) {
		t.Fatalf("status.OIDC was mutated during unavailability:\n before=%+v\n after=%+v",
			before.OwnedMappings, after.OwnedMappings)
	}

	// Unavailability short-circuits BEFORE any seam mutation: zero mutators.
	if fake.retrieveCalls != 0 {
		t.Fatalf("expected zero group retrievals on unavailability, got %d", fake.retrieveCalls)
	}
	if len(fake.creates) != 0 {
		t.Fatalf("expected zero CreateGroup calls on unavailability, got %v", fake.creates)
	}
	if len(fake.maps) != 0 {
		t.Fatalf("expected zero AddMapping2 calls on unavailability, got %v", fake.maps)
	}
	if len(fake.deletes) != 0 {
		t.Fatalf("expected zero DeleteMappingByUuid calls on unavailability, got %v", fake.deletes)
	}

	// A Warning advisory event must announce the transient outage.
	ev := drainFakeEvents(rec)
	if !containsSubstring(ev, "OIDCUnavailable") {
		t.Fatalf("expected a Warning OIDCUnavailable event, got %v", ev)
	}
}

// ---------------------------------------------------------------------------
// scrubOwnedOIDCMappings -- S04 deletion contracts
// ---------------------------------------------------------------------------

// TestScrubOwnedOIDCMappings_DeletesEveryOwnedEdge pins the happy-path scrub:
// every owned mapping edge is DELETE-MAPPING-BY-UUID'd exactly once, in
// recorded order, with ZERO provisioning mutators (no creates/maps/retrieves)
// fired -- the scrub is a pure reclamation sweep, not a re-provision.
func TestScrubOwnedOIDCMappings_DeletesEveryOwnedEdge(t *testing.T) {
	fake := &fakeOIDCAClient{available: true}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "edge-washer"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"irrelevant"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{
				OwnedMappings: []dependencytrackv1alpha1.OwnedOIDCMapping{
					{MappingUUID: "m-aaa", GroupName: "g1"},
					{MappingUUID: "m-bbb", GroupName: "g2"},
					{MappingUUID: "m-ccc", GroupName: "g3"},
				},
			},
		},
	}

	err := r.scrubOwnedOIDCMappings(context.Background(), fake, team)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Exactly one delete per owned edge, in enumeration order.
	if got, want := len(fake.deletes), 3; got != want {
		t.Fatalf("DeleteMappingByUuid calls = %d, want %d (%v)", got, want, fake.deletes)
	}
	for i, want := range []string{"m-aaa", "m-bbb", "m-ccc"} {
		if fake.deletes[i] != want {
			t.Fatalf("deletes[%d] = %q, want %q", i, fake.deletes[i], want)
		}
	}

	// Pure reclaim: no provisioning mutators must fire.
	if len(fake.creates) != 0 {
		t.Fatalf("expected zero CreateGroup calls, got %v", fake.creates)
	}
	if len(fake.maps) != 0 {
		t.Fatalf("expected zero AddMapping2 calls, got %v", fake.maps)
	}
	if fake.retrieveCalls != 0 {
		t.Fatalf("expected zero RetrieveGroups calls, got %d", fake.retrieveCalls)
	}
}

// TestScrubOwnedOIDCMappings_NoopWhenNoOwnership locks the fast path: a team
// that opts into OIDC but owns no mappings (status.OIDC nil) scrubs nothing and
// SKIPS the availability probe -- zero seam traffic of any kind. Distinct from
// the vanilla-passthrough case (which has Spec.OIDC == nil).
func TestScrubOwnedOIDCMappings_NoopWhenNoOwnership(t *testing.T) {
	fake := &fakeOIDCAClient{available: true}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "none-yet"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"something"}},
		},
		// status.OIDC intentionally nil: nothing was ever owned/reclaimed.
	}

	err := r.scrubOwnedOIDCMappings(context.Background(), fake, team)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertZeroScrubTraffic(t, fake, "noop-when-no-ownership")
}

// TestScrubOwnedOIDCMappings_VanillaTeamPassthrough proves the self-guard:
// even when LEFTOVER status.OIDC mappings linger (e.g. from a rolled-back
// enablement), a team whose Spec.OIDC == nil scrubs nothing and issues zero
// API traffic -- the guard fires before the availability probe, so orphaned
// status rows are harmlessly inert.
func TestScrubOwnedOIDCMappings_VanillaTeamPassthrough(t *testing.T) {
	fake := &fakeOIDCAClient{available: true}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "vanilla"},
		Spec:       dependencytrackv1alpha1.TeamSpec{}, // OIDC == nil
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{
				OwnedMappings: []dependencytrackv1alpha1.OwnedOIDCMapping{
					{MappingUUID: "orphan-left-behind", GroupName: "straggler"},
				},
			},
		},
	}

	err := r.scrubOwnedOIDCMappings(context.Background(), fake, team)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	assertZeroScrubTraffic(t, fake, "vanilla-team-passthrough")
}

// TestScrubOwnedOIDCMappings_UnavailableDefersWithoutMutators pins the deferral
// contract: when OIDC is not provisioned, the scrub returns
// ErrOIDCUnavailable and performs NO mutations whatsoever (not even the
// availability probe's side-effect-free counterpart RetrieveGroups), leaving the
// caller free to retain the finalizer and requeue.
func TestScrubOwnedOIDCMappings_UnavailableDefersWithoutMutators(t *testing.T) {
	fake := &fakeOIDCAClient{available: false}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "unprovisioned"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"ghost"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{
				OwnedMappings: []dependencytrackv1alpha1.OwnedOIDCMapping{
					{MappingUUID: "will-not-be-deleted", GroupName: "ghost"},
				},
			},
		},
	}

	err := r.scrubOwnedOIDCMappings(context.Background(), fake, team)
	if !errors.Is(err, dependencytrack.ErrOIDCUnavailable) {
		t.Fatalf("err = %v, want ErrOIDCUnavailable", err)
	}

	assertZeroScrubTraffic(t, fake, "unavailable-deferral")
}

// TestScrubOwnedOIDCMappings_FailStopsOnError locks the fail-stop contract:
// the shim records the offending UUID BEFORE surfacing the injected error, so
// the first failing edge is observed exactly once and subsequent edges are
// never attempted. The returned error must wrap the inner transport fault so
// callers can discriminate it via errors.Is.
func TestScrubOwnedOIDCMappings_FailStopsOnError(t *testing.T) {
	sentinel := errors.New("transport refused")
	fake := &fakeOIDCAClient{
		available: true,
		deleteErr: sentinel,
	}
	r := &TeamReconciler{Recorder: record.NewFakeRecorder(16)}

	team := &dependencytrackv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{Name: "fragile"},
		Spec: dependencytrackv1alpha1.TeamSpec{
			OIDC: &dependencytrackv1alpha1.TeamOIDCConfig{Groups: []string{"always-fails"}},
		},
		Status: dependencytrackv1alpha1.TeamStatus{
			OIDC: &dependencytrackv1alpha1.TeamOIDCStatus{
				OwnedMappings: []dependencytrackv1alpha1.OwnedOIDCMapping{
					{MappingUUID: "edge-first", GroupName: "g1"},
					{MappingUUID: "edge-second", GroupName: "g2"},
					{MappingUUID: "edge-third", GroupName: "g3"},
				},
			},
		},
	}

	err := r.scrubOwnedOIDCMappings(context.Background(), fake, team)
	if err == nil {
		t.Fatal("expected a wrapped error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want it to wrap %v", err, sentinel)
	}

	// Fail-stop: exactly one attempt, on the first edge, recorded pre-error.
	if got, want := len(fake.deletes), 1; got != want {
		t.Fatalf("DeleteMappingByUuid calls = %d, want %d (fail-stop after first edge)", got, want)
	}
	if fake.deletes[0] != "edge-first" {
		t.Fatalf("deletes[0] = %q, want %q", fake.deletes[0], "edge-first")
	}

	// Abort-on-first-edge: no further mutators fire.
	if len(fake.creates) != 0 {
		t.Fatalf("expected zero CreateGroup calls, got %v", fake.creates)
	}
	if len(fake.maps) != 0 {
		t.Fatalf("expected zero AddMapping2 calls, got %v", fake.maps)
	}
	if fake.retrieveCalls != 0 {
		t.Fatalf("expected zero RetrieveGroups calls, got %d", fake.retrieveCalls)
	}
}

// assertZeroScrubTraffic centralizes the "total absence of seam activity"
// assertion shared by the noop/deferral/fast-path scrub contracts.
func assertZeroScrubTraffic(t *testing.T, fake *fakeOIDCAClient, tag string) {
	t.Helper()
	if len(fake.deletes) != 0 {
		t.Fatalf("[%s] expected zero deletes, got %v", tag, fake.deletes)
	}
	if len(fake.creates) != 0 {
		t.Fatalf("[%s] expected zero creates, got %v", tag, fake.creates)
	}
	if len(fake.maps) != 0 {
		t.Fatalf("[%s] expected zero maps, got %v", tag, fake.maps)
	}
	if fake.retrieveCalls != 0 {
		t.Fatalf("[%s] expected zero retrieveCalls, got %d", tag, fake.retrieveCalls)
	}
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// newFakeCtrlClient builds a controller-runtime fake client with the
// dependencytrackv1alpha1 scheme registered, so controller paths that persist
// status (e.g. the ErrOIDCUnavailable branch of reconcileOIDC) are exercisable
// without standing up envtest.
func newFakeCtrlClient(t *testing.T) client.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := dependencytrackv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("registering dependencytrackv1alpha1 scheme: %v", err)
	}
	return fakeclient.NewClientBuilder().WithScheme(scheme).Build()
}

// containsSubstring reports whether any drained FakeRecorder event string
// contains substr. The FakeRecorder encodes events as "<EventType> <Reason>
// <Message>", so matching on the reason token is stable across message wording.
func containsSubstring(events []string, substr string) bool {
	for _, e := range events {
		if strings.Contains(e, substr) {
			return true
		}
	}
	return false
}
