# Canonical Verification Gates

Structured rerun manifest for milestone **M001**, slice **S06**. Produced by
task **T05** ("Compile/harden the test-e2e pathway and emit structured
verification evidence"). Each of these Canonical verification gates is independently re-runnable;
the milestone validator re-executes the Command column against the current
source revision and reads the Current Evidence table for the last asserted status.

> `make test-e2e` is a canonical CI gate but **cannot execute in this sandbox**
> (no Docker/Kind). Its precondition is documented in §2 below. The real
> Dependency-Track lifecycle proof therefore remains gated on a
> Docker/Kind-capable runner — no loopback stub substitutes for the live
> fixture. T05 restores the *compile/contract* assurance and the structured
> evidence surface so the gate can be rerun deterministically elsewhere.

## 1. Canonical Gate Matrix

Legend: ✅ Green = asserted at the quoted source revision · 🧭 Precond = environmental prerequisite · ⏳ Blocked = requires an external runtime absent here.

| # | Gate | Criterion | Command (rerun) | Status | Evidence ref |
|---|------|-----------|------------------|--------|--------------|
| G1 | Webhook marker reproducible | `git diff --quiet config/webhook/manifests.yaml` exits 0 (marker/idempotent) | `git diff --quiet config/webhook/manifests.yaml` | ✅ Green | commit `04562d4` locks `+kubebuilder:webhook`; `a2204d7` mounts `-webhook-server-cert` (MEM023) |
| G2 | Manifests/generate idempotent | `make manifests generate fmt` leaves no uncommitted diff | `make manifests generate fmt && git diff --quiet` | ✅ Green | inherited from T02 marker lock + T04 pipeline |
| G3 | Helm/sync in-sync | chart artefacts stay aligned with kustomize/RBAC/CRD sources | `make test` (runs `test/distribution`) | ✅ Green | `test/distribution` package ok in `make test` |
| G4 | Build | `go build ./...` succeeds | `go build ./...` | ✅ Green | `BUILD_EXIT:0` |
| G5 | Lint zero-issue | golangci-lint reports 0 issues | `golangci-lint run` | ✅ Green | `0 issues.` rc=0 |
| G6 | Envtest green | controller envtest suite passes | `make test` | ✅ Green | rc=0; controller pkgs `ok` |
| G7 | e2e compiles | `go vet ./test/e2e/... ./test/utils/...` rc=0 | `go vet ./test/e2e/... ./test/utils/...` | ✅ Green | rc=0 |

Notes:
- G1 is the slice-level lockout surface: any marker/regeneration drift flips the
  diff nonzero, failing the gate. Verified clean at HEAD.
- G2/G3 reuse the generation+distribution contracts established by T02 (marker
  lock) and the `test/distribution` contract test exercising kustomize↔Helm parity.
- G7 does **not** imply `make test-e2e` runs; it asserts the e2e package type-
  checks and the OIDC lifecycle wiring compiles. Runtime proof lives behind the
  precondition in §2.

## 2. Credential-Safe Fixture Precondition (make test-e2e)

`make test-e2e` builds a Kind cluster, loads the operator image, installs
PostgreSQL + Dependency-Track + Dex(OIDC) via Helm, and drives the live DT API.
It is the ONLY place the real-DT ownership-preservation + vanilla-Team lifecycle
is proved. Environmental prerequisites:

🧭 Precond — a container runtime (`docker` daemon reachable, or `podman`) **+** `kind` **+** `KIND_CLUSTER` exported.

| Step | Purpose | Command |
|------|---------|---------|
| Fresh cluster | provision/clean Kind | `make setup-test-e2e` (`KIND_CLUSTER=<name>`) |
| Run e2e | drive the live fixture | `make test-e2e` (= `KIND_CLUSTER=... go test ./test/e2e/ -v -ginkgo.v`) |
| Fast iterate | reuse a preserved cluster | `make test-e2e-fast` (sets `E2E_SKIP_CLUSTER_TEARDOWN=true`) |
| Teardown | destroy Kind cluster | `make cleanup-test-e2e` |

Credential safety (hardening delivered by T05):
- DT credentials, usernames, passwords, UUIDs and group names are **streamed over
  stdin** to probe scripts exec-ed into the DT API pod — never placed in argv or
  emitted to test logs.
- Probes use `curl --fail --silent --show-error` with explicit `--max-time`, so a
  flaky/unreachable backend fails loud instead of hanging the suite.
- Login failure is propagated via `|| exit 1` in `oidcLoginPreamble`.
- Intentional 404s (scrubbed mappings) are expressed as HTTP-status equality
  (`mappingHTTPStatus(...) == "404"`) so a defensive negative assertion does not
  trip `utils.Run`'s error path.
- Implementation: `test/e2e/e2e_test.go` constants `oidcLoginPreamble`,
  `oidcGroupByNameScript`, `oidcGroupTeamsScript`, `oidcMappingDeleteScript`;
  runners `runOIDCAPI` / `runPolicyAPI` stream stdin buffers.

## 3. Real Dependency-Track Lifecycle Coverage

Location: `test/e2e/e2e_test.go`, `Context("OIDC group-to-team mapping lifecycle...")`.
Three phases map to the S06 ownership-preservation contract:

- **(a) create-N** — `spec.oidc.groups=[oidc-admins, oidc-devs]` converges two
  owned mappings (`status.oidc.ownedMappings` reaches 2); both OIDC groups exist
  and both map the team. Asserts `Expect(ownedMappings()).To(HaveLen(2))`.
- **(b) shrink-revise-404** — patched to `[oidc-admins]`; the retained binding
  survives verbatim (`retained.MappingUUID == adminMap.MappingUUID`), the revoked
  `oidc-devs` mapping asserts HTTP **404** (scrubbed), and the shared group
  *object* persists (only edges are managed, not principal objects).
- **(c) delete-scrub-owner-only** — deleting the Team fires the finalizer; every
  owned mapping 404s afterward and both group objects persist.

Vanilla-Team cases (permission-less / empty-`permissions:[]` creation) are covered
by the sibling `Context("Team integration with real DependencyTrack")`.

## 4. Current Evidence (snapshot at task T05)

Quoted source revision: `HEAD` ≈ `04562d4` family (post-MEM023 mounting of
`-webhook-server-cert`, post-T02 `+kubebuilder:webhook` lock).

| Gate | Result | How | Duration |
|------|--------|-----|----------|
| G1 webhook marker | exit 0 (no drift) | `git diff --quiet config/webhook/manifests.yaml` | instant |
| G4 lint | `0 issues.` rc=0 | `golangci-lint run` | ~0.9 s |
| G5 envtest | rc=0 | `make test` | ~12.7 s |
| G6 e2e-compile | rc=0 | `go vet ./test/e2e/... ./test/utils/...` | ~0.1 s |
| G7 build | rc=0 | `go build ./...` | ~2.8 s |

Live-DT lifecycle proof remains gated on a Docker/Kind-capable runner
(see §2); the compile/contract gates above are the in-sandbox portion.

## Appendix: Quality-Gate Responses (Q5 / Q6 / Q7)

### Failure Modes

External dependencies of the hardened e2e pathway and their failure handling
(all verified by inspection of `test/e2e/e2e_test.go` + `test/utils/utils.go`):

1. **Container runtime (Docker/podman).** `dockerDaemonReachable()` returns false
   when the `docker` binary is absent or the daemon is unreachable; the loader
   then falls back to `podman save`+`kind load image-archive`. Neither present ⇒
   `LoadImageToKindClusterWithName` returns `"neither docker nor podman is
   available…"`, aborting `BeforeSuite` loudly rather than silently mis-deploying.
2. **Kind cluster.** `setup-test-e2e` destroys/recreates the named cluster; a
   mismatched `KIND_CLUSTER` name is threaded consistently to create/load/delete.
3. **kubectl / cluster reachability.** Every kubectl op is routed through
   `utils.Run`, which captures combined output and wraps errors; `Expect …
   NotTo(HaveOccurred)` fails fast with the offending command echoed to
   `GinkgoWriter`.
4. **DependencyTrack API.** `InstallDependencyTrack` polls `/api/version` for
   120×1 s and returns a timeout error on failure; `VerifyOIDCAvailable` halts the
   suite when `GET /api/v1/oidc/available ≠ true`, preventing misleading flakes.
5. **Dex OIDC provider.** `ProvisionDex` waits `--for=condition=Ready` on the Dex
   pod with a 2-minute timeout before DT boots, so the issuer is resolvable when
   DT fetches its well-known metadata.
6. **Credentials / secrets.** Decoded from `deptrack-credentials` via `json`
   Unmarshal of `kubectl get secret … -o json`; streamed over stdin to probes so
   secrets never enter argv or logs.
7. **Malformed/over-timeout API responses.** Probe scripts combine
   `curl --fail … --show-error --max-time` with exact-string `grep -Ft`; login
   failure propagates via `|| exit 1`; an intentionally-absent mapping is
   asserted as `== "404"` rather than allowed to raise a non-nil `Run` error.

### Load Profile

_N/A — T05 is a compile + documentation unit; the e2e suite is driven
interactively/by CI, not under synthetic load. There is no runtime saturation
dimension to size pools/rate-limit/cache._

### Negative Tests

Defences asserting error paths in the hardened e2e pathway:

- **`should reject a Team CR with an invalid schema`** (`test/e2e/e2e_test.go`):
  `spec.permissions: "not-an-array"` is rejected by the kubebuilder-validation
  webhook; `applyTeam(…, expectFail=true)` inverts the expectation.
- **`should reject a Policy CR with an invalid enum value`**:
  `spec.violationState: NOOP` is rejected; same inverted-apply pattern.
- **Scrub-defect guard**: `Expect(mappingHTTPStatus(devMap.MappingUUID)).To(Equal("404"))`
  after the shrink phase — a residual binding (left-over edge) fails the gate.
- **Credential-propagation negative**: `oidcLoginPreamble` terminates with
  `|| exit 1`, so a bogus login aborts the probe (login failure ⇒ non-zero exit
  ⇒ `Expect` fails), instead of producing a silent empty-token 401 cascade.
- **Unit-layer negative**: `utils.splitImage` rejects image references lacking a
  tag (`"image %q must include a tag"`), guarding the Kind image-name plumbing.
