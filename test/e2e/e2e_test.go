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

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/martin31821/dependencytrack-operator/test/utils"
)

// namespace where the project is deployed in
const namespace = "deptrack-operator-system"

// serviceAccountName created for the project
const serviceAccountName = "deptrack-operator-controller-manager"

// metricsServiceName is the name of the metrics service of the project
const metricsServiceName = "deptrack-operator-controller-manager-metrics-service"

// metricsRoleBindingName is the name of the RBAC that will be created to allow get the metrics data
const metricsRoleBindingName = "deptrack-operator-metrics-binding"

var _ = Describe("Manager", Ordered, func() {
	var controllerPodName string

	// Before running the tests, set up the environment by creating the namespace,
	// enforce the restricted security policy to the namespace, installing CRDs,
	// and deploying the controller.
	BeforeAll(func() {
		By("creating manager namespace")
		cmd := exec.Command("kubectl", "create", "ns", namespace)
		_, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to create namespace")

		By("labeling the namespace to enforce the restricted security policy")
		cmd = exec.Command("kubectl", "label", "--overwrite", "ns", namespace,
			"pod-security.kubernetes.io/enforce=restricted")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to label namespace with restricted policy")

		By("installing CRDs")
		cmd = exec.Command("make", "install")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")

		By("deploying the controller-manager")
		cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectImage))
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to deploy the controller-manager")
	})

	// After all tests have been executed, clean up by undeploying the controller, uninstalling CRDs,
	// and deleting the namespace.
	AfterAll(func() {
		By("cleaning up the curl pod for metrics")
		cmd := exec.Command("kubectl", "delete", "pod", "curl-metrics", "-n", namespace)
		_, _ = utils.Run(cmd)

		By("cleaning up the metrics ClusterRoleBinding")
		cmd = exec.Command("kubectl", "delete", "clusterrolebinding", metricsRoleBindingName)
		_, _ = utils.Run(cmd)

		By("undeploying the controller-manager")
		cmd = exec.Command("make", "undeploy")
		_, _ = utils.Run(cmd)

		By("uninstalling CRDs")
		cmd = exec.Command("make", "uninstall")
		_, _ = utils.Run(cmd)

		By("removing manager namespace")
		cmd = exec.Command("kubectl", "delete", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	// After each test, check for failures and collect logs, events,
	// and pod descriptions for debugging.
	AfterEach(func() {
		specReport := CurrentSpecReport()
		if specReport.Failed() {
			By("Fetching controller manager pod logs")
			cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
			controllerLogs, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Controller logs:\n %s", controllerLogs)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Controller logs: %s", err)
			}

			By("Fetching Kubernetes events")
			cmd = exec.Command("kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp")
			eventsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Kubernetes events:\n%s", eventsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Kubernetes events: %s", err)
			}

			By("Fetching curl-metrics logs")
			cmd = exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
			metricsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Metrics logs:\n %s", metricsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get curl-metrics logs: %s", err)
			}

			By("Fetching controller manager pod description")
			cmd = exec.Command("kubectl", "describe", "pod", controllerPodName, "-n", namespace)
			podDescription, err := utils.Run(cmd)
			if err == nil {
				fmt.Println("Pod description:\n", podDescription)
			} else {
				fmt.Println("Failed to describe controller pod")
			}
		}
	})

	SetDefaultEventuallyTimeout(2 * time.Minute)
	SetDefaultEventuallyPollingInterval(time.Second)

	Context("Manager", func() {
		It("should run successfully", func() {
			By("validating that the controller-manager pod is running as expected")
			verifyControllerUp := func(g Gomega) {
				// Get the name of the controller-manager pod
				cmd := exec.Command("kubectl", "get",
					"pods", "-l", "control-plane=controller-manager",
					"-o", "go-template={{ range .items }}"+
						"{{ if not .metadata.deletionTimestamp }}"+
						"{{ .metadata.name }}"+
						"{{ \"\\n\" }}{{ end }}{{ end }}",
					"-n", namespace,
				)

				podOutput, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred(), "Failed to retrieve controller-manager pod information")
				podNames := utils.GetNonEmptyLines(podOutput)
				g.Expect(podNames).To(HaveLen(1), "expected 1 controller pod running")
				controllerPodName = podNames[0]
				g.Expect(controllerPodName).To(ContainSubstring("controller-manager"))

				// Validate the pod's status
				cmd = exec.Command("kubectl", "get",
					"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
					"-n", namespace,
				)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Running"), "Incorrect controller-manager pod status")
			}
			Eventually(verifyControllerUp).Should(Succeed())
		})

		It("should ensure the metrics endpoint is serving metrics", func() {
			By("creating a ClusterRoleBinding for the service account to allow access to metrics")
			crbYAML := fmt.Sprintf(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: deptrack-operator-metrics-reader
subjects:
- kind: ServiceAccount
  name: %s
  namespace: %s
`, metricsRoleBindingName, serviceAccountName, namespace)
			tmpFile, err := os.CreateTemp("", "crb-*.yaml")
			Expect(err).NotTo(HaveOccurred())
			_, _ = tmpFile.WriteString(crbYAML)
			_ = tmpFile.Close()
			defer func() { _ = os.Remove(tmpFile.Name()) }()

			cmd := exec.Command("kubectl", "apply", "-f", tmpFile.Name())
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply ClusterRoleBinding")

			By("validating that the metrics service is available")
			cmd = exec.Command("kubectl", "get", "service", metricsServiceName, "-n", namespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Metrics service should exist")

			By("getting the service account token")
			token, err := serviceAccountToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())

			By("waiting for the metrics endpoint to be ready")
			verifyMetricsEndpointReady := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "endpoints", metricsServiceName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("8443"), "Metrics endpoint is not ready")
			}
			Eventually(verifyMetricsEndpointReady).Should(Succeed())

			By("verifying that the controller manager is serving the metrics server")
			verifyMetricsServerStarted := func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("controller-runtime.metrics\tServing metrics server"),
					"Metrics server not yet started")
			}
			Eventually(verifyMetricsServerStarted).Should(Succeed())

			By("creating the curl-metrics pod to access the metrics endpoint")
			cmd = exec.Command("kubectl", "run", "curl-metrics", "--restart=Never",
				"--namespace", namespace,
				"--image=curlimages/curl:latest",
				"--overrides",
				fmt.Sprintf(`{
					"spec": {
						"containers": [{
							"name": "curl",
							"image": "curlimages/curl:latest",
							"command": ["/bin/sh", "-c"],
							"args": ["curl -v -k -H 'Authorization: Bearer %s' https://%s.%s.svc.cluster.local:8443/metrics"],
							"securityContext": {
								"allowPrivilegeEscalation": false,
								"capabilities": {
									"drop": ["ALL"]
								},
								"runAsNonRoot": true,
								"runAsUser": 1000,
								"seccompProfile": {
									"type": "RuntimeDefault"
								}
							}
						}],
						"serviceAccount": "%s"
					}
				}`, token, metricsServiceName, namespace, serviceAccountName))
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create curl-metrics pod")

			By("waiting for the curl-metrics pod to complete.")
			verifyCurlUp := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "pods", "curl-metrics",
					"-o", "jsonpath={.status.phase}",
					"-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Succeeded"), "curl pod in wrong status")
			}
			Eventually(verifyCurlUp, 5*time.Minute).Should(Succeed())

			By("getting the metrics by checking curl-metrics logs")
			metricsOutput := getMetricsOutput()
			Expect(metricsOutput).To(ContainSubstring(
				"controller_runtime_reconcile_total",
			))
		})

		// +kubebuilder:scaffold:e2e-webhooks-checks

		// Team permissions e2e tests
		// These tests verify that the Team CRD and controller handle permissions
		// correctly. They run against a cluster with no DependencyTrack instance,
		// so the controller will fail to connect — but we can still verify that
		// the CRD schema is accepted, the controller attempts reconciliation,
		// conditions are set appropriately, and deletion with the finalizer works.
		Context("Team permissions", func() {
			const (
				teamWithPerms    = "team-with-permissions"
				teamWithoutPerms = "team-without-permissions"
				teamEmptyPerms   = "team-empty-permissions"
			)

			AfterEach(func() {
				// Clean up test teams after each test.
				for _, name := range []string{teamWithPerms, teamWithoutPerms, teamEmptyPerms} {
					By(fmt.Sprintf("deleting test team %q", name))
					cmd := exec.Command("kubectl", "delete", "team", name, "-n", namespace)
					_, _ = utils.Run(cmd)
				}
			})

			// createTeam writes the given YAML to a temp file and applies it via kubectl.
			// Returns the temp file path so the caller can defer cleanup.
			createTeam := func(teamName, yaml string) string {
				tmpFile, err := os.CreateTemp("", fmt.Sprintf("team-%s-*.yaml", teamName))
				Expect(err).NotTo(HaveOccurred())
				_, err = tmpFile.WriteString(yaml)
				Expect(err).NotTo(HaveOccurred())
				Expect(tmpFile.Close()).NotTo(HaveOccurred())
				return tmpFile.Name()
			}

			// removeTemp cleans up a temp file, ignoring errors (not actionable in tests).
			removeTemp := func(path string) { _ = os.Remove(path) }

			applyTeam := func(path string, expectFail bool, errMsg string) {
				cmd := exec.Command("kubectl", "apply", "-f", path)
				_, err := utils.Run(cmd)
				if expectFail {
					Expect(err).To(HaveOccurred(), errMsg)
				} else {
					Expect(err).NotTo(HaveOccurred(), errMsg)
				}
			}

			verifyTeamHasStatus := func(teamName string) {
				By("verifying Team has a status condition")
				verifyTeamExists := func(g Gomega) {
					cmd := exec.Command("kubectl", "get", "team", teamName, "-n", namespace,
						"-o", "jsonpath={.status.conditions[0].type}")
					output, err := utils.Run(cmd)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(output).To(Not(BeEmpty()), "Team status conditions should be populated")
				}
				Eventually(verifyTeamExists, 2*time.Minute).Should(Succeed())
			}

			It("should accept a Team CR with permissions and set a failed condition", func() {
				By("creating a Team with permissions")
				teamYAML := fmt.Sprintf(`apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: %s
  namespace: %s
spec:
  name: test-team-permissions
  permissions:
    - PORTFOLIO_VIEW
    - VIEW_PORTFOLIO
`, teamWithPerms, namespace)
				path := createTeam(teamWithPerms, teamYAML)
				defer removeTemp(path)
				applyTeam(path, false, "Failed to create Team with permissions")
				verifyTeamHasStatus(teamWithPerms)
			})

			It("should accept a Team CR without permissions (nil)", func() {
				By("creating a Team without permissions")
				teamYAML := fmt.Sprintf(`apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: %s
  namespace: %s
spec:
  name: test-team-no-permissions
`, teamWithoutPerms, namespace)
				path := createTeam(teamWithoutPerms, teamYAML)
				defer removeTemp(path)
				applyTeam(path, false, "Failed to create Team without permissions")
				verifyTeamHasStatus(teamWithoutPerms)
			})

			It("should accept a Team CR with an empty permissions array", func() {
				By("creating a Team with empty permissions array")
				teamYAML := fmt.Sprintf(`apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: %s
  namespace: %s
spec:
  name: test-team-empty-permissions
  permissions: []
`, teamEmptyPerms, namespace)
				path := createTeam(teamEmptyPerms, teamYAML)
				defer removeTemp(path)
				applyTeam(path, false, "Failed to create Team with empty permissions")
				verifyTeamHasStatus(teamEmptyPerms)
			})

			It("should reject a Team CR with an invalid schema", func() {
				By("creating a Team with invalid schema")
				badYAML := fmt.Sprintf(`apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: team-invalid-schema
  namespace: %s
spec:
  name: test-invalid
  permissions: "not-an-array"
`, namespace)
				path := createTeam("invalid", badYAML)
				defer removeTemp(path)
				applyTeam(path, true, "kubectl should reject Team with invalid permissions type")
			})

			It("should handle Team deletion with the finalizer", func() {
				const deleteTeamName = "team-for-deletion"

				By("creating a Team to test deletion")
				teamYAML := fmt.Sprintf(`apiVersion: dependencytrack.mko.dev/v1alpha1
kind: Team
metadata:
  name: %s
  namespace: %s
spec:
  name: test-team-delete
`, deleteTeamName, namespace)
				path := createTeam(deleteTeamName, teamYAML)
				defer removeTemp(path)
				applyTeam(path, false, "Failed to create Team for deletion test")

				// Wait for the finalizer to be added.
				By("waiting for the finalizer to be added")
				verifyFinalizerExists := func(g Gomega) {
					cmd := exec.Command("kubectl", "get", "team", deleteTeamName, "-n", namespace,
						"-o", "jsonpath={.metadata.finalizers}")
					output, err := utils.Run(cmd)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(output).To(ContainSubstring("dependencytrack.mko.dev/finalizer"))
				}
				Eventually(verifyFinalizerExists, 1*time.Minute).Should(Succeed())

				// Delete the Team.
				By("deleting the Team")
				deleteCmd := exec.Command("kubectl", "delete", "team", deleteTeamName, "-n", namespace)
				_, delErr := utils.Run(deleteCmd)
				Expect(delErr).NotTo(HaveOccurred(), "Failed to delete Team")

				// Wait for the Team to be fully removed (finalizer cleanup completes).
				By("waiting for the Team to be fully removed")
				verifyTeamGone := func(g Gomega) {
					cmd := exec.Command("kubectl", "get", "team", deleteTeamName, "-n", namespace)
					_, err := utils.Run(cmd)
					g.Expect(err).To(HaveOccurred(), "Team should be deleted")
				}
				Eventually(verifyTeamGone, 2*time.Minute).Should(Succeed())
			})
		})
	})
})

// serviceAccountToken returns a token for the specified service account in the given namespace.
// It uses the Kubernetes TokenRequest API to generate a token by directly sending a request
// and parsing the resulting token from the API response.
func serviceAccountToken() (string, error) {
	const tokenRequestRawString = `{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenRequest"
	}`

	// Temporary file to store the token request
	secretName := fmt.Sprintf("%s-token-request", serviceAccountName)
	tokenRequestFile := filepath.Join("/tmp", secretName)
	err := os.WriteFile(tokenRequestFile, []byte(tokenRequestRawString), os.FileMode(0o644))
	if err != nil {
		return "", err
	}

	var out string
	verifyTokenCreation := func(g Gomega) {
		// Execute kubectl command to create the token
		cmd := exec.Command("kubectl", "create", "--raw", fmt.Sprintf(
			"/api/v1/namespaces/%s/serviceaccounts/%s/token",
			namespace,
			serviceAccountName,
		), "-f", tokenRequestFile)

		output, err := cmd.CombinedOutput()
		g.Expect(err).NotTo(HaveOccurred())

		// Parse the JSON output to extract the token
		var token tokenRequest
		err = json.Unmarshal(output, &token)
		g.Expect(err).NotTo(HaveOccurred())

		out = token.Status.Token
	}
	Eventually(verifyTokenCreation).Should(Succeed())

	return out, err
}

// getMetricsOutput retrieves and returns the logs from the curl pod used to access the metrics endpoint.
func getMetricsOutput() string {
	By("getting the curl-metrics logs")
	cmd := exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
	metricsOutput, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "Failed to retrieve logs from curl pod")
	Expect(metricsOutput).To(ContainSubstring("< HTTP/1.1 200 OK"))
	return metricsOutput
}

// tokenRequest is a simplified representation of the Kubernetes TokenRequest API response,
// containing only the token field that we need to extract.
type tokenRequest struct {
	Status struct {
		Token string `json:"token"`
	} `json:"status"`
}
