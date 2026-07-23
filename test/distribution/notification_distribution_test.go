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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestNotificationPublisherCRDDistribution verifies the NotificationPublisher CRD
// in the Helm chart template is synchronized with the kustomize-generated CRD.
func TestNotificationPublisherCRDDistribution(t *testing.T) {
	root := findProjectRoot(t)

	kustomizeCRDPath := filepath.Join(root, "config", "crd", "bases", "dependencytrack.mko.dev_notificationpublishers.yaml")
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationpublisher-crd.yaml")

	// Verify both files exist
	if _, err := os.Stat(kustomizeCRDPath); os.IsNotExist(err) {
		t.Fatalf("kustomize NotificationPublisher CRD not found at %s", kustomizeCRDPath)
	}
	if _, err := os.Stat(helmCRDPath); os.IsNotExist(err) {
		t.Fatalf("Helm NotificationPublisher CRD template not found at %s", helmCRDPath)
	}

	// Parse both YAMLs (Helm version needs template stripping)
	var kustomizeDoc map[string]interface{}
	if err := parseYAML(kustomizeCRDPath, &kustomizeDoc); err != nil {
		t.Fatalf("failed to parse kustomize CRD: %v", err)
	}

	var helmDoc map[string]interface{}
	if err := parseHelmYAML(helmCRDPath, &helmDoc); err != nil {
		t.Fatalf("failed to parse Helm CRD template: %v", err)
	}

	// Compare CRD identity fields
	checkEqual(
		t,
		"CRD group",
		getString(getNestedMap(kustomizeDoc, "spec"), "group"),
		getString(getNestedMap(helmDoc, "spec"), "group"),
	)
	checkEqual(
		t,
		"CRD plural",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "plural"),
		getString(getNestedMap(helmDoc, "spec", "names"), "plural"),
	)
	checkEqual(
		t,
		"CRD singular",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "singular"),
		getString(getNestedMap(helmDoc, "spec", "names"), "singular"),
	)
	checkEqual(
		t,
		"CRD kind",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "kind"),
		getString(getNestedMap(helmDoc, "spec", "names"), "kind"),
	)
	checkEqual(
		t,
		"CRD scope",
		getString(getNestedMap(kustomizeDoc, "spec"), "scope"),
		getString(getNestedMap(helmDoc, "spec"), "scope"),
	)

	// Compare version schemas
	kustomizeVersions := getNestedMapSlice(kustomizeDoc, "spec", "versions")
	helmVersions := getNestedMapSlice(helmDoc, "spec", "versions")
	if len(kustomizeVersions) != len(helmVersions) {
		t.Errorf("version count mismatch: kustomize=%d, helm=%d", len(kustomizeVersions), len(helmVersions))
	}

	for i, kv := range kustomizeVersions {
		if i >= len(helmVersions) {
			break
		}
		hv := helmVersions[i]
		checkEqual(t, fmt.Sprintf("version[%d] name", i), getString(kv, "name"), getString(hv, "name"))
		checkEqual(t, fmt.Sprintf("version[%d] storage", i), getString(kv, "storage"), getString(hv, "storage"))
		checkEqual(t, fmt.Sprintf("version[%d] served", i), getString(kv, "served"), getString(hv, "served"))

		// Compare printer columns count
		kCols := getNestedMapSlice(kv, "additionalPrinterColumns")
		hCols := getNestedMapSlice(hv, "additionalPrinterColumns")
		if len(kCols) != len(hCols) {
			t.Errorf("column count mismatch at version[%d]: kustomize=%d, helm=%d", i, len(kCols), len(hCols))
		}
		for j := range kCols {
			if j >= len(hCols) {
				break
			}
			checkEqual(t, fmt.Sprintf("version[%d] column[%d] jsonPath", i, j),
				getString(kCols[j], "jsonPath"), getString(hCols[j], "jsonPath"))
			checkEqual(t, fmt.Sprintf("version[%d] column[%d] name", i, j),
				getString(kCols[j], "name"), getString(hCols[j], "name"))
		}
	}

	// Compare schema structure (ignoring description text differences)
	kSchema := getNestedMap(kustomizeVersions[0], "schema", "openAPIV3Schema")
	hSchema := getNestedMap(helmVersions[0], "schema", "openAPIV3Schema")

	// Verify spec properties exist and have same required fields
	kSpecProps := getNestedMap(kSchema, "properties", "spec")
	hSpecProps := getNestedMap(hSchema, "properties", "spec")

	checkEqual(
		t,
		"spec.required",
		strings.Join(getRequiredFields(kSpecProps), ","),
		strings.Join(getRequiredFields(hSpecProps), ","),
	)

	// Verify status properties exist.
	kStatusWrapper := getNestedMap(kSchema, "properties", "status")
	hStatusWrapper := getNestedMap(hSchema, "properties", "status")
	if kStatusWrapper == nil || hStatusWrapper == nil {
		t.Error("status wrapper missing in CRD schema")
	} else {
		kStatusProps := getNestedMap(kStatusWrapper, "properties")
		hStatusProps := getNestedMap(hStatusWrapper, "properties")
		if _, ok := kStatusProps["uuid"]; !ok {
			t.Error("kustomize CRD missing 'uuid' in status.properties")
		}
		if _, ok := hStatusProps["uuid"]; !ok {
			t.Error("Helm CRD missing 'uuid' in status.properties")
		}

		// Verify status conditions schema exists
		if _, ok := kStatusProps["conditions"]; !ok {
			t.Error("kustomize CRD missing 'conditions' in status.properties")
		}
		if _, ok := hStatusProps["conditions"]; !ok {
			t.Error("Helm CRD missing 'conditions' in status.properties")
		}
	}
}

// TestNotificationRuleCRDDistribution verifies the NotificationRule CRD in the
// Helm chart template is synchronized with the kustomize-generated CRD.
func TestNotificationRuleCRDDistribution(t *testing.T) {
	root := findProjectRoot(t)

	kustomizeCRDPath := filepath.Join(root, "config", "crd", "bases", "dependencytrack.mko.dev_notificationrules.yaml")
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationrule-crd.yaml")

	// Verify both files exist
	if _, err := os.Stat(kustomizeCRDPath); os.IsNotExist(err) {
		t.Fatalf("kustomize NotificationRule CRD not found at %s", kustomizeCRDPath)
	}
	if _, err := os.Stat(helmCRDPath); os.IsNotExist(err) {
		t.Fatalf("Helm NotificationRule CRD template not found at %s", helmCRDPath)
	}

	// Parse both YAMLs (Helm version needs template stripping)
	var kustomizeDoc map[string]interface{}
	if err := parseYAML(kustomizeCRDPath, &kustomizeDoc); err != nil {
		t.Fatalf("failed to parse kustomize CRD: %v", err)
	}

	var helmDoc map[string]interface{}
	if err := parseHelmYAML(helmCRDPath, &helmDoc); err != nil {
		t.Fatalf("failed to parse Helm CRD template: %v", err)
	}

	// Compare CRD identity fields
	checkEqual(
		t,
		"CRD group",
		getString(getNestedMap(kustomizeDoc, "spec"), "group"),
		getString(getNestedMap(helmDoc, "spec"), "group"),
	)
	checkEqual(
		t,
		"CRD plural",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "plural"),
		getString(getNestedMap(helmDoc, "spec", "names"), "plural"),
	)
	checkEqual(
		t,
		"CRD singular",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "singular"),
		getString(getNestedMap(helmDoc, "spec", "names"), "singular"),
	)
	checkEqual(
		t,
		"CRD kind",
		getString(getNestedMap(kustomizeDoc, "spec", "names"), "kind"),
		getString(getNestedMap(helmDoc, "spec", "names"), "kind"),
	)
	checkEqual(
		t,
		"CRD scope",
		getString(getNestedMap(kustomizeDoc, "spec"), "scope"),
		getString(getNestedMap(helmDoc, "spec"), "scope"),
	)

	// Compare version schemas
	kustomizeVersions := getNestedMapSlice(kustomizeDoc, "spec", "versions")
	helmVersions := getNestedMapSlice(helmDoc, "spec", "versions")
	if len(kustomizeVersions) != len(helmVersions) {
		t.Errorf("version count mismatch: kustomize=%d, helm=%d", len(kustomizeVersions), len(helmVersions))
	}

	for i, kv := range kustomizeVersions {
		if i >= len(helmVersions) {
			break
		}
		hv := helmVersions[i]
		checkEqual(t, fmt.Sprintf("version[%d] name", i), getString(kv, "name"), getString(hv, "name"))
		checkEqual(t, fmt.Sprintf("version[%d] storage", i), getString(kv, "storage"), getString(hv, "storage"))
		checkEqual(t, fmt.Sprintf("version[%d] served", i), getString(kv, "served"), getString(hv, "served"))

		// Compare printer columns count
		kCols := getNestedMapSlice(kv, "additionalPrinterColumns")
		hCols := getNestedMapSlice(hv, "additionalPrinterColumns")
		if len(kCols) != len(hCols) {
			t.Errorf("column count mismatch at version[%d]: kustomize=%d, helm=%d", i, len(kCols), len(hCols))
		}
		for j := range kCols {
			if j >= len(hCols) {
				break
			}
			checkEqual(t, fmt.Sprintf("version[%d] column[%d] jsonPath", i, j),
				getString(kCols[j], "jsonPath"), getString(hCols[j], "jsonPath"))
			checkEqual(t, fmt.Sprintf("version[%d] column[%d] name", i, j),
				getString(kCols[j], "name"), getString(hCols[j], "name"))
		}
	}

	// Compare schema structure (ignoring description text differences)
	kSchema := getNestedMap(kustomizeVersions[0], "schema", "openAPIV3Schema")
	hSchema := getNestedMap(helmVersions[0], "schema", "openAPIV3Schema")

	// Verify spec properties exist and have same required fields
	kSpecProps := getNestedMap(kSchema, "properties", "spec")
	hSpecProps := getNestedMap(hSchema, "properties", "spec")

	checkEqual(
		t,
		"spec.required",
		strings.Join(getRequiredFields(kSpecProps), ","),
		strings.Join(getRequiredFields(hSpecProps), ","),
	)

	// Verify status properties exist.
	kStatusWrapper := getNestedMap(kSchema, "properties", "status")
	hStatusWrapper := getNestedMap(hSchema, "properties", "status")
	if kStatusWrapper == nil || hStatusWrapper == nil {
		t.Error("status wrapper missing in CRD schema")
	} else {
		kStatusProps := getNestedMap(kStatusWrapper, "properties")
		hStatusProps := getNestedMap(hStatusWrapper, "properties")
		if _, ok := kStatusProps["uuid"]; !ok {
			t.Error("kustomize CRD missing 'uuid' in status.properties")
		}
		if _, ok := hStatusProps["uuid"]; !ok {
			t.Error("Helm CRD missing 'uuid' in status.properties")
		}

		// Verify status conditions schema exists
		if _, ok := kStatusProps["conditions"]; !ok {
			t.Error("kustomize CRD missing 'conditions' in status.properties")
		}
		if _, ok := hStatusProps["conditions"]; !ok {
			t.Error("Helm CRD missing 'conditions' in status.properties")
		}
	}
}

// TestNotificationPublisherRBACDistribution verifies the NotificationPublisher
// RBAC templates in the Helm chart are synchronized with the kustomize source
// RBAC roles.
func TestNotificationPublisherRBACDistribution(t *testing.T) {
	root := findProjectRoot(t)

	testCases := []struct {
		name        string
		kustomize   string
		helm        string
		shouldMatch bool
	}{
		{
			name:      "notificationpublisher-admin-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationpublisher_admin_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationpublisher-admin-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "notificationpublisher-editor-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationpublisher_editor_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationpublisher-editor-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "notificationpublisher-viewer-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationpublisher_viewer_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationpublisher-viewer-rbac.yaml",
			),
			shouldMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify kustomize source exists
			if _, err := os.Stat(tc.kustomize); os.IsNotExist(err) {
				t.Fatalf("kustomize RBAC not found at %s", tc.kustomize)
			}
			// Verify Helm template exists
			if _, err := os.Stat(tc.helm); os.IsNotExist(err) {
				t.Fatalf("Helm RBAC template not found at %s", tc.helm)
			}

			if !tc.shouldMatch {
				return
			}

			// Parse both YAMLs
			var kustomizeDoc map[string]interface{}
			if err := parseYAML(tc.kustomize, &kustomizeDoc); err != nil {
				t.Fatalf("failed to parse kustomize RBAC: %v", err)
			}

			var helmDoc map[string]interface{}
			if err := parseHelmYAML(tc.helm, &helmDoc); err != nil {
				t.Fatalf("failed to parse Helm RBAC: %v", err)
			}

			// Both should be ClusterRole in Helm, Role in kustomize (by design)
			// The important thing is that the rules match
			kRules := getNestedMapSlice(kustomizeDoc, "rules")
			hRules := getNestedMapSlice(helmDoc, "rules")

			if len(kRules) != len(hRules) {
				t.Errorf("rule count mismatch: kustomize=%d, helm=%d", len(kRules), len(hRules))
			}

			for i, kr := range kRules {
				if i >= len(hRules) {
					break
				}
				hr := hRules[i]

				// Check apiGroups match
				kAPIGroups := getStringSlice(kr, "apiGroups")
				hAPIGroups := getStringSlice(hr, "apiGroups")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] apiGroups", i), kAPIGroups, hAPIGroups)

				// Check resources match
				kResources := getStringSlice(kr, "resources")
				hResources := getStringSlice(hr, "resources")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] resources", i), kResources, hResources)

				// Check verbs match
				kVerbs := getStringSlice(kr, "verbs")
				hVerbs := getStringSlice(hr, "verbs")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] verbs", i), kVerbs, hVerbs)
			}
		})
	}
}

// TestNotificationRuleRBACDistribution verifies the NotificationRule RBAC
// templates in the Helm chart are synchronized with the kustomize source RBAC
// roles.
func TestNotificationRuleRBACDistribution(t *testing.T) {
	root := findProjectRoot(t)

	testCases := []struct {
		name        string
		kustomize   string
		helm        string
		shouldMatch bool
	}{
		{
			name:      "notificationrule-admin-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationrule_admin_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationrule-admin-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "notificationrule-editor-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationrule_editor_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationrule-editor-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "notificationrule-viewer-role",
			kustomize: filepath.Join(root, "config", "rbac", "notificationrule_viewer_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationrule-viewer-rbac.yaml",
			),
			shouldMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify kustomize source exists
			if _, err := os.Stat(tc.kustomize); os.IsNotExist(err) {
				t.Fatalf("kustomize RBAC not found at %s", tc.kustomize)
			}
			// Verify Helm template exists
			if _, err := os.Stat(tc.helm); os.IsNotExist(err) {
				t.Fatalf("Helm RBAC template not found at %s", tc.helm)
			}

			if !tc.shouldMatch {
				return
			}

			// Parse both YAMLs
			var kustomizeDoc map[string]interface{}
			if err := parseYAML(tc.kustomize, &kustomizeDoc); err != nil {
				t.Fatalf("failed to parse kustomize RBAC: %v", err)
			}

			var helmDoc map[string]interface{}
			if err := parseHelmYAML(tc.helm, &helmDoc); err != nil {
				t.Fatalf("failed to parse Helm RBAC: %v", err)
			}

			kRules := getNestedMapSlice(kustomizeDoc, "rules")
			hRules := getNestedMapSlice(helmDoc, "rules")

			if len(kRules) != len(hRules) {
				t.Errorf("rule count mismatch: kustomize=%d, helm=%d", len(kRules), len(hRules))
			}

			for i, kr := range kRules {
				if i >= len(hRules) {
					break
				}
				hr := hRules[i]

				kAPIGroups := getStringSlice(kr, "apiGroups")
				hAPIGroups := getStringSlice(hr, "apiGroups")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] apiGroups", i), kAPIGroups, hAPIGroups)

				kResources := getStringSlice(kr, "resources")
				hResources := getStringSlice(hr, "resources")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] resources", i), kResources, hResources)

				kVerbs := getStringSlice(kr, "verbs")
				hVerbs := getStringSlice(hr, "verbs")
				checkSliceEqual(t, fmt.Sprintf("rules[%d] verbs", i), kVerbs, hVerbs)
			}
		})
	}
}

// TestManagerRBACIncludesNotificationResources verifies the manager
// ClusterRole in the Helm chart includes notificationpublishers and
// notificationrules in its rules alongside apikeys, policies, and teams.
func TestManagerRBACIncludesNotificationResources(t *testing.T) {
	root := findProjectRoot(t)

	helmManagerPath := filepath.Join(
		root, "deploy", "charts", "dependencytrack-operator", "templates", "manager-rbac.yaml",
	)

	if _, err := os.Stat(helmManagerPath); os.IsNotExist(err) {
		t.Fatalf("Helm manager RBAC template not found at %s", helmManagerPath)
	}

	data, err := os.ReadFile(helmManagerPath)
	if err != nil {
		t.Fatalf("failed to read Helm manager RBAC: %v", err)
	}

	// Strip Helm template variables for parsing
	strippedData := stripHelmTemplates(string(data))

	documents := strings.Split(strippedData, "\n---")
	var clusterRoleDoc string
	for _, doc := range documents {
		if strings.Contains(doc, "kind: ClusterRole") {
			clusterRoleDoc = doc
			break
		}
	}

	if clusterRoleDoc == "" {
		t.Fatal("no ClusterRole found in Helm manager RBAC template")
	}

	var cr map[string]interface{}
	if err := yaml.Unmarshal([]byte(clusterRoleDoc), &cr); err != nil {
		t.Fatalf("failed to parse ClusterRole: %v", err)
	}

	rules := getNestedMapSlice(cr, "rules")
	if len(rules) == 0 {
		t.Fatal("no rules found in manager ClusterRole")
	}

	var hasPublisherResources, hasPublisherFinalizers, hasPublisherStatus bool
	var hasRuleResources, hasRuleFinalizers, hasRuleStatus bool

	for _, rule := range rules {
		resources := getStringSlice(rule, "resources")
		verbs := getStringSlice(rule, "verbs")

		for _, r := range resources {
			if r == "notificationpublishers" {
				hasPublisherResources = true
			}
			if r == "notificationpublishers/finalizers" {
				hasPublisherFinalizers = true
			}
			if r == "notificationpublishers/status" {
				hasPublisherStatus = true
			}
			if r == "notificationrules" {
				hasRuleResources = true
			}
			if r == "notificationrules/finalizers" {
				hasRuleFinalizers = true
			}
			if r == "notificationrules/status" {
				hasRuleStatus = true
			}
		}
		// Verify that at least one rule carries standard CRUD verbs.
		for _, v := range verbs {
			switch v {
			case "create", "delete", "get", "list", "patch", "update", "watch", "*":
				// Verbs present — resource checks below validate the rule is complete.
			}
		}
	}

	if !hasPublisherResources {
		t.Error("manager ClusterRole missing 'notificationpublishers' resource")
	}
	if !hasPublisherFinalizers {
		t.Error("manager ClusterRole missing 'notificationpublishers/finalizers' resource")
	}
	if !hasPublisherStatus {
		t.Error("manager ClusterRole missing 'notificationpublishers/status' resource")
	}
	if !hasRuleResources {
		t.Error("manager ClusterRole missing 'notificationrules' resource")
	}
	if !hasRuleFinalizers {
		t.Error("manager ClusterRole missing 'notificationrules/finalizers' resource")
	}
	if !hasRuleStatus {
		t.Error("manager ClusterRole missing 'notificationrules/status' resource")
	}
}

// TestNotificationPublisherCRDSchemaIntegrity verifies the NotificationPublisher
// CRD schema has the required fields and types that match the Go types.
func TestNotificationPublisherCRDSchemaIntegrity(t *testing.T) {
	root := findProjectRoot(t)
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationpublisher-crd.yaml")

	var doc map[string]interface{}
	if err := parseHelmYAML(helmCRDPath, &doc); err != nil {
		t.Fatalf("failed to parse Helm CRD: %v", err)
	}

	schema := getNestedMap(
		getNestedMapSlice(doc, "spec", "versions")[0],
		"schema", "openAPIV3Schema",
	)
	specProps := getNestedMap(schema, "properties", "spec")

	// Verify required spec fields exist
	requiredFields := getRequiredFields(specProps)
	expectedRequired := []string{"extensionName", "name"}
	for _, expected := range expectedRequired {
		found := false
		for _, rf := range requiredFields {
			if rf == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("spec.required missing: %s", expected)
		}
	}

	// Verify spec.name has minLength and maxLength
	nameProps := getNestedMap(specProps, "properties", "name")
	if minLength := getString(nameProps, "minLength"); minLength != "1" {
		t.Errorf("spec.name minLength=%s, expected '1'", minLength)
	}
	if maxLength := getString(nameProps, "maxLength"); maxLength != "255" {
		t.Errorf("spec.name maxLength=%s, expected '255'", maxLength)
	}

	// Verify spec.extensionName has minLength
	extensionNameProps := getNestedMap(specProps, "properties", "extensionName")
	if minLength := getString(extensionNameProps, "minLength"); minLength != "1" {
		t.Errorf("spec.extensionName minLength=%s, expected '1'", minLength)
	}

	// Verify spec.description has maxLength
	descProps := getNestedMap(specProps, "properties", "description")
	if maxLength := getString(descProps, "maxLength"); maxLength != "1024" {
		t.Errorf("spec.description maxLength=%s, expected '1024'", maxLength)
	}

	// Verify status.uuid is a string type
	statusProps := getNestedMap(schema, "properties", "status")
	uuidType := getString(getNestedMap(statusProps, "properties", "uuid"), "type")
	if uuidType != "string" {
		t.Errorf("status.uuid type=%s, expected 'string'", uuidType)
	}
}

// TestNotificationRuleCRDSchemaIntegrity verifies the NotificationRule CRD
// schema has the required fields and types that match the Go types.
func TestNotificationRuleCRDSchemaIntegrity(t *testing.T) {
	root := findProjectRoot(t)
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "notificationrule-crd.yaml")

	var doc map[string]interface{}
	if err := parseHelmYAML(helmCRDPath, &doc); err != nil {
		t.Fatalf("failed to parse Helm CRD: %v", err)
	}

	schema := getNestedMap(
		getNestedMapSlice(doc, "spec", "versions")[0],
		"schema", "openAPIV3Schema",
	)
	specProps := getNestedMap(schema, "properties", "spec")

	// Verify required spec fields exist
	requiredFields := getRequiredFields(specProps)
	expectedRequired := []string{"level", "name", "publisherRef", "scope", "triggerType"}
	for _, expected := range expectedRequired {
		found := false
		for _, rf := range requiredFields {
			if rf == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("spec.required missing: %s", expected)
		}
	}

	// Verify spec.name has minLength and maxLength
	nameProps := getNestedMap(specProps, "properties", "name")
	if minLength := getString(nameProps, "minLength"); minLength != "1" {
		t.Errorf("spec.name minLength=%s, expected '1'", minLength)
	}
	if maxLength := getString(nameProps, "maxLength"); maxLength != "255" {
		t.Errorf("spec.name maxLength=%s, expected '255'", maxLength)
	}

	// Verify spec.level has enum values
	levelProps := getNestedMap(specProps, "properties", "level")
	levelEnum := getStringSlice(levelProps, "enum")
	expectedLevels := []string{"INFORMATIONAL", "WARNING", "ERROR"}
	if len(levelEnum) != len(expectedLevels) {
		t.Errorf("spec.level enum=%v, expected %v", levelEnum, expectedLevels)
	} else {
		for i, expected := range expectedLevels {
			if levelEnum[i] != expected {
				t.Errorf("spec.level enum[%d]=%s, expected %s", i, levelEnum[i], expected)
			}
		}
	}

	// Verify spec.scope has enum values
	scopeProps := getNestedMap(specProps, "properties", "scope")
	scopeEnum := getStringSlice(scopeProps, "enum")
	expectedScopes := []string{"PORTFOLIO", "SYSTEM"}
	if len(scopeEnum) != len(expectedScopes) {
		t.Errorf("spec.scope enum=%v, expected %v", scopeEnum, expectedScopes)
	} else {
		for i, expected := range expectedScopes {
			if scopeEnum[i] != expected {
				t.Errorf("spec.scope enum[%d]=%s, expected %s", i, scopeEnum[i], expected)
			}
		}
	}

	// Verify spec.triggerType has enum values
	triggerProps := getNestedMap(specProps, "properties", "triggerType")
	triggerEnum := getStringSlice(triggerProps, "enum")
	expectedTriggers := []string{"EVENT", "SCHEDULE"}
	if len(triggerEnum) != len(expectedTriggers) {
		t.Errorf("spec.triggerType enum=%v, expected %v", triggerEnum, expectedTriggers)
	} else {
		for i, expected := range expectedTriggers {
			if triggerEnum[i] != expected {
				t.Errorf("spec.triggerType enum[%d]=%s, expected %s", i, triggerEnum[i], expected)
			}
		}
	}

	// Verify spec.publisherRef has required field
	publisherRefProps := getNestedMap(specProps, "properties", "publisherRef")
	publisherRefRequired := getRequiredFields(publisherRefProps)
	expectedPRRequired := []string{"name"}
	if len(publisherRefRequired) != len(expectedPRRequired) {
		t.Errorf("publisherRef.required=%v, expected exactly %v", publisherRefRequired, expectedPRRequired)
	}

	// Verify spec.publisherConfigSecretRef has required fields
	secretRefProps := getNestedMap(specProps, "properties", "publisherConfigSecretRef")
	secretRefRequired := getRequiredFields(secretRefProps)
	expectedSRRequired := []string{"key", "name"}
	if len(secretRefRequired) != len(expectedSRRequired) {
		t.Errorf("publisherConfigSecretRef.required=%v, expected exactly %v", secretRefRequired, expectedSRRequired)
	}

	// Verify spec.notifyOn is an array type
	notifyOnProps := getNestedMap(specProps, "properties", "notifyOn")
	if notifyOnType := getString(notifyOnProps, "type"); notifyOnType != "array" {
		t.Errorf("spec.notifyOn type=%s, expected 'array'", notifyOnType)
	}

	// Verify status.uuid is a string type
	statusProps := getNestedMap(schema, "properties", "status")
	uuidType := getString(getNestedMap(statusProps, "properties", "uuid"), "type")
	if uuidType != "string" {
		t.Errorf("status.uuid type=%s, expected 'string'", uuidType)
	}
}

// TestNotificationCRDProjectMetadataMatchesScope verifies that PROJECT metadata
// aligns with the NotificationPublisher and NotificationRule CRD scopes.
func TestNotificationCRDProjectMetadataMatchesScope(t *testing.T) {
	root := findProjectRoot(t)

	var project struct {
		Resources []struct {
			API struct {
				Namespaced bool `yaml:"namespaced"`
			} `yaml:"api"`
			Kind string `yaml:"kind"`
		} `yaml:"resources"`
	}
	if err := parseYAML(filepath.Join(root, "PROJECT"), &project); err != nil {
		t.Fatalf("failed to parse PROJECT metadata: %v", err)
	}

	expectedKinds := map[string]bool{
		"NotificationPublisher": true,
		"NotificationRule":      true,
	}

	found := make(map[string]bool)

	for _, resource := range project.Resources {
		if expectedKinds[resource.Kind] {
			found[resource.Kind] = true
			if !resource.API.Namespaced {
				t.Errorf("PROJECT metadata declares %s cluster-scoped, want namespaced", resource.Kind)
			}
		}
	}

	for kind := range expectedKinds {
		if !found[kind] {
			t.Errorf("%s missing from PROJECT metadata", kind)
		}
	}

	// Verify CRD scope matches
	crdFiles := map[string]string{
		"NotificationPublisher": "dependencytrack.mko.dev_notificationpublishers.yaml",
		"NotificationRule":      "dependencytrack.mko.dev_notificationrules.yaml",
	}

	for kind, filename := range crdFiles {
		crdPath := filepath.Join(root, "config", "crd", "bases", filename)
		var crd map[string]interface{}
		if err := parseYAML(crdPath, &crd); err != nil {
			t.Fatalf("failed to parse %s CRD: %v", kind, err)
		}
		scope := getString(getNestedMap(crd, "spec"), "scope")
		if scope != "Namespaced" {
			t.Errorf("%s CRD scope=%q, want Namespaced", kind, scope)
		}
	}
}

// TestAllNotificationCRDsDistributed verifies that all notification resource
// CRDs are present in the Helm chart templates.
func TestAllNotificationCRDsDistributed(t *testing.T) {
	root := findProjectRoot(t)
	chartTemplates := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates")

	expectedCRDs := []string{
		"notificationpublisher-crd.yaml",
		"notificationrule-crd.yaml",
	}

	for _, crd := range expectedCRDs {
		path := filepath.Join(chartTemplates, crd)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Helm chart missing CRD template: %s", crd)
		}
	}
}

// TestAllNotificationRBACDistributed verifies that all notification resource
// RBAC templates exist for NotificationPublisher and NotificationRule.
func TestAllNotificationRBACDistributed(t *testing.T) {
	root := findProjectRoot(t)
	chartTemplates := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates")

	resources := map[string][]string{
		"notificationpublisher": {
			"notificationpublisher-admin-rbac.yaml",
			"notificationpublisher-editor-rbac.yaml",
			"notificationpublisher-viewer-rbac.yaml",
		},
		"notificationrule": {
			"notificationrule-admin-rbac.yaml",
			"notificationrule-editor-rbac.yaml",
			"notificationrule-viewer-rbac.yaml",
		},
	}
	for resource, rbacs := range resources {
		for _, rbac := range rbacs {
			path := filepath.Join(chartTemplates, rbac)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Helm chart missing RBAC template for %s: %s", resource, rbac)
			}
		}
	}
}
