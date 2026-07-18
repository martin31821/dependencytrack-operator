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
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// helmTemplatePattern matches Go template expressions like {{- include "..." . | nindent 4 }}
var helmTemplatePattern = regexp.MustCompile(`\{\{.*?\}\}`)

// stripHelmTemplates replaces Go template expressions with empty strings for YAML parsing.
func stripHelmTemplates(content string) string {
	// Replace multi-line template expressions first
	multiLine := regexp.MustCompile(`\{\{[-\s\S]*?\}\}`)
	result := multiLine.ReplaceAllString(content, "")
	// Then single-line
	result = helmTemplatePattern.ReplaceAllString(result, "")
	return result
}

// readAndStrip reads a file and strips Helm template variables for parsing.
func readAndStrip(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	stripped := stripHelmTemplates(string(data))
	return []byte(stripped), nil
}

// parseHelmYAML parses a Helm template file by stripping template variables first.
func parseHelmYAML(path string, out interface{}) error {
	data, err := readAndStrip(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

// TestPolicyCRDDistribution verifies the Policy CRD in the Helm chart
// template is synchronized with the kustomize-generated CRD.
func TestPolicyCRDDistribution(t *testing.T) {
	root := findProjectRoot(t)

	kustomizeCRDPath := filepath.Join(root, "config", "crd", "bases", "dependencytrack.mko.dev_policies.yaml")
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "policy-crd.yaml")

	// Verify both files exist
	if _, err := os.Stat(kustomizeCRDPath); os.IsNotExist(err) {
		t.Fatalf("kustomize Policy CRD not found at %s", kustomizeCRDPath)
	}
	if _, err := os.Stat(helmCRDPath); os.IsNotExist(err) {
		t.Fatalf("Helm Policy CRD template not found at %s", helmCRDPath)
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
	// The CRD uses nested properties: status -> properties -> properties{uuid, conditions}.
	// getNestedMap("properties", "status") returns {type, description, properties}.
	// We need to descend one more level into .properties for the field keys.
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

// TestPolicyRBACDistribution verifies the Policy RBAC templates in the Helm chart
// are synchronized with the kustomize source RBAC roles.
func TestPolicyRBACDistribution(t *testing.T) {
	root := findProjectRoot(t)

	testCases := []struct {
		name        string
		kustomize   string
		helm        string
		shouldMatch bool
	}{
		{
			name:      "policy-admin-role",
			kustomize: filepath.Join(root, "config", "rbac", "policy_admin_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "policy-admin-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "policy-editor-role",
			kustomize: filepath.Join(root, "config", "rbac", "policy_editor_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "policy-editor-rbac.yaml",
			),
			shouldMatch: true,
		},
		{
			name:      "policy-viewer-role",
			kustomize: filepath.Join(root, "config", "rbac", "policy_viewer_role.yaml"),
			helm: filepath.Join(
				root, "deploy", "charts", "dependencytrack-operator", "templates", "policy-viewer-rbac.yaml",
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

// TestManagerRBACIsClusterScoped verifies the kustomize deployment grants the
// manager cluster-wide access because the manager watches resources in all namespaces.
func TestManagerRBACIsClusterScoped(t *testing.T) {
	root := findProjectRoot(t)

	var role map[string]interface{}
	if err := parseYAML(filepath.Join(root, "config", "rbac", "role.yaml"), &role); err != nil {
		t.Fatalf("failed to parse manager RBAC role: %v", err)
	}
	if kind := getString(role, "kind"); kind != "ClusterRole" {
		t.Errorf("manager RBAC kind = %q, want ClusterRole", kind)
	}

	var binding map[string]interface{}
	if err := parseYAML(filepath.Join(root, "config", "rbac", "role_binding.yaml"), &binding); err != nil {
		t.Fatalf("failed to parse manager RBAC binding: %v", err)
	}
	if kind := getString(binding, "kind"); kind != "ClusterRoleBinding" {
		t.Errorf("manager RBAC binding kind = %q, want ClusterRoleBinding", kind)
	}
	if roleRefKind := getString(getNestedMap(binding, "roleRef"), "kind"); roleRefKind != "ClusterRole" {
		t.Errorf("manager RBAC roleRef kind = %q, want ClusterRole", roleRefKind)
	}
}

// TestPolicyProjectMetadataMatchesCRDScope verifies Kubebuilder metadata stays
// aligned with the generated namespaced Policy CRD.
func TestPolicyProjectMetadataMatchesCRDScope(t *testing.T) {
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

	policyFound := false
	for _, resource := range project.Resources {
		if resource.Kind == "Policy" {
			policyFound = true
			if !resource.API.Namespaced {
				t.Error("PROJECT metadata declares Policy cluster-scoped, want namespaced")
			}
			break
		}
	}
	if !policyFound {
		t.Fatal("Policy resource missing from PROJECT metadata")
	}

	var crd map[string]interface{}
	policyCRDPath := filepath.Join(
		root, "config", "crd", "bases", "dependencytrack.mko.dev_policies.yaml",
	)
	if err := parseYAML(policyCRDPath, &crd); err != nil {
		t.Fatalf("failed to parse Policy CRD: %v", err)
	}
	if scope := getString(getNestedMap(crd, "spec"), "scope"); scope != "Namespaced" {
		t.Errorf("Policy CRD scope = %q, want Namespaced", scope)
	}
}

// TestManagerRBACIncludesPolicy verifies the manager ClusterRole in the Helm chart
// includes policies in its rules alongside apikeys and teams.
func TestManagerRBACIncludesPolicy(t *testing.T) {
	root := findProjectRoot(t)

	helmManagerPath := filepath.Join(
		root, "deploy", "charts", "dependencytrack-operator", "templates", "manager-rbac.yaml",
	)

	if _, err := os.Stat(helmManagerPath); os.IsNotExist(err) {
		t.Fatalf("Helm manager RBAC template not found at %s", helmManagerPath)
	}

	// Parse the YAML
	// The file has two documents: ClusterRole and ClusterRoleBinding
	data, err := os.ReadFile(helmManagerPath)
	if err != nil {
		t.Fatalf("failed to read Helm manager RBAC: %v", err)
	}

	// Strip Helm template variables for parsing
	strippedData := stripHelmTemplates(string(data))

	// Find the ClusterRole document
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

	// Check that policies appear in the rules
	var hasPoliciesResources, hasPoliciesVerbs, hasPoliciesFinalizers, hasPoliciesStatus bool
	for _, rule := range rules {
		resources := getStringSlice(rule, "resources")
		verbs := getStringSlice(rule, "verbs")

		for _, r := range resources {
			if r == "policies" {
				hasPoliciesResources = true
			}
			if r == "policies/finalizers" {
				hasPoliciesFinalizers = true
			}
			if r == "policies/status" {
				hasPoliciesStatus = true
			}
		}
		// Verify that at least one rule carries standard CRUD verbs.
		// In this RBAC, 'policies' resources always carry create/delete/get/list/patch/update/watch.
		for _, v := range verbs {
			switch v {
			case "create", "delete", "get", "list", "patch", "update", "watch", "*":
				hasPoliciesVerbs = true
			}
		}
	}

	if !hasPoliciesResources {
		t.Error("manager ClusterRole missing 'policies' resource")
	}
	if !hasPoliciesFinalizers {
		t.Error("manager ClusterRole missing 'policies/finalizers' resource")
	}
	if !hasPoliciesStatus {
		t.Error("manager ClusterRole missing 'policies/status' resource")
	}
	if !hasPoliciesVerbs {
		t.Error("manager ClusterRole missing policies verbs (create/delete/get/list/patch/update/watch)")
	}
}

// TestPolicyCRDSchemaIntegrity verifies the Policy CRD schema has the required fields
// and types that match the Go types.
func TestPolicyCRDSchemaIntegrity(t *testing.T) {
	root := findProjectRoot(t)
	helmCRDPath := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates", "policy-crd.yaml")

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
	expectedRequired := []string{"name", "priority", "failureAction", "conditions"}
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

	// Verify conditions has minItems
	conditionsProps := getNestedMap(specProps, "properties", "conditions")
	if minItems := getString(conditionsProps, "minItems"); minItems != "1" {
		t.Errorf("conditions.minItems=%s, expected '1'", minItems)
	}

	// Verify conditions items have required fields
	items := getNestedMap(conditionsProps, "items")
	condRequired := getRequiredFields(items)
	condExpected := []string{"type", "comparator", "value"}
	for _, expected := range condExpected {
		found := false
		for _, rf := range condRequired {
			if rf == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("conditions[].required missing: %s", expected)
		}
	}

	// Verify status.uuid is a string type
	statusProps := getNestedMap(schema, "properties", "status")
	uuidType := getString(getNestedMap(statusProps, "properties", "uuid"), "type")
	if uuidType != "string" {
		t.Errorf("status.uuid type=%s, expected 'string'", uuidType)
	}
}

// TestAllCRDsDistributed verifies that all three resource CRDs (APIKey, Policy, Team)
// are present in the Helm chart templates.
func TestAllCRDsDistributed(t *testing.T) {
	root := findProjectRoot(t)
	chartTemplates := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates")

	expectedCRDs := []string{
		"apikey-crd.yaml",
		"policy-crd.yaml",
		"team-crd.yaml",
	}

	for _, crd := range expectedCRDs {
		path := filepath.Join(chartTemplates, crd)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Helm chart missing CRD template: %s", crd)
		}
	}
}

// TestAllRBACDistributed verifies that all three resource RBAC templates exist for Policy.
func TestAllRBACDistributed(t *testing.T) {
	root := findProjectRoot(t)
	chartTemplates := filepath.Join(root, "deploy", "charts", "dependencytrack-operator", "templates")

	// Policy RBAC templates
	policyRBACs := []string{
		"policy-admin-rbac.yaml",
		"policy-editor-rbac.yaml",
		"policy-viewer-rbac.yaml",
	}
	for _, rbac := range policyRBACs {
		path := filepath.Join(chartTemplates, rbac)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Helm chart missing RBAC template: %s", rbac)
		}
	}

	// Verify all three RBAC types exist for every resource
	resources := map[string][]string{
		"apikey": {"apikey-admin-rbac.yaml", "apikey-editor-rbac.yaml", "apikey-viewer-rbac.yaml"},
		"policy": {"policy-admin-rbac.yaml", "policy-editor-rbac.yaml", "policy-viewer-rbac.yaml"},
		"team":   {"team-admin-rbac.yaml", "team-editor-rbac.yaml", "team-viewer-rbac.yaml"},
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

// --- helpers ---

func findProjectRoot(t *testing.T) string {
	// Walk up from this test file to find the project root (go.mod)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("could not find project root (go.mod not found)")
	return ""
}

func parseYAML(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func checkEqual(t *testing.T, label, a, b string) {
	if a != b {
		t.Errorf("%s: got %q, want %q", label, a, b)
	}
}

func checkSliceEqual(t *testing.T, label string, a, b []string) {
	if len(a) != len(b) {
		t.Errorf("%s: got %v, want %v", label, a, b)
		return
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("%s[%d]: got %q, want %q", label, i, a[i], b[i])
		}
	}
}

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if m == nil {
		return nil
	}
	if v, ok := m[key]; ok {
		var result []string
		for _, item := range v.([]interface{}) {
			result = append(result, fmt.Sprintf("%v", item))
		}
		return result
	}
	return nil
}

func getNestedMap(m map[string]interface{}, keys ...string) map[string]interface{} {
	for _, key := range keys {
		if m == nil {
			return nil
		}
		if v, ok := m[key]; ok {
			if nm, ok := v.(map[string]interface{}); ok {
				m = nm
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return m
}

func getNestedMapSlice(m map[string]interface{}, keys ...string) []map[string]interface{} {
	// Walk the path, collecting intermediate maps.
	// If we land on an array at the final key, return it.
	// If we land on a map, look for "items" or any []interface{} values.
	node := m
	for i, key := range keys {
		if node == nil {
			return nil
		}
		v, ok := node[key]
		if !ok {
			return nil
		}
		// If this is not the last key, we need a map to descend further.
		if i < len(keys)-1 {
			if _, isArr := v.([]interface{}); isArr {
				// Intermediate value is an array but we need a map.
				return nil
			}
			if nm, ok := v.(map[string]interface{}); ok {
				node = nm
			} else {
				return nil
			}
		} else {
			// This is the last key — check for array directly
			if arr, ok := v.([]interface{}); ok {
				return toMapSlice(arr)
			}
			// Fallback: look inside a map for "items" key
			if nm, ok := v.(map[string]interface{}); ok {
				if items, ok := nm["items"]; ok {
					if arr, ok := items.([]interface{}); ok {
						return toMapSlice(arr)
					}
				}
				// Also scan for any direct []interface{} values
				for _, val := range nm {
					if arr, ok := val.([]interface{}); ok {
						return toMapSlice(arr)
					}
				}
			}
			return nil
		}
	}
	// All keys exhausted and node is a map — look for arrays
	if node == nil {
		return nil
	}
	for _, v := range node {
		if arr, ok := v.([]interface{}); ok {
			return toMapSlice(arr)
		}
	}
	return nil
}

func toMapSlice(arr []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range arr {
		if nm, ok := item.(map[string]interface{}); ok {
			result = append(result, nm)
		}
	}
	return result
}

func getRequiredFields(props map[string]interface{}) []string {
	if props == nil {
		return nil
	}
	required := getStringSlice(props, "required")
	return required
}
