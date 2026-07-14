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

Once Helm has been set up correctly, add the repo as follows:

helm repo add dependencytrack-operator https://martin31821.github.io/dependencytrack-operator

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages. You can then run `helm search repo
dependencytrack-operator` to see the charts.

To install the dependencytrack-operator chart:

    helm install my-dependencytrack-operator dependencytrack-operator/dependencytrack-operator

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
