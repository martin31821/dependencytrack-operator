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

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2" // nolint:revive,staticcheck
)

const (
	prometheusOperatorVersion = "v0.77.1"
	prometheusOperatorURL     = "https://github.com/prometheus-operator/prometheus-operator/" +
		"releases/download/%s/bundle.yaml"
)

var (
	// postgresqlManifest is the YAML for a lightweight PostgreSQL
	// deployment. DependencyTrack expects a service named "postgresql".
	postgresqlManifest = `apiVersion: v1
kind: Namespace
metadata:
  name: dependency-track
---
apiVersion: v1
kind: Secret
metadata:
  name: dt-postgres
  namespace: dependency-track
type: Opaque
stringData:
  POSTGRES_DB: dtrack
  POSTGRES_USER: dtrack
  POSTGRES_PASSWORD: dtrack123
---
apiVersion: v1
kind: Service
metadata:
  name: postgresql
  namespace: dependency-track
  labels:
    app: postgresql
spec:
  ports:
    - port: 5432
      targetPort: 5432
      protocol: TCP
  selector:
    app: postgresql
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgresql
  namespace: dependency-track
  labels:
    app: postgresql
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgresql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: postgresql
    spec:
      containers:
        - name: postgres
          image: postgres:16
          envFrom:
            - secretRef:
                name: dt-postgres
          ports:
            - containerPort: 5432
              name: postgres
          readinessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - dtrack
            initialDelaySeconds: 10
            periodSeconds: 5
          livenessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - dtrack
            initialDelaySeconds: 30
            periodSeconds: 10
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              memory: 512Mi
      volumes:
        - name: postgres-data
          emptyDir: {}
          `
)

// dexManifest provisions a minimal Dex OIDC provider in the dependency-track
// namespace. Dex uses in-memory storage and its builtin local connector,
// registering DependencyTrack as a static OIDC client so the e2e fixture can
// exercise the full OIDC group-to-team mapping lifecycle (S05).
var dexManifest = `apiVersion: v1
kind: ConfigMap
metadata:
  name: dex-config
  namespace: dependency-track
data:
  config.yaml: |
    issuer: http://dex.dependency-track.svc.cluster.local:5556/dex
    storage:
      type: memory
    web:
      http: 0.0.0.0:5556
    staticClients:
      - id: dependency-track
        redirectURIs:
          - 'http://my-dependency-track-api-server.dependency-track.svc.cluster.local:8080/callback'
        name: 'Dependency Track'
        public: true
    enablePasswordDB: true
    staticPasswords:
      - email: "admin@example.com"
        username: "admin"
        userID: "08a8684b-db88-4b73-90a9-3cd1661f5466"
        hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
    oauth2:
      passwordConnector: local
---
apiVersion: v1
kind: Service
metadata:
  name: dex
  namespace: dependency-track
  labels:
    app: dex
spec:
  ports:
    - port: 5556
      targetPort: 5556
      protocol: TCP
  selector:
    app: dex
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dex
  namespace: dependency-track
  labels:
    app: dex
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dex
  template:
    metadata:
      labels:
        app: dex
    spec:
      containers:
        - name: dex
          image: ghcr.io/dexidp/dex:v2.45.1
          command:
            - /usr/local/bin/dex
          args:
            - serve
            - /config/config.yaml
          ports:
            - containerPort: 5556
              name: dex
          readinessProbe:
            tcpSocket:
              port: 5556
            initialDelaySeconds: 3
            periodSeconds: 3
          volumeMounts:
            - name: config
              mountPath: /config
      volumes:
        - name: config
          configMap:
            name: dex-config
`

func warnError(err error) {
	_, _ = fmt.Fprintf(GinkgoWriter, "warning: %v\n", err)
}

// Run executes the provided command within this context
func Run(cmd *exec.Cmd) (string, error) {
	dir, _ := GetProjectDir()
	cmd.Dir = dir

	if err := os.Chdir(cmd.Dir); err != nil {
		_, _ = fmt.Fprintf(GinkgoWriter, "chdir dir: %q\n", err)
	}

	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	command := strings.Join(cmd.Args, " ")
	_, _ = fmt.Fprintf(GinkgoWriter, "running: %q\n", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("%q failed with error %q: %w", command, string(output), err)
	}

	return string(output), nil
}

// InstallPrometheusOperator installs the prometheus Operator to be used to export the enabled metrics.
func InstallPrometheusOperator() error {
	url := fmt.Sprintf(prometheusOperatorURL, prometheusOperatorVersion)
	cmd := exec.Command("kubectl", "create", "-f", url)
	_, err := Run(cmd)
	return err
}

// UninstallPrometheusOperator uninstalls the prometheus
func UninstallPrometheusOperator() {
	url := fmt.Sprintf(prometheusOperatorURL, prometheusOperatorVersion)
	cmd := exec.Command("kubectl", "delete", "-f", url)
	if _, err := Run(cmd); err != nil {
		warnError(err)
	}
}

// IsPrometheusCRDsInstalled checks if any Prometheus CRDs are installed
// by verifying the existence of key CRDs related to Prometheus.
func IsPrometheusCRDsInstalled() bool {
	// List of common Prometheus CRDs
	prometheusCRDs := []string{
		"prometheuses.monitoring.coreos.com",
		"prometheusrules.monitoring.coreos.com",
		"prometheusagents.monitoring.coreos.com",
	}

	cmd := exec.Command("kubectl", "get", "crds", "-o", "custom-columns=NAME:.metadata.name")
	output, err := Run(cmd)
	if err != nil {
		return false
	}
	crdList := GetNonEmptyLines(output)
	for _, crd := range prometheusCRDs {
		for _, line := range crdList {
			if strings.Contains(line, crd) {
				return true
			}
		}
	}

	return false
}

// LoadImageToKindClusterWithName loads a local container image into the kind
// cluster. It prefers the docker-aware "kind load docker-image" path when a
// Docker daemon is reachable; when Docker is unavailable (hosts relying on
// podman exclusively), it falls back to exporting the image with
// "podman save" and importing it via "kind load image-archive". The latter
// sidesteps kind's experimental podman provider, which cannot enumerate images
// stored in podman's local repository.
func LoadImageToKindClusterWithName(name string) error {
	cluster := "kind"
	if v, ok := os.LookupEnv("KIND_CLUSTER"); ok {
		cluster = v
	}

	if dockerDaemonReachable() {
		kindOptions := []string{"load", "docker-image", name, "--name", cluster}
		_, err := Run(exec.Command("kind", kindOptions...))
		return err
	}

	podmanBin, err := exec.LookPath("podman")
	if err != nil {
		return fmt.Errorf("neither docker nor podman is available to load image %q into kind cluster %q", name, cluster)
	}
	return loadImageViaArchive(podmanBin, name, cluster)
}

// dockerDaemonReachable reports whether the docker CLI can contact a daemon.
// Absence of the binary or an unreachable daemon yields false so callers can
// fall back to an alternate runtime.
func dockerDaemonReachable() bool {
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	cmd := exec.Command("docker", "info")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

// loadImageViaArchive exports the image with the podman CLI to a tarball and
// feeds it to "kind load image-archive", bypassing kind's experimental podman
// provider which cannot see podman-local images.
func loadImageViaArchive(podmanBin, name, cluster string) error {
	archive, err := os.CreateTemp("", "kind-image-*.tar")
	if err != nil {
		return fmt.Errorf("create temp archive: %w", err)
	}
	archivePath := archive.Name()
	defer func() { _ = os.Remove(archivePath) }()
	if err := archive.Close(); err != nil {
		return fmt.Errorf("close temp archive: %w", err)
	}

	if _, err := Run(exec.Command(podmanBin, "save", "-o", archivePath, name)); err != nil {
		return fmt.Errorf("podman save %s: %w", name, err)
	}

	if _, err := Run(exec.Command("kind", "load", "image-archive", archivePath, "--name", cluster)); err != nil {
		return fmt.Errorf("kind load image-archive %s: %w", archivePath, err)
	}
	return nil
}

// GetNonEmptyLines converts given command output string into individual objects
// according to line breakers, and ignores the empty elements in it.
func GetNonEmptyLines(output string) []string {
	var res []string
	elements := strings.Split(output, "\n")
	for _, element := range elements {
		if element != "" {
			res = append(res, element)
		}
	}

	return res
}

// GetProjectDir will return the directory where the project is
func GetProjectDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return wd, fmt.Errorf("failed to get current working directory: %w", err)
	}
	wd = strings.ReplaceAll(wd, "/test/e2e", "")
	return wd, nil
}

// GenerateKEK returns a random base64-encoded 32-byte key suitable for
// DependencyTrack's secretManagement.database.kek.value setting.
func GenerateKEK() (string, error) {
	cmd := exec.Command("openssl", "rand", "-base64", "32")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("openssl rand: %w", err)
	}
	return string(output), nil
}

// ProvisionDex deploys a Dex OIDC provider (Deployment + ClusterIP Service +
// ConfigMap) into the dependency-track namespace. Dex uses in-memory storage
// and its builtin local connector, registering DependencyTrack as a static
// OIDC client.
// The issuer URL is http://dex.dependency-track.svc.cluster.local:5556/dex.
// Called from InstallDependencyTrack so the issuer is resolvable before DT
// fetches its well-known OIDC metadata.
func ProvisionDex() error {
	_, _ = fmt.Fprintf(GinkgoWriter, "Provisioning Dex OIDC provider...\\n")
	tmpFile, err := os.CreateTemp("", "dex-*.yaml")
	if err != nil {
		return fmt.Errorf("create dex temp file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()
	if _, err := tmpFile.WriteString(dexManifest); err != nil {
		return fmt.Errorf("write dex manifest: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close dex manifest: %w", err)
	}
	if _, err := Run(exec.Command("kubectl", "apply", "-f", tmpFile.Name())); err != nil {
		return fmt.Errorf("apply dex manifest: %w", err)
	}
	// Wait for Dex to become ready so the issuer is resolvable before DT boots.
	_, _ = fmt.Fprintf(GinkgoWriter, "Waiting for Dex to be Ready...\\n")
	waitCmd := exec.Command("kubectl", "wait", "--for=condition=Ready",
		"pod", "-l", "app=dex",
		"-n", "dependency-track",
		"--timeout", "2m")
	if _, err := Run(waitCmd); err != nil {
		return fmt.Errorf("timeout waiting for Dex pod: %w", err)
	}
	return nil
}

// InstallDependencyTrack installs PostgreSQL and DependencyTrack into the
// cluster. DependencyTrack requires a PostgreSQL service named "postgresql"
// and a Key Encryption Key (KEK). Both are provisioned here.
func InstallDependencyTrack() error {
	_, _ = fmt.Fprintf(GinkgoWriter, "Installing DependencyTrack via Helm...\\n")

	// Create the namespace and PostgreSQL deployment (lightweight, no Helm needed).
	_, _ = fmt.Fprintf(GinkgoWriter, "Installing PostgreSQL...\\n")
	tmpFile, err := os.CreateTemp("", "dt-postgres-*.yaml")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()
	_, err = tmpFile.WriteString(postgresqlManifest)
	if err != nil {
		return fmt.Errorf("write postgresql manifest: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close postgresql manifest: %w", err)
	}

	_, err = Run(exec.Command("kubectl", "apply", "-f", tmpFile.Name()))
	if err != nil {
		return fmt.Errorf("apply postgresql manifest: %w", err)
	}

	// Wait for PostgreSQL to be Ready.
	_, _ = fmt.Fprintf(GinkgoWriter, "Waiting for PostgreSQL to be Ready...\\n")
	pgWait := exec.Command("kubectl", "wait", "--for=condition=Ready",
		"pod",
		"-l", "app=postgresql",
		"-n", "dependency-track",
		"--timeout", "3m")
	_, err = Run(pgWait)
	if err != nil {
		return fmt.Errorf("timeout waiting for PostgreSQL pod: %w", err)
	}

	// Add the DependencyTrack Helm chart repo.
	_, _ = fmt.Fprintf(GinkgoWriter, "Installing DependencyTrack...\\n")
	addCmd := exec.Command("helm", "repo", "add", "dependencytrack",
		"https://dependencytrack.github.io/helm-charts")
	_, _ = Run(addCmd) // idempotent

	updateCmd := exec.Command("helm", "repo", "update")
	_, err = Run(updateCmd)
	if err != nil {
		return fmt.Errorf("helm repo update: %w", err)
	}

	// Generate a KEK (Key Encryption Key) required by DependencyTrack.
	keq, err := GenerateKEK()
	if err != nil {
		return fmt.Errorf("generate KEK: %w", err)
	}

	// Pass database credentials and the generated KEK through a mode-0600
	// values file so they never appear in process arguments or test logs.
	secretValues, err := json.Marshal(map[string]any{
		"database": map[string]any{
			"username": "dtrack",
			"password": "dtrack123",
		},
		"secretManagement": map[string]any{
			"database": map[string]any{
				"kek": map[string]any{"value": strings.TrimSpace(keq)},
			},
		},
		"apiServer": map[string]any{
			"extraEnv": []map[string]any{
				{"name": "DT_OIDC_ENABLED", "value": "true"},
				{"name": "DT_OIDC_ISSUER", "value": "http://dex.dependency-track.svc.cluster.local:5556/dex"},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("encode DependencyTrack Helm values: %w", err)
	}
	valuesFile, err := os.CreateTemp("", "dependency-track-values-*.json")
	if err != nil {
		return fmt.Errorf("create DependencyTrack Helm values: %w", err)
	}
	defer func() { _ = os.Remove(valuesFile.Name()) }()
	if _, err := valuesFile.Write(secretValues); err != nil {
		_ = valuesFile.Close()
		return fmt.Errorf("write DependencyTrack Helm values: %w", err)
	}
	if err := valuesFile.Close(); err != nil {
		return fmt.Errorf("close DependencyTrack Helm values: %w", err)
	}

	// Provision Dex OIDC provider so the issuer is resolvable before DT
	// boots and fetches its well-known OIDC metadata.
	if err := ProvisionDex(); err != nil {
		return fmt.Errorf("provision dex: %w", err)
	}

	installCmd := exec.Command("helm", "upgrade", "--install", "my-dependency-track",
		"dependencytrack/dependency-track",
		"--version", "2.0.0-rc.2",
		"--namespace", "dependency-track",
		"--values", valuesFile.Name(),
	)
	_, err = Run(installCmd)
	if err != nil {
		return fmt.Errorf("helm install dependency-track: %w", err)
	}

	// Wait for the API server pod to be Running.
	_, _ = fmt.Fprintf(GinkgoWriter, "Waiting for DependencyTrack API pod to be Ready...\\n")
	waitCmd := exec.Command("kubectl", "wait", "--for=condition=Ready",
		"pod",
		"-l", "app.kubernetes.io/instance=my-dependency-track,app.kubernetes.io/component=api-server",
		"-n", "dependency-track",
		"--timeout", "5m")
	_, err = Run(waitCmd)
	if err != nil {
		return fmt.Errorf("timeout waiting for DependencyTrack API pod: %w", err)
	}

	// Wait for the API to answer requests. Run the probe in the API pod because
	// the service DNS name is only resolvable inside the cluster. The per-request
	// timeout prevents a failed probe from hanging the suite indefinitely.
	_, _ = fmt.Fprintf(GinkgoWriter, "Waiting for DependencyTrack API to accept connections...\\n")
	var lastErr error
	for i := 0; i < 120; i++ {
		checkCmd := exec.Command("kubectl", "exec",
			"deployment/my-dependency-track-api-server",
			"--namespace", "dependency-track",
			"--", "curl", "--fail", "--silent", "--show-error", "--max-time", "5",
			"--output", "/dev/null", "http://127.0.0.1:8080/api/version")
		if _, err := Run(checkCmd); err == nil {
			_, _ = fmt.Fprintf(GinkgoWriter, "DependencyTrack API is ready.\\n")
			break
		} else {
			lastErr = err
		}
		time.Sleep(1 * time.Second)
	}
	if lastErr != nil {
		return fmt.Errorf("DependencyTrack API not ready after 2 minutes: %w", lastErr)
	}

	return nil
}

// IsDependencyTrackReady checks whether the DependencyTrack API pod is ready
// by attempting to reach the readiness endpoint from a kubectl exec context.
func IsDependencyTrackReady() bool {
	cmd := exec.Command("kubectl", "exec",
		"deployment/my-dependency-track-api-server",
		"--namespace", "dependency-track",
		"--", "curl", "--fail", "--silent", "--max-time", "3",
		"--output", "/dev/null", "http://127.0.0.1:8080/api/version")
	return cmd.Run() == nil
}

// VerifyOIDCAvailable checks whether the DependencyTrack API reports OIDC
// as available by exec-ing into the API pod and curling GET /api/v1/oidc/available.
// Prints the result to GinkgoWriter; returns false (halting the suite in
// BeforeSuite) when OIDC is not yet wired up, preventing misleading flake.
func VerifyOIDCAvailable() bool {
	cmd := exec.Command("kubectl", "exec",
		"deployment/my-dependency-track-api-server",
		"--namespace", "dependency-track",
		"--", "curl", "--silent", "--max-time", "5",
		"http://127.0.0.1:8080/api/v1/oidc/available")
	out, err := Run(cmd)
	if err != nil {
		_, _ = fmt.Fprintf(GinkgoWriter, "VerifyOIDCAvailable: probe failed: %v\n", err)
		return false
	}
	trimmed := strings.TrimSpace(out)
	available := trimmed == "true"
	_, _ = fmt.Fprintf(GinkgoWriter, "VerifyOIDCAvailable: oidc.available=%s (available=%v)\n", trimmed, available)
	return available
}

// InstallOrUpgradeDependencyTrack installs PostgreSQL and DependencyTrack if they
// are not already present.  When DT is already running (detected via API probe)
// this is a no-op so that preserved clusters skip the expensive install step.
func InstallOrUpgradeDependencyTrack() error {
	// Quick probe: if the API is already answering, skip installation.
	if IsDependencyTrackReady() {
		_, _ = fmt.Fprintf(GinkgoWriter, "DependencyTrack is already running – skipping install.\\n")
		return nil
	}
	return InstallDependencyTrack()
}

// UninstallDependencyTrack removes the DependencyTrack Helm release and namespace.
func UninstallDependencyTrack() {
	_, _ = fmt.Fprintf(GinkgoWriter, "Uninstalling DependencyTrack...\\n")

	// Clean up Dex OIDC provider resources.
	_, _ = Run(exec.Command("kubectl", "delete", "deployment", "dex",
		"-n", "dependency-track", "--ignore-not-found=true"))
	_, _ = Run(exec.Command("kubectl", "delete", "service", "dex",
		"-n", "dependency-track", "--ignore-not-found=true"))
	_, _ = Run(exec.Command("kubectl", "delete", "configmap", "dex-config",
		"-n", "dependency-track", "--ignore-not-found=true"))

	// Uninstall Helm release.
	uninstallCmd := exec.Command("helm", "uninstall", "my-dependency-track", "--namespace", "dependency-track")
	_, _ = Run(uninstallCmd)

	// Delete the namespace (also removes PostgreSQL).
	_, _ = fmt.Fprintf(GinkgoWriter, "Waiting for pods to terminate...\\n")
	_, _ = Run(exec.Command("kubectl", "delete", "namespace", "dependency-track", "--ignore-not-found=true"))

	// Remove the Helm repo.
	removeCmd := exec.Command("helm", "repo", "remove", "dependencytrack")
	_, _ = Run(removeCmd)

	_, _ = fmt.Fprintf(GinkgoWriter, "DependencyTrack uninstalled.\\n")
}

// minPasswordLen mirrors the constant in internal/auth/rotator.go so we
// can detect when PasswordRotationRunnable has finished rotating.
const minPasswordLen = 30

// WaitForPasswordRotation polls the operator namespace for the credentials
// secret and blocks until the password has been rotated to at least
// minPasswordLen characters (the minimum enforced by
// PasswordRotationRunnable).
func WaitForPasswordRotation(secretName, namespace string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		cmd := exec.Command("kubectl", "get", "secret", secretName, "-n", namespace,
			"-o", "jsonpath={.data.password}")
		out, err := Run(cmd)
		if err != nil {
			// Secret may not exist yet (PasswordRotationRunnable creates it).
			time.Sleep(2 * time.Second)
			continue
		}
		// Decode base64 password and check length.
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(out))
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if len(decoded) >= minPasswordLen {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("password rotation timed out after %v", timeout)
}

// DeployOperatorHelm installs the operator Helm chart with the given project
// directory, image reference, DependencyTrack endpoint, and target namespace.
// It overrides the default image (controller:latest) with the provided image
// and sets imagePullPolicy to IfNotPresent so the locally-loaded Kind image
// is used instead of pulling from a registry.
func DeployOperatorHelm(projectDir, image, deptrackURL, namespace string) error {
	repository, tag, err := splitImage(image)
	if err != nil {
		return err
	}
	chartPath := filepath.Join(projectDir, "deploy", "charts", "dependencytrack-operator")
	cmd := exec.Command("helm", "install", "deptrack-operator", chartPath,
		"--namespace", namespace,
		"--set", "fullnameOverride=deptrack-operator",
		"--set", "controllerManager.manager.image.repository="+repository,
		"--set", "controllerManager.manager.image.tag="+tag,
		"--set", "controllerManager.manager.image.pullPolicy=IfNotPresent",
		"--set", "controllerManager.manager.env.deptrackUrl="+deptrackURL,
		"--set", "controllerManager.manager.env.deptrackCredentialsSecret=deptrack-credentials",
		"--wait", "--timeout", "3m")
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("helm install operator: %w", err)
	}
	return nil
}

// DeleteOperatorHelm removes the operator Helm release from the given namespace.
func DeleteOperatorHelm(namespace string) {
	_, _ = Run(exec.Command("helm", "uninstall", "deptrack-operator", "--namespace", namespace))
}

// OperatorPodReady checks whether the controller-manager deployment has at
// least one ready replica in the given namespace.
func OperatorPodReady(namespace string) bool {
	cmd := exec.Command("kubectl", "get", "deployment",
		"deptrack-operator-controller-manager",
		"-n", namespace, "-o", "jsonpath={.status.readyReplicas}")
	out, err := Run(cmd)
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) == "1"
}

// InstallOperatorHelm installs the packaged operator chart with the locally
// built image and live Dependency-Track endpoint used by the e2e suite.
// Uses helm upgrade --install so it is idempotent when called multiple times
// (e.g. when the notification lifecycle tests deploy the operator first and
// the Manager BeforeAll runs later).
func InstallOperatorHelm(projectDir, image, deptrackURL, namespace string) error {
	repository, tag, err := splitImage(image)
	if err != nil {
		return err
	}
	chartPath := filepath.Join(projectDir, "deploy", "charts", "dependencytrack-operator")
	cmd := exec.Command("helm", "upgrade", "--install", "deptrack-operator", chartPath,
		"--namespace", namespace,
		"--set", "fullnameOverride=deptrack-operator",
		"--set", "controllerManager.manager.image.repository="+repository,
		"--set", "controllerManager.manager.image.tag="+tag,
		"--set", "controllerManager.manager.image.pullPolicy=IfNotPresent",
		"--set", "controllerManager.manager.env.deptrackUrl="+deptrackURL,
		"--set", "controllerManager.manager.env.deptrackCredentialsSecret=deptrack-credentials",
		"--wait", "--timeout", "3m")
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("helm upgrade --install operator: %w", err)
	}
	return nil
}

// IsOperatorDeployed checks whether the operator deployment exists in the given namespace.
func IsOperatorDeployed(namespace string) bool {
	cmd := exec.Command("kubectl", "get", "deployment",
		"deptrack-operator-controller-manager", "-n", namespace)
	return cmd.Run() == nil
}

// UpgradeOperatorHelm exercises the packaged chart upgrade path while retaining
// the values used for the initial installation.
func UpgradeOperatorHelm(projectDir, image, deptrackURL, namespace string) error {
	repository, tag, err := splitImage(image)
	if err != nil {
		return err
	}
	chartPath := filepath.Join(projectDir, "deploy", "charts", "dependencytrack-operator")
	cmd := exec.Command("helm", "upgrade", "deptrack-operator", chartPath,
		"--namespace", namespace,
		"--reuse-values",
		"--set", "controllerManager.manager.image.repository="+repository,
		"--set", "controllerManager.manager.image.tag="+tag,
		"--set", "controllerManager.manager.image.pullPolicy=IfNotPresent",
		"--set", "controllerManager.manager.env.deptrackUrl="+deptrackURL,
		"--set", "controllerManager.manager.env.deptrackCredentialsSecret=deptrack-credentials",
		"--wait", "--timeout", "3m")
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("helm upgrade operator: %w", err)
	}
	return nil
}

// UninstallOperatorHelm removes the operator Helm release.
func UninstallOperatorHelm(namespace string) {
	_, _ = Run(exec.Command("helm", "uninstall", "deptrack-operator", "--namespace", namespace))
}

func splitImage(image string) (string, string, error) {
	separator := strings.LastIndex(image, ":")
	if separator <= strings.LastIndex(image, "/") || separator == len(image)-1 {
		return "", "", fmt.Errorf("image %q must include a tag", image)
	}
	return image[:separator], image[separator+1:], nil
}

// DependencyTrackHost returns the Kubernetes service hostname for the
// DependencyTrack instance, suitable for use as DEPTRACK_URL
// (e.g. http://my-dependency-track-api-server.dependency-track.svc.cluster.local:8080).
func DependencyTrackHost() string {
	return "http://my-dependency-track-api-server.dependency-track.svc.cluster.local:8080"
}

// DeleteKindCluster tears down the Kind cluster used by the e2e suite.
// The cluster name comes from the KIND_CLUSTER environment variable
// (defaulting to "kind").
func DeleteKindCluster() {
	cluster := "kind"
	if v, ok := os.LookupEnv("KIND_CLUSTER"); ok {
		cluster = v
	}
	_, _ = fmt.Fprintf(GinkgoWriter, "Deleting Kind cluster %q...\n", cluster)
	cmd := exec.Command("kind", "delete", "cluster", "--name", cluster)
	if _, err := Run(cmd); err != nil {
		warnError(fmt.Errorf("kind delete cluster: %w", err))
	}
	_, _ = fmt.Fprintf(GinkgoWriter, "Kind cluster %q deleted.\n", cluster)
}

const certManagerVersion = "v1.14.5"
const certManagerManifestURL = "https://github.com/cert-manager/cert-manager/releases/download/%s/cert-manager.yaml"

// ProvisionCertManager installs cert-manager into the cluster so the operator
// chart's Certificate/Issuer CRs reconcile before the webhook serves. Skipped
// when CERT_MANAGER_INSTALL_SKIP=true (mirrors the e2e Makefile contract).
func ProvisionCertManager() error {
	if os.Getenv("CERT_MANAGER_INSTALL_SKIP") == "true" {
		_, _ = fmt.Fprintf(GinkgoWriter, "Skipping cert-manager install (CERT_MANAGER_INSTALL_SKIP=true).\n")
		return nil
	}
	_, _ = fmt.Fprintf(GinkgoWriter, "Installing cert-manager %s...\n", certManagerVersion)
	url := fmt.Sprintf(certManagerManifestURL, certManagerVersion)
	if _, err := Run(exec.Command("kubectl", "apply", "-f", url)); err != nil {
		return fmt.Errorf("apply cert-manager manifest: %w", err)
	}
	if err := waitForCRD("certificates.cert-manager.io", "2m"); err != nil {
		return fmt.Errorf("cert-manager CRD not ready: %w", err)
	}
	for _, dep := range []string{"cert-manager", "cert-manager-cainjector", "cert-manager-webhook"} {
		wait := exec.Command("kubectl", "rollout", "status", "deployment/"+dep,
			"-n", "cert-manager", "--timeout", "2m")
		if _, err := Run(wait); err != nil {
			return fmt.Errorf("cert-manager deployment %s not ready: %w", dep, err)
		}
	}
	return nil
}

// UninstallCertManager removes cert-manager from the cluster (idempotent).
func UninstallCertManager() {
	url := fmt.Sprintf(certManagerManifestURL, certManagerVersion)
	_, _ = fmt.Fprintf(GinkgoWriter, "Uninstalling cert-manager...\n")
	_, _ = Run(exec.Command("kubectl", "delete", "-f", url))
}

// waitForCRD polls until the named CRD reaches the Established condition or times out.
func waitForCRD(name, timeout string) error {
	cmd := exec.Command("kubectl", "wait", "--for=condition=Established",
		"crd", name, "--timeout", timeout)
	if _, err := Run(cmd); err != nil {
		return fmt.Errorf("waiting for CRD %s: %w", name, err)
	}
	return nil
}

// UncommentCode searches for target in the file and remove the comment prefix
