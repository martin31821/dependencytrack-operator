# Image URL to use for building/pushing the manager container image
IMG ?= controller:latest

# Container tool (docker or podman)
CONTAINER_TOOL ?= docker

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: test
test: manifests generate fmt envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell setup-envtest use $(ENVTEST_K8S_VERSION) -p path)" go test $$(go list ./... | grep -v -e /e2e -e /gen/dtapi/test -e /gen/dtv2/test) -coverprofile cover.out

# TODO(user): To use a different vendor for e2e tests, modify the setup under 'tests/e2e'.
# The default setup assumes Kind is pre-installed and builds/loads the Manager Docker image locally.
# CertManager is installed by default; skip with:
# - CERT_MANAGER_INSTALL_SKIP=true
KIND ?= kind
KIND_CLUSTER ?= deptrack-operator-test-e2e

# E2E cluster preservation (fast iteration between runs):
#   E2E_SKIP_CLUSTER_TEARDOWN=true  – keep the Kind cluster after tests
#   E2E_SKIP_DT_TEARDOWN=true       – keep DependencyTrack between suite runs
# When both are true, only teardown the operator (if the main AfterSuite runs),
# making it possible to iterate on operator code without rebuilding DT.

.PHONY: setup-test-e2e
setup-test-e2e: ## Set up an isolated Kind cluster for e2e tests (always creates fresh)
	@command -v $(KIND) >/dev/null 2>&1 || { \
		echo "Kind is not installed. Please install Kind manually."; \
		exit 1; \
	}
	@if [ "$$($(KIND) get clusters 2>/dev/null)" = "$(KIND_CLUSTER)" ]; then \
		echo "Deleting existing Kind cluster '$(KIND_CLUSTER)' for a clean run..."; \
		$(KIND) delete cluster --name $(KIND_CLUSTER) 2>/dev/null || true; \
	fi
	@echo "Creating Kind cluster '$(KIND_CLUSTER)'..."
	$(KIND) create cluster --name $(KIND_CLUSTER)

.PHONY: test-distribution
test-distribution: ## Run distribution contract tests between kustomize and Helm chart artifacts.
	go test ./test/distribution/ -v

.PHONY: test-e2e
test-e2e: setup-test-e2e manifests generate fmt ## Run the e2e tests. Expected an isolated environment using Kind.
	KIND_CLUSTER=$(KIND_CLUSTER) go test ./test/e2e/ -v -ginkgo.v

.PHONY: test-e2e-fast
test-e2e-fast: manifests generate fmt ## Run e2e tests preserving the Kind cluster (fast iteration).
	@if [ "$$($(KIND) get clusters 2>/dev/null)" != "$(KIND_CLUSTER)" ]; then \
		echo "Creating Kind cluster '$(KIND_CLUSTER)'..."; \
		$(KIND) create cluster --name $(KIND_CLUSTER); \
	fi
	E2E_SKIP_CLUSTER_TEARDOWN=true \
	E2E_SKIP_DT_TEARDOWN=true \
	KIND_CLUSTER=$(KIND_CLUSTER) \
	go test ./test/e2e/ -v -ginkgo.v

.PHONY: cleanup-test-e2e
cleanup-test-e2e: ## Tear down the Kind cluster used for e2e tests
	@$(KIND) delete cluster --name $(KIND_CLUSTER)

.PHONY: cleanup-e2e-services
cleanup-e2e-services: ## Remove DT and operator Helm releases from the preserved cluster
	helm uninstall my-dependency-track --namespace dependency-track || true
	helm uninstall deptrack-operator --namespace deptrack-operator-system || true
	kubectl delete namespace dependency-track deptrack-operator-system --ignore-not-found=true || true

.PHONY: lint
lint: ## Run golangci-lint linter
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint linter and perform fixes
	golangci-lint run --fix

.PHONY: lint-config
lint-config: ## Verify golangci-lint linter configuration
	golangci-lint config verify

##@ Build

.PHONY: build
build: manifests generate fmt ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name deptrack-operator-builder
	$(CONTAINER_TOOL) buildx use deptrack-operator-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm deptrack-operator-builder
	rm Dockerfile.cross

.PHONY: build-installer
build-installer: manifests generate ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	./hack/kustomize-build-with-image.sh "${IMG}" > dist/install.yaml

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kustomize build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	./hack/kustomize-build-with-image.sh "${IMG}" | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Dependencies

## Tool versions (used by setup-envtest to determine which K8s binaries to download)
#ENVTEST_VERSION is the version of controller-runtime release branch to fetch the envtest setup script (i.e. release-0.20)
ENVTEST_VERSION ?= $(shell go list -m -f "{{ .Version }}" sigs.k8s.io/controller-runtime | awk -F'[v.]' '{printf "release-%d.%d", $$2, $$3}')
#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries (i.e. 1.31)
ENVTEST_K8S_VERSION ?= $(shell go list -m -f "{{ .Version }}" k8s.io/api | awk -F'[v.]' '{printf "1.%d", $$3}')

.PHONY: envtest
envtest: ## Download envtest binaries (etcd, kube-apiserver) for the target K8s version.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@setup-envtest use $(ENVTEST_K8S_VERSION) -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}


##@ Helm Chart

CHART_DIR ?= deploy/charts/dependencytrack-operator
.PHONY: helm-chart
helm-chart: manifests ## Generate a Helm chart from kustomize output.
	mkdir -p $(CHART_DIR)
	./hack/kustomize-build-with-image.sh "$(IMG)" | helmify $(CHART_DIR) 2>/dev/null
	@# Remove duplicate hardcoded selector/template labels that conflict with helpers
	@sed -i '/matchLabels:/,/^    {{/{/app\.kubernetes\.io\/name: deptrack-operator/d}' $(CHART_DIR)/templates/deployment.yaml
	@sed -i '/labels:/,/^    {{/{/app\.kubernetes\.io\/name: deptrack-operator/d}' $(CHART_DIR)/templates/deployment.yaml
	@# Add imagePullPolicy support (helmify does not generate it)
	@awk '/AppVersion/{print; print "        {{- if .Values.controllerManager.manager.image.pullPolicy }}"; print "        imagePullPolicy: {{ .Values.controllerManager.manager.image.pullPolicy }}"; print "        {{- end }}"; next}1' $(CHART_DIR)/templates/deployment.yaml > $(CHART_DIR)/templates/deployment.yaml.tmp && mv $(CHART_DIR)/templates/deployment.yaml.tmp $(CHART_DIR)/templates/deployment.yaml
