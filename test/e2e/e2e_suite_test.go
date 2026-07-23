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
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/martin31821/dependencytrack-operator/test/utils"
)

var (
	// projectImage is the name of the image which will be build and loaded
	// with the code source changes to be tested.
	projectImage = "example.com/deptrack-operator:v0.0.1"
)

// TestE2E runs the end-to-end (e2e) tests for the project. These tests execute in an isolated,
// temporary environment to validate project changes with the purpose of being used in CI jobs.
// The setup requires KIND_CLUSTER to opt in because it builds and loads a manager image,
// installs cluster dependencies, and tears the selected Kind cluster down. The Makefile's
// test-e2e target supplies KIND_CLUSTER; ordinary `go test ./...` runs skip this package.
func TestE2E(t *testing.T) {
	if os.Getenv("KIND_CLUSTER") == "" {
		t.Skip("set KIND_CLUSTER or run `make test-e2e` to execute the e2e suite")
	}

	RegisterFailHandler(Fail)
	_, _ = fmt.Fprintf(GinkgoWriter, "Starting deptrack-operator integration test suite\n")
	RunSpecs(t, "e2e suite")
}

var _ = BeforeSuite(func() {
	By("building the manager(Operator) image")
	cmd := exec.Command("make", "docker-build", fmt.Sprintf("IMG=%s", projectImage))
	_, err := utils.Run(cmd)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to build the manager(Operator) image")

	// TODO(user): If you want to change the e2e test vendor from Kind, ensure the image is
	// built and available before running the tests. Also, remove the following block.
	By("loading the manager(Operator) image on Kind")
	err = utils.LoadImageToKindClusterWithName(projectImage)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to load the manager(Operator) image into Kind")

	// Deploy a real DependencyTrack instance via Helm so the operator can
	// reconcile CRs against a live backend.
	By("installing DependencyTrack via Helm")
	Expect(utils.InstallOrUpgradeDependencyTrack()).To(Succeed(), "Failed to install DependencyTrack")
})

var _ = AfterSuite(func() {
	skipCluster := os.Getenv("E2E_SKIP_CLUSTER_TEARDOWN") == "true"

	// Only tear down DependencyTrack when keeping the cluster – if the
	// cluster is going away, the namespace is deleted anyway.
	if skipCluster {
		_, _ = fmt.Fprintf(GinkgoWriter, "Skipping DependencyTrack teardown (cluster teardown pending)\n")
	} else if os.Getenv("E2E_SKIP_DT_TEARDOWN") != "true" {
		utils.UninstallDependencyTrack()
	} else {
		_, _ = fmt.Fprintf(GinkgoWriter, "Skipping DependencyTrack teardown (E2E_SKIP_DT_TEARDOWN=true)\n")
	}

	// Tear down the Kind cluster when not preserving the cluster between runs.
	if !skipCluster {
		utils.DeleteKindCluster()
	} else {
		_, _ = fmt.Fprintf(GinkgoWriter, "Skipping Kind cluster teardown (E2E_SKIP_CLUSTER_TEARDOWN=true)\n")
	}
})
