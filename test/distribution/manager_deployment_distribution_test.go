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

package distribution

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestManagerWebhookCertMounted guards MEM023: the manager Deployment in the Helm
// chart MUST mount the Team validating webhook serving cert at
// /tmp/k8s-webhook-server/serving-certs and reference the same Secret that the
// webhook-server-cert template publishes. It also pins the kustomize source patch
// that drives the mount, so a dropped mount, renamed secret, or unwired patch
// regresses the build.
func TestManagerWebhookCertMounted(t *testing.T) {
	root := findProjectRoot(t)

	const wantMountPath = "/tmp/k8s-webhook-server/serving-certs"

	depPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "deployment.yaml")
	certPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "webhook-server-cert.yaml")
	patchPath := filepath.Join(root, "config", "default", "webhook_server_cert_mount_patch.yaml")
	kustPath := filepath.Join(root, "config", "default", "kustomization.yaml")

	for _, p := range []string{depPath, certPath, patchPath, kustPath} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Fatalf("missing expected file: %s", p)
		}
	}

	// Drive the mount through kustomize (source of truth for the chart).
	ks, err := os.ReadFile(kustPath)
	if err != nil {
		t.Fatalf("read kustomization: %v", err)
	}
	if !bytes.Contains(ks, []byte("webhook_server_cert_mount_patch.yaml")) {
		t.Error("config/default/kustomization.yaml does not wire webhook_server_cert_mount_patch.yaml")
	}

	// Parse the Helm deployment template. The image line straddles two template
	// fragments separated by a literal ':' that survives stripping as a dangling
	// "image: :" and breaks YAML parsing, so it is normalized to a plain scalar.
	var dep map[string]interface{}
	{
		raw, rerr := os.ReadFile(depPath)
		if rerr != nil {
			t.Fatalf("read %s: %v", depPath, rerr)
		}
		sanitized := stripHelmTemplates(string(raw))
		imageColon := regexp.MustCompile(`(?m)^(\s*image:\s+).*$`)
		if err := yaml.Unmarshal([]byte(imageColon.ReplaceAllString(sanitized, "${1}controller:latest")), &dep); err != nil {
			t.Fatalf("failed to parse Helm deployment template: %v", err)
		}
	}

	template := getNestedMap(dep, "spec", "template", "spec")
	if template == nil {
		t.Fatal("deployment spec.template.spec missing")
	}

	containers := getNestedMapSlice(template, "containers")
	if len(containers) == 0 {
		t.Fatal("no containers in deployment spec")
	}
	manager := containers[0]
	if got := getString(manager, "name"); got != "manager" {
		t.Fatalf("first container name = %q, want manager", got)
	}

	// Positive: the manager container mounts the webhook serving cert read-only.
	vmounts := getNestedMapSlice(manager, "volumeMounts")
	var foundVM *map[string]interface{}
	for i := range vmounts {
		if getString(vmounts[i], "name") == "webhook-server-cert" {
			foundVM = &vmounts[i]
			break
		}
	}
	if foundVM == nil {
		t.Fatal("manager container missing volumeMount named webhook-server-cert")
	}
	if got := getString(*foundVM, "mountPath"); got != wantMountPath {
		t.Errorf("webhook volumeMount.mountPath = %q, want %q", got, wantMountPath)
	}
	if ro := (*foundVM)["readOnly"]; ro != true {
		t.Errorf("webhook volumeMount.readOnly = %v, want true", ro)
	}

	// Positive: the pod spec declares the matching secret-backed volume.
	volumes := getNestedMapSlice(template, "volumes")
	var foundVol *map[string]interface{}
	for i := range volumes {
		if getString(volumes[i], "name") == "webhook-server-cert" {
			foundVol = &volumes[i]
			break
		}
	}
	if foundVol == nil {
		t.Fatal("pod spec missing volume named webhook-server-cert")
	}
	sec := getNestedMap(*foundVol, "secret")
	if sec == nil {
		t.Fatal("webhook volume missing secret block")
	}
	sn := getString(sec, "secretName")
	if !strings.HasSuffix(sn, "-webhook-server-cert") {
		t.Errorf("webhook volume.secretName = %q, want to end with -webhook-server-cert", sn)
	}

	// Correlation: the published Secret shares the cert suffix the volume references.
	var cert map[string]interface{}
	if err := parseHelmYAML(certPath, &cert); err != nil {
		t.Fatalf("failed to parse webhook-server-cert template: %v", err)
	}
	certName := strings.TrimSpace(getString(getNestedMap(cert, "metadata"), "name"))
	if !strings.HasSuffix(certName, "-webhook-server-cert") {
		t.Errorf("webhook-server-cert Secret name = %q, want to end with -webhook-server-cert", certName)
	}
}
