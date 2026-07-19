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

- Automate initial setup (credential rotation, API keys, OIDC, Teams, ...)
- Integrate seamlessly with GitOps workflows (ArgoCD, Flux, etc.)

Note, that we currently see this as an intermediate solution until the gap is
closed in upstream DependencyTrack.

## Custom Resources

The operator provides three CRDs in the `dependencytrack.mko.dev/v1alpha1` API group.

### Team

Creates and manages a **Team** in DependencyTrack.

| Field | Type | Required | Description |
|---|---|---|---|
| `spec.name` | string | No | Human-readable team name |
| `spec.permissions` | []string | No | List of permission names to assign (omit to leave unchanged, empty array to clear all) |
| `status.uuid` | string | — | DependencyTrack UUID assigned to the team |
| `status.permissions` | string | — | Comma-separated list of permissions last synced (observability only) |
| `status.conditions` | []Condition | — | Reconciliation state |

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

### APIKey

Creates and manages an **API access key** in DependencyTrack, scoped to a Team. The generated key value is stored in a Kubernetes `Secret`.

| Field | Type | Required | Description |
|---|---|---|---|
| `spec.teamRef` | string | Yes | Name of the `Team` CR (same namespace) this key belongs to |
| `spec.secretName` | string | Yes | Kubernetes `Secret` where the generated key is stored |
| `spec.comment` | string | No | Human-readable label for the key in DependencyTrack |
| `status.publicId` | string | — | DependencyTrack's stable key identifier (for updates/deletes) |
| `status.conditions` | []Condition | — | Reconciliation state |

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

| Field | Type | Required | Description |
|---|---|---|---|
| `spec.name` | string | Yes | Human-readable policy name; must be globally unique in DependencyTrack |
| `spec.description` | string | No | Human-readable description of the policy |
| `spec.priority` | string | Yes | Policy priority: `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`, or `INFO` |
| `spec.failureAction` | string | Yes | Action for violations: `BLOCK_RELEASE`, `BLOCK_DEPLOY`, `REPORT`, or `IGNORE` |
| `spec.conditions` | []PolicyCondition | Yes | One or more inline conditions evaluated by DependencyTrack |
| `spec.conditions[].type` | string | Yes | DependencyTrack condition subject, such as `SEVERITY`, `LICENSE`, `CPE`, `PURL`, or `VULNERABILITY` |
| `spec.conditions[].comparator` | string | Yes | Comparison operator: `EQ`, `NE`, `GT`, `GTE`, `LT`, or `LTE` |
| `spec.conditions[].value` | string | Yes | Value compared by the condition |
| `status.uuid` | string | — | DependencyTrack UUID used as the authoritative remote identity |
| `status.conditions` | []Condition | — | Reconciliation state |

**Example:**

```yaml
apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Policy
metadata:
  name: critical-vulnerability-policy
  namespace: default
spec:
  name: Critical Vulnerability Policy
  description: Report components with critical vulnerabilities
  priority: CRITICAL
  failureAction: REPORT
  conditions:
    - type: SEVERITY
      comparator: EQ
      value: CRITICAL
```

The operator creates the policy first and then persists each inline condition through DependencyTrack's condition API. It records the remote UUID in `status.uuid`, uses that UUID for subsequent updates and deletion, and reports failures through the `Ready` status condition.

> **Dependency-Track v5.0.2 compatibility:** use subjects and condition behavior supported by that release. In particular, `CVSS` conditions and suppression conditions (`isSuppression: true`) are not supported by Dependency-Track v5.0.2. The aliases `PURL` and `VULNERABILITY` are translated to `PACKAGE_URL` and `VULNERABILITY_ID` respectively.

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

### Environment variables

The operator container requires these environment variables:

| Variable                      | Description                                                                                                        | Default                 |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------ | ----------------------- |
| `DEPTRACK_URL`                | HTTP(S) URL of the DependencyTrack instance to manage (e.g. `http://dtrack-service:8080`)                           | `http://dtrack-service:8080` |
| `DEPTRACK_CREDENTIALS_SECRET` | Name of the Kubernetes `Secret` that holds the `username` and `password` used to authenticate with DependencyTrack | `deptrack-credentials`  |
| `POD_NAMESPACE`               | Namespace the operator runs in (auto-injected by Kubernetes)                                                       | auto-injected           |

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

### Cert-manager (optional)

For production deployments you may want to enable cert-manager so the metrics
and webhook endpoints use CA-signed TLS certificates. Uncomment the relevant
lines in `config/default/kustomization.yaml` and `config/prometheus/kustomization.yaml`
before building the Helm chart.
