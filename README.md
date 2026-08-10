# dependencytrack-operator

A Kubernetes Operator for the declarative bootstrapping, provisioning, and
lifecycle management of [OWASP Dependency-Track](https://dependencytrack.org/) instances.

## Description

Dependency-Track is an intelligent Software Composition Analysis (SCA) platform
that identifies open-source risks and vulnerabilities in software projects.
However, as noted in [Dependency-Track #6193](https://github.com/DependencyTrack/dependency-track/issues/6193),
the project lacks a built-in, GitOps-friendly mechanism to automate the initial
bootstrapping and provisioning of a fresh installation.

`dependencytrack-operator` closes this gap by providing a **Kubernetes-native,
declarative approach** to:

- Automate initial setup (credential rotation, API keys, OIDC, Teams, notification publishers, notification rules, policies, ...)
- Integrate seamlessly with GitOps workflows (ArgoCD, Flux, etc.)

Note, that we currently see this as an intermediate solution until the gap is
closed in upstream DependencyTrack.

## Custom Resources

The operator provides five CRDs in the `dependencytrack.mko.dev/v1alpha1` API group.

### Team

Creates and manages a **Team** in DependencyTrack.

| Field                       | Type               | Required | Description                                                                                                                                                                                                                                                                                                             |
| --------------------------- | ------------------ | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `spec.name`                 | string             | No       | Human-readable team name                                                                                                                                                                                                                                                                                                |
| `spec.permissions`          | []string           | No       | List of permission names to assign (omit to leave unchanged, empty array to clear all)                                                                                                                                                                                                                                  |
| `spec.oidc`                 | object             | No       | Optional OIDC group-mapping configuration. A nil value disables OIDC management (zero API traffic). See [OpenID Connect (OIDC)](#openid-connect-oidc-group-to-team-mapping) below.                                                                                                                                      |
| `spec.oidc.groups`          | []string           | No       | Ordered list of OIDC group names (Identity Provider claim values) to bind to this team. Compared verbatim—case is preserved; omit to leave existing mappings unchanged, empty array to clear all mappings. Validated by the admission webhook (blank/whitespace-only and trimmed-exact-duplicate entries are rejected). |
| `status.uuid`               | string             | —        | DependencyTrack UUID assigned to the team                                                                                                                                                                                                                                                                               |
| `status.permissions`        | string             | —        | Comma-separated list of permissions last synced (observability only)                                                                                                                                                                                                                                                    |
| `status.oidc`               | object             | —        | Observed OIDC mapping state; `nil` before first reconciliation, empty `ownedMappings` afterward.                                                                                                                                                                                                                        |
| `status.oidc.ownedMappings` | []OwnedOIDCMapping | —        | Bindings this operator created/claims: `groupName`, `groupUuid`, `teamUuid`, `mappingUuid`. Treated as the authoritative diff anchor.                                                                                                                                                                                   |
| `status.conditions`         | []Condition        | —        | Reconciliation state                                                                                                                                                                                                                                                                                                    |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: my-team
  namespace: default
spec:
  name: My Team
  permissions:
    - PORTFOLIO_MANAGEMENT
    - VIEW_PORTFOLIO
```

An OIDC-mapped team differs only by the `spec.oidc.groups` stanza. While it is
present, the controller ensures every named group is mapped to this team in
Dependency-Track (see [OpenID Connect (OIDC)](#openid-connect-oidc-group-to-team-mapping)):

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: my-team
  namespace: default
spec:
  name: My Team
  permissions:
    - PORTFOLIO_MANAGEMENT
    - VIEW_PORTFOLIO
  oidc:
    groups:
      - admins
      - developers
```

### APIKey

Creates and manages an **API access key** in DependencyTrack, scoped to a Team. The generated key value is stored in a Kubernetes `Secret`.

| Field               | Type        | Required | Description                                                   |
| ------------------- | ----------- | -------- | ------------------------------------------------------------- |
| `spec.teamRef`      | string      | Yes      | Name of the `Team` CR (same namespace) this key belongs to    |
| `spec.secretName`   | string      | Yes      | Kubernetes `Secret` where the generated key is stored         |
| `spec.comment`      | string      | No       | Human-readable label for the key in DependencyTrack           |
| `status.publicId`   | string      | —        | DependencyTrack's stable key identifier (for updates/deletes) |
| `status.conditions` | []Condition | —        | Reconciliation state                                          |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: APIKey
metadata:
  name: my-api-key
  namespace: default
spec:
  teamRef: my-team
  secretName: my-team-api-key
  comment: "CI/CD pipeline key"
```

After reconciliation, the operator creates a `Secret` with the API key value. The `Team` must exist before the `APIKey` is reconciled — the controller references the `Team` by name to create the key under that team in DependencyTrack.

### Policy

Creates and manages a global **Policy** and its conditions in DependencyTrack. The Kubernetes resource is namespaced, but DependencyTrack policies are global; policy names must therefore be unique across all namespaces managed by the operator.

| Field                        | Type              | Required | Description                                                                                          |
| ---------------------------- | ----------------- | -------- | ---------------------------------------------------------------------------------------------------- |
| `spec.name`                  | string            | Yes      | Human-readable policy name; must be globally unique in DependencyTrack                               |
| `spec.operator`              | string            | Yes      | Condition matching mode: `ANY` if one condition must match, or `ALL` if every condition must match   |
| `spec.violationState`        | string            | Yes      | Dependency-Track violation state: `INFO` (Inform), `WARN` (Warn), or `FAIL` (Fail)                   |
| `spec.conditions`            | []PolicyCondition | Yes      | One or more inline conditions evaluated by DependencyTrack                                           |
| `spec.conditions[].subject`  | string            | Yes      | Dependency-Track subject, such as `SEVERITY`, `LICENSE`, `CPE`, `PACKAGE_URL`, or `VULNERABILITY_ID` |
| `spec.conditions[].operator` | string            | Yes      | Comparison operator: `IS` or `IS_NOT`                                                                |
| `spec.conditions[].value`    | string            | Yes      | Value compared against the subject                                                                   |
| `status.uuid`                | string            | —        | DependencyTrack UUID used as the authoritative remote identity                                       |
| `status.conditions`          | []Condition       | —        | Reconciliation state                                                                                 |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Policy
metadata:
  name: critical-vulnerability-policy
  namespace: default
spec:
  name: Critical Vulnerability Policy
  operator: ANY
  violationState: WARN
  conditions:
    - subject: SEVERITY
      operator: IS
      value: CRITICAL
```

The operator creates the policy first and then persists each inline condition through DependencyTrack's condition API. It records the remote UUID in `status.uuid`, uses that UUID for subsequent updates and deletion, and reports failures through the `Ready` status condition.

> **Dependency-Track v5.0.2 compatibility:** condition subjects use Dependency-Track's native names. `CVSS` and suppression conditions are not supported; use a supported subject such as `SEVERITY`, `LICENSE`, `PACKAGE_URL`, or `VULNERABILITY_ID`.

### NotificationPublisher

Creates and manages a **notification publisher** in DependencyTrack — a configurable endpoint (Slack, email, webhook, etc.) that receives notification events.

| Field                   | Type        | Required | Description                                                                                                                                                                                        |
| ----------------------- | ----------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `spec.name`             | string      | Yes      | Display name for the publisher in DependencyTrack                                                                                                                                                  |
| `spec.extensionName`    | string      | Yes      | Publisher extension identifier (e.g. `slack`, `email`, `webhook`, `opsgenie`)                                                                                                                      |
| `spec.description`      | string      | No       | Human-readable description (max 1024 chars)                                                                                                                                                        |
| `spec.template`         | string      | No       | Custom notification message template applied by the publisher. Unlike the legacy per-rule `message`, this is unconstrained in length. When omitted, DependencyTrack applies the extension default. |
| `spec.templateMimeType` | string      | No       | Media type of the template body (e.g. `text/plain`, `text/html`, `application/json`); max 255 chars. When omitted, DependencyTrack applies the extension default.                                  |
| `status.uuid`           | string      | —        | DependencyTrack UUID assigned to the publisher                                                                                                                                                     |
| `status.name`           | string      | —        | Name last synced to DependencyTrack                                                                                                                                                                |
| `status.conditions`     | []Condition | —        | Reconciliation state                                                                                                                                                                               |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: NotificationPublisher
metadata:
  name: slack-publisher
  namespace: default
spec:
  name: Slack Notifications
  extensionName: slack
  description: "Publishes critical vulnerability events to #security-alerts"
```

A custom notification body is configured on the publisher via `template` and `templateMimeType`. This is the supported mechanism for bespoke message bodies and is unconstrained in length (the legacy per-rule `message` field was capped at 4096 chars and is not modeled by this operator):

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: NotificationPublisher
metadata:
  name: webhook-publisher
  namespace: default
spec:
  name: Webhook Publisher
  extensionName: webhook
  template: '{"text":"{{project.name}}: {{notification.subject}}"}'
  templateMimeType: application/json
```

The publisher must exist and be `Ready` before any `NotificationRule` can reference it.

### NotificationRule

Creates and manages a **notification rule** in DependencyTrack — a policy that routes notification events to a configured publisher.

| Field                           | Type        | Required | Description                                                                                  |
| ------------------------------- | ----------- | -------- | -------------------------------------------------------------------------------------------- |
| `spec.name`                     | string      | Yes      | Display name for the rule (max 255 chars)                                                    |
| `spec.scope`                    | string      | Yes      | Applies to: `SYSTEM` or `PORTFOLIO`                                                          |
| `spec.triggerType`              | string      | Yes      | Fires on: `EVENT` or `SCHEDULE`                                                              |
| `spec.level`                    | string      | Yes      | Filter by severity: `INFORMATIONAL`, `WARNING`, or `ERROR`                                   |
| `spec.publisherRef.name`        | string      | Yes      | Name of the `NotificationPublisher` CR in the same namespace                                 |
| `spec.enabled`                  | bool        | No       | Whether the rule is active (default: `true`)                                                 |
| `spec.notifyOn`                 | []string    | No       | Event types that trigger the rule (e.g. `NEW_VULNERABILITY`, `VULNERABILITY_SCAN_COMPLETED`) |
| `spec.filterExpression`         | string      | No       | QL filter string for the rule (max 1024 chars)                                               |
| `spec.publisherConfigSecretRef` | object      | No       | Secret containing publisher-specific config JSON (see below)                                 |
| `spec.logSuccessfulPublish`     | bool        | No       | Log successful publishes; defaults to false                                                  |
| `spec.notifyChildren`           | bool        | No       | Apply to child projects (only for PORTFOLIO/SYSTEM scope)                                    |
| `spec.scheduleCron`             | string      | No       | Cron expression for scheduled rules; required when `triggerType: SCHEDULE`                   |
| `spec.scheduleSkipUnchanged`    | bool        | No       | Skip emitting notifications if result is unchanged (schedule only)                           |
| `spec.teams`                    | []string    | No       | Team CR names whose remote UUID is associated with this rule                                 |
| `spec.projects`                 | []string    | No       | Project UUIDs to associate with this rule (ignored for PORTFOLIO/SYSTEM scope)               |
| `status.uuid`                   | string      | —        | DependencyTrack UUID assigned to the rule                                                    |
| `status.name`                   | string      | —        | Name last synced to DependencyTrack                                                          |
| `status.conditions`             | []Condition | —        | Reconciliation state                                                                         |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: NotificationRule
metadata:
  name: critical-vuln-rule
  namespace: default
spec:
  name: Critical Vulnerability Alert
  scope: PORTFOLIO
  triggerType: EVENT
  level: ERROR
  publisherRef:
    name: slack-publisher
  notifyOn:
    - NEW_VULNERABILITY
    - VULNERABILITY_SCAN_COMPLETED
```

**Publisher config:** Some extensions (like Slack) require configuration (webhook URL, channel, etc.). Store this JSON in a Kubernetes `Secret` and reference it:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: slack-config
  namespace: default
type: Opaque
stringData:
  config.json: '{"webhookUrl": "https://hooks.slack.com/services/T00000/B00000/XXXX", "channel": "#alerts"}'
---
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: NotificationRule
metadata:
  name: slack-alert-rule
  namespace: default
spec:
  name: Slack Critical Alert
  scope: PORTFOLIO
  triggerType: EVENT
  level: ERROR
  publisherRef:
    name: slack-publisher
  publisherConfigSecretRef:
    name: slack-config
    key: config.json
```

The operator validates the JSON config against the publisher extension schema and reports failures via the `Ready` status condition.

## OpenID Connect (OIDC) Group-to-Team Mapping

The operator integrates with Dependency-Track's
[OIDC group mapping](https://docs.dependencytrack.org/integrations/oauth2_openid_connect/)
feature so that a `Team`'s `spec.oidc.groups` declaratively drives which Identity
Provider (IdP) groups are bound to that team. Mapping is **edge-oriented**: the
operator manages the OIDC _bindings_ (group↔team), never the group _objects_,
so groups shared across teams or authored by the platform survive a team
shrink or deletion.

### Prerequisites

Dependency-Track must have OIDC provisioned before the operator can wire
mappings. When installing Dependency-Track (e.g. via its
[`dependency-track`](https://github.com/DependencyTrack/helm-charts) Helm chart),
enable OIDC on the API server by injecting these environment variables. They
use the single-`O` spelling recognized by Dependency-Track:

| Variable          | Purpose                                                             | Example                       |
| ----------------- | ------------------------------------------------------------------- | ----------------------------- |
| `DT_OIDC_ENABLED` | Enables the `/v1/oidc/*` endpoints. Must be `"true"`.               | `true`                        |
| `DT_OIDC_ISSUER`  | Base URL of the OIDC Issuer (must expose a well-known JWKS config). | `https://dex.example.com/dex` |

Through the Helm chart these are injected via `apiServer.extraEnv` as a standard
Kubernetes `EnvVar` list:

```yaml
apiServer:
  extraEnv:
    - name: DT_OIDC_ENABLED
      value: "true"
    - name: DT_OIDC_ISSUER
      value: https://dex.example.com/dex
```

Any compliant OIDC issuer works (Dex, Keycloak, Okta, Auth0, Azure AD, …). The
[e2e fixture](../test/utils/utils.go) ships a minimal
[ghcr.io/dexidp/dex](https://github.com/dexidp/dex) provider as a reference
IdP backed by its builtin local connector.

Once `DT_OIDC_ENABLED=true` and the issuer responds, Dependency-Track answers
`GET /v1/oidc/available` with `true`; the operator probes this before issuing
any mutation and short-circuits gracefully (see [Availability](#availability))
when it does not.

### How reconciliation works

When `spec.oidc` is non-nil and the operator obtains an authenticated
Dependency-Track client, it diffs the **desired** group list against the
**previously-owned** bindings recorded in `status.oidc.ownedMappings` and
converges:

1. **Dedupe** the desired list (blanks dropped, insertion-order preserved).
2. **Diff** against `status.oidc.ownedMappings` keyed by `GroupName`:
   - _Survivors_ (still desired) are carried over verbatim—no server round-trip,
     so reordering a list never churns the cluster.
   - _Stale_ edges (no longer desired) are pruned individually by mapping UUID
     via `DELETE /v1/oidc/mapping/{uuid}`. Group objects are **never** deleted.
   - _New_ names are established by looking up an existing group by exact-case
     name, or creating it, then `PUT /v1/oidc/mapping` to bind it to the team.
3. **Persist** the resulting binding set to `status.oidc.ownedMappings` as the
   authoritative diff anchor for the next reconcile.

Concretely, shrinking `developers` out of a team deletes only that mapping edge;
the `developers` group object remains. Deleting the `Team` CR scrubs every
binding the operator owns (again, by mapping UUID) before removing the team
principal—so shared groups outlive any single team.

#### Availability

`GET /v1/oidc/*` requires an OAuth2/session bearer token; Dependency-Track refuses
API-key auth on those routes (returning `401`). The controller therefore threads
the authenticated context obtained from `ClientProvider.Get` into every OIDC
call—not the ambient reconcile context, which carries no token and would make
every probe fail with `401`.

When `IsAvailable` reports `false` (OIDC not yet provisioned), the operator
makes **no** mutations and leaves `status.oidc` untouched so the last-known
ownership survives to retry. Upserts requeue; deletions defer with the
finalizer held until OIDC becomes available.

#### Status fields

| Path                                      | Meaning                                                  |
| ----------------------------------------- | -------------------------------------------------------- |
| `status.oidc`                             | `nil` until first reconciliation; populated afterward.   |
| `status.oidc.ownedMappings[].groupName`   | IdPs group claim, stored verbatim (case preserved).      |
| `status.oidc.ownedMappings[].groupUuid`   | Server-issued group UUID.                                |
| `status.oidc.ownedMappings[].teamUuid`    | Server-issued team UUID.                                 |
| `status.oidc.ownedMappings[].mappingUuid` | Server-issued binding UUID; the only thing ever deleted. |

Inspect it with:

```sh
kubectl get team <name> -o jsonpath='{.status.oidc}'
```

#### Idempotent restore (wipe-safe)

If `status.oidc` is ever cleared (operator upgrade resetting the status
subresource, or a never-reconciled team), the controller enters
restore-only mode: it (re)creates every desired binding via
lookup-or-create-then-bind, but emits **zero** `DELETE` calls—there is no
trusted prior ownership to difference against, so pruning cannot erase
peer-authored or platform-managed bindings. Adoption races (two actors creating
the same group concurrently) resolve idempotently: a `409 Conflict` on create
triggers a relist-and-adopt-by-exact-case-name, never a clobber.

#### Admission validation

A fail-closed validating webhook (`team.validating.dependencytrack.mko.dev`)
guards `spec.oidc.groups` before it reaches the controller:

- A `nil` `spec.oidc` disables OIDC management (no-op).
- A non-nil config with a `nil`/empty `groups` slice is **valid**—it is an
  intentional "bind this team to zero groups," converged to clearing by the
  controller.
- Each entry is judged on its trimmed form: blank/whitespace-only entries are
  rejected, and trimmed-exact-duplicate entries collide (but `admins` ≠ `ADMINS`—
  comparison is case-sensitive; only identically-spelled trimmed values clash).

Original casing is preserved on the wire so operators retain full control over
matching semantics against upstream IdP claims.

### Observability

The controller emits Kubernetes Events reflecting OIDC outcomes:

| Reason               | Trigger                                                   |
| -------------------- | --------------------------------------------------------- |
| `OIDCMappingCreated` | Mappings ensured (carries the count).                     |
| `OIDCMappingSkipped` | Nothing to create (empty/zero desired).                   |
| `OIDCUnavailable`    | Warning: OIDC not provisioned; mappings left untouched.   |
| `OIDCError`          | Warning: reconciliation faulted (details in the message). |
| `OIDCMappingRemoved` | Finalizer scrubbed owned mappings during team deletion.   |

Failures also surface on the `Ready` condition (`reason=OIDCUnavailable` or
`reason=OIDCError`) with the underlying error in the message.

## Getting Started

### Prerequisites

- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- helm version v3.0+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

#### Using the provided helm chart

[Helm](https://helm.sh) must be installed to use the charts. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

The chart is published as an [OCI artifact](https://helm.sh/docs/topics/registries/) on GHCR:

    helm install my-dependencytrack-operator oci://ghcr.io/martin31821/charts/dependencytrack-operator

To install a specific version:

    helm install my-dependencytrack-operator oci://ghcr.io/martin31821/charts/dependencytrack-operator --version 0.1.5

To uninstall the chart:

    helm uninstall my-dependencytrack-operator

#### From scratch

A Helm chart is provided under `deploy/charts/dependencytrack-operator`. Regenerate it from the Kustomize
output with:

```sh
make helm-chart IMG=<some-registry>/deptrack-operator:tag
```

This uses [helmify](https://github.com/arttor/helmify) to convert the
Kustomize output into a Helm chart. You can also run it directly:

```sh
kustomize build config/default | helmify deploy/charts/dependencytrack-operator
```

To install the operator from the chart:

```sh
helm install deptrack-operator ./deploy/charts/dependencytrack-operator \
  --set controllerManager.manager.image.repository=<your-registry>/deptrack-operator \
  --set controllerManager.manager.image.tag=v0.0.1
```

**NOTE:** After modifying Kustomize manifests, regenerate the Helm chart by
running `make helm-chart` again. The chart templates are auto-generated — any
manual changes to `deploy/charts/dependencytrack-operator` will be overwritten on regeneration. Preserve
custom values in `values.yaml` overrides or apply them via `helm install --values`.

### Updating CRDs (Upgrade Guide)

The CRDs are **generated into a `crds/` directory**
(`deploy/charts/dependencytrack-operator/crds/{team,policy,apikey,notificationpublisher,notificationrule}-crd.yaml`)
rather than into `templates/`, via `helmify -crd-dir` in the `helm-chart`
Makefile target. Per Helm best practice, files in `crds/` are installed once
(on `helm install`) and are **never** reconciled by `helm upgrade` nor deleted
by `helm uninstall`. The canonical source of truth remains `config/crd/bases/`
(all five CRD bases are registered in `config/crd/kustomization.yaml`); the
chart is reproduced from it by `make helm-chart`.

#### Procedure

1. Edit the desired CRD definition in `config/crd/bases/`.
   Definitions originate from the `// +kubebuilder:object:root=true` types in
   `api/v1alpha1/*_types.go`; regenerate rather than hand-editing:
   ```sh
   make manifests              # refreshes config/crd/bases via controller-gen
   make helm-chart IMG=<your-registry>/deptrack-operator:<tag>   # writes crds/*-crd.yaml
   ```
2. **(Recommended)** Back up existing CRs before the schema widens, so a
   tighter schema cannot strand live data:
   ```sh
   for c in teams policies apikeys notificationpublishers notificationrules; do
     kubectl get "$c.dependencytrack.mko.dev" -A \
       -o yaml > "backup-$c-$(date +%F).yaml"
   done
   ```
3. Apply the **updated CRDs first**, ahead of the operator image. Because
   `crds/` is ignored by `helm upgrade`, you must apply them out of band.
   An older controller against a newly-tightened schema can silently drop
   unknown fields; a brand-new controller against a stale CRD skips new
   validation. Wait for each CRD's `Established` condition before cutting over:
   ```sh
   kubectl apply -f deploy/charts/dependencytrack-operator/crds/
   ```
4. Roll the operator Deployment:
   ```sh
   helm upgrade my-dependencytrack-operator \
     ./deploy/charts/dependencytrack-operator \
     --set controllerManager.manager.image.tag=<tag> --reuse-values
   kubectl -n deptrack-operator-system rollout status \
     deploy/deptrack-operator-controller-manager
   ```

#### Things that break silent upgrades

- **Removing required fields** that existing CRs populate flips the CRD to
  `Accepted=False` ("field is required"). Strip the field from offending CRs
  first, or re-add it before deprecating.
- **Changing `spec.versions`.** Deprecating a version requires migrating every
  stored object away from it first, then patching `status.storedVersions`.
- **Structural-schema pruning** drops unknown fields unless
  `x-kubernetes-preserve-unknown-fields: true`. Stick to additive-only edits.

#### Verify

```sh
kubectl get crd policies.dependencytrack.mko.dev \
  -o jsonpath='{.status.conditions[*].type}={.status.conditions[*].status}{"\n"}'
kubectl -n deptrack-operator-system get pods
kubectl describe policy/<sample>   # observe a fresh status.conditions stamp
kubectl logs deploy/deptrack-operator-controller-manager --tail=50
```

### Environment variables

The operator container requires these environment variables:

| Variable                      | Description                                                                                                        | Default                      |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------ | ---------------------------- |
| `DEPTRACK_URL`                | HTTP(S) URL of the DependencyTrack instance to manage (e.g. `http://dtrack-service:8080`)                          | `http://dtrack-service:8080` |
| `DEPTRACK_CREDENTIALS_SECRET` | Name of the Kubernetes `Secret` that holds the `username` and `password` used to authenticate with DependencyTrack | `deptrack-credentials`       |
| `POD_NAMESPACE`               | Namespace the operator runs in (auto-injected by Kubernetes)                                                       | auto-injected                |

The credentials `Secret` must contain two keys:

| Key        | Description                                                                                                                                          |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| `username` | Admin username for DependencyTrack (typically `admin`)                                                                                               |
| `password` | Password for that user. The operator automatically rotates weak passwords (fewer than 30 characters) to a cryptographically random value on startup. |

**Bootstrapping with a fresh DependencyTrack instance:** If the credentials
`Secret` does not exist, the operator creates one with `admin:admin` and then
immediately rotates the password in DependencyTrack via the `ForceChangePassword`
API. After the first successful rotation, the operator never reverts to the
hard-coded defaults.

### Helm configuration reference

When deploying with the provided Helm chart you can override the defaults
via `--set` flags or a custom `values.yaml` file:

```yaml
controllerManager:
  manager:
    env:
      deptrackUrl: https://dtrack.example.com # your DependencyTrack URL
      deptrackCredentialsSecret: deptrack-credentials # secret name (unchanged)
    image:
      repository: ghcr.io/your-org/dependencytrack-operator
      tag: v1.0.0
    resources:
      limits:
        cpu: 500m
        memory: 128Mi
      requests:
        cpu: 10m
        memory: 64Mi
  replicas: 2 # set > 1 for HA leader election
```

## End-to-End Tests

The e2e suite (`test/e2e`) drives an ephemeral [Kind](https://kind.sigs.k8s.io/)
cluster: it builds the manager image with `make docker-build`, pushes it into
the Kind node with `kind load docker-image`, installs DependencyTrack (with Dex
as the OIDC issuer) and the operator via Helm, then reconciles the sample CRs.

The `Makefile` defaults to **Podman**:

```makefile
CONTAINER_TOOL ?= podman
```

### `DOCKER_HOST` requirement when using Podman

Building the image goes through Podman's native CLI, so it works regardless of
`DOCKER_HOST`. Loading the freshly built image **into the Kind cluster** is
where `DOCKER_HOST` becomes mandatory: `kind load docker-image` speaks the
Docker Engine API, and it reaches the runtime exclusively through whatever
`DOCKER_HOST` advertises (falling back to the default Docker socket
`/var/run/docker.sock`, which does not exist under a pure-Podman setup). If
`DOCKER_HOST` is unset or points nowhere, the build succeeds but the
load step stalls or fails, leaving the Kind node unable to schedule the operator
because the image was never copied.

Expose Podman's Docker-compatible REST socket and point `DOCKER_HOST` at it
**before** invoking `make test-e2e`:

**Linux (per-user socket):**

```sh
systemctl --user start podman.socket   # Fedora/RHEL
export DOCKER_HOST=unix:///run/user/$(id -u)/podman/podman.sock
```

**Linux (system socket):**

```sh
sudo systemctl start podman.socket
export DOCKER_HOST=unix:///run/podman/podman.sock
```

**macOS / Windows (Podman Machine):** `podman machine init && podman machine
start` exports `DOCKER_HOST` for you; otherwise set it to the forwarded TCP
endpoint printed by `podman env`.

Verify connectivity before running the suite:

```sh
podman info >/dev/null && curl --unix-socket "${DOCKER_HOST#unix://}" version
```

Then run the tests normally:

```sh
make test-e2e
```

For faster iteration without rebuilding DependencyTrack between runs, preserve
the cluster and the DT stack:

```sh
E2E_SKIP_CLUSTER_TEARDOWN=true E2E_SKIP_DT_TEARDOWN=true make test-e2e-fast
```

### Cert-manager (optional)

For production deployments you may want to enable cert-manager so the metrics
and webhook endpoints use CA-signed TLS certificates. Uncomment the relevant
lines in `config/default/kustomization.yaml` and `config/prometheus/kustomization.yaml`
before building the Helm chart.
