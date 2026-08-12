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

// TestManagerWebhookCertMounted guards MEM023: the manager Deployment in the
// Helm chart MUST mount the Team validating webhook serving cert at
// /tmp/k8s-webhook-server/serving-certs and reference the same Secret that the
// cert-manager Certificate publishes. It also pins the cert-manager
// Certificate/Issuer pairing and the ValidatingWebhookConfiguration
// inject-ca-from annotation, so a dropped mount, renamed secret, unaligned
// Issuer, or value-driven caBundle regresses the build.
func TestManagerWebhookCertMounted(t *testing.T) {
	root := findProjectRoot(t)

	const wantMountPath = "/tmp/k8s-webhook-server/serving-certs"

	depPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "deployment.yaml")
	certMgrCertPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "certmanager-certificate.yaml")
	certMgrIssuerPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "certmanager-issuer.yaml")
	whcPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "webhook-validating-webhook.yaml")
	patchPath := filepath.Join(root, "config", "default", "webhook_server_cert_mount_patch.yaml")
	kustPath := filepath.Join(root, "config", "default", "kustomization.yaml")

	for _, p := range []string{depPath, certMgrCertPath, certMgrIssuerPath, whcPath, patchPath, kustPath} {
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

	// Correlation: the cert-manager Certificate must publish the SAME Secret the
	// volume references, and pair with a self-signed Issuer.
	var cert map[string]interface{}
	if err := parseHelmYAML(certMgrCertPath, &cert); err != nil {
		t.Fatalf("failed to parse certmanager-certificate template: %v", err)
	}
	certMeta := getNestedMap(cert, "metadata")
	certName := strings.TrimSpace(getString(certMeta, "name"))
	if !strings.HasSuffix(certName, "-webhook-server-cert") {
		t.Errorf("Certificate metadata.name = %q, want to end with -webhook-server-cert", certName)
	}
	certSpec := getNestedMap(cert, "spec")
	if certSpec == nil {
		t.Fatal("Certificate spec missing")
	}
	certSecretName := strings.TrimSpace(getString(certSpec, "secretName"))
	if certSecretName != sn {
		t.Errorf("Certificate secretName = %q, want the volume secretName %q", certSecretName, sn)
	}
	issuerRef := getNestedMap(certSpec, "issuerRef")
	if issuerRef == nil {
		t.Fatal("Certificate spec.issuerRef missing")
	}
	if got := getString(issuerRef, "kind"); got != "Issuer" {
		t.Errorf("Certificate issuerRef.kind = %q, want Issuer", got)
	}
	issuerRefName := strings.TrimSpace(getString(issuerRef, "name"))
	if !strings.HasSuffix(issuerRefName, "-self-signed-issuer") {
		t.Errorf("Certificate issuerRef.name = %q, want to end with -self-signed-issuer", issuerRefName)
	}

	// The Issuer must exist and bear the name the Certificate references.
	var issuer map[string]interface{}
	if err := parseHelmYAML(certMgrIssuerPath, &issuer); err != nil {
		t.Fatalf("failed to parse certmanager-issuer template: %v", err)
	}
	ispec := getNestedMap(issuer, "spec")
	if ispec == nil {
		t.Fatal("Issuer spec missing")
	}
	if getNestedMap(ispec, "selfSigned") == nil {
		t.Error("Issuer spec.selfSigned missing; expected a self-signed Issuer")
	}
	issuerName := strings.TrimSpace(getString(getNestedMap(issuer, "metadata"), "name"))
	if issuerName != issuerRefName {
		t.Errorf("Issuer name = %q, want the Certificate issuerRef.name %q", issuerName, issuerRefName)
	}

	// The ValidatingWebhookConfiguration must delegate caBundle to cainjector
	// (annotated) and must NOT carry a value-driven caBundle/fail block.
	whcRaw, err := os.ReadFile(whcPath)
	if err != nil {
		t.Fatalf("read %s: %v", whcPath, err)
	}
	whcS := stripHelmTemplates(string(whcRaw))
	if !bytes.Contains([]byte(whcS), []byte("cert-manager.io/inject-ca-from")) {
		t.Error("ValidatingWebhookConfiguration is missing the cert-manager.io/inject-ca-from annotation")
	}
	if !bytes.Contains([]byte(whcS), []byte("deptrack-operator-webhook-server-cert")) {
		t.Error("ValidatingWebhookConfiguration inject-ca-from does not reference the cert-manager serving Secret")
	}
	if bytes.Contains([]byte(whcS), []byte("caBundle:")) {
		t.Error("ValidatingWebhookConfiguration still embeds a value-driven caBundle block; expected cainjector injection only")
	}
}
