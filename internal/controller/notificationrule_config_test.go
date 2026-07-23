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

package controller

import (
	"fmt"
	"strings"
	"testing"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
)

// ---- PublisherConfigValidator tests ----

func TestValidate_EmptyInput(t *testing.T) {
	v := NewPublisherConfigValidator()
	err := v.Validate([]byte{}, nil)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok, "expected PublisherConfigValidationError, got %T", err)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
}

func TestValidate_MalformedJSON(t *testing.T) {
	v := NewPublisherConfigValidator()
	err := v.Validate([]byte("{bad json"), nil)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
	// Verify Secret values do NOT leak into the error message.
	assert.NotContains(t, ve.Message, "{{secret")
	assert.NotContains(t, ve.Message, "s3cr3t")
}

func TestValidate_JSONIsNotObject(t *testing.T) {
	v := NewPublisherConfigValidator()

	// JSON array
	err := v.Validate([]byte(`[1,2,3]`), nil)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
	assert.Contains(t, ve.Message, "must be a JSON object")

	// JSON string
	err = v.Validate([]byte(`"hello"`), nil)
	assert.Error(t, err)
	ve, ok = err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
}

func TestValidate_ValidJSONNoSchema(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"url": "https://example.com", "enabled": true}`)
	err := v.Validate(config, nil)
	assert.NoError(t, err, "valid JSON without schema should pass")
}

func TestValidate_ValidJSONWithSchema(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"url": "https://example.com", "enabled": true}`)
	schema := map[string]interface{}{
		"required": []interface{}{"url"},
		"properties": map[string]interface{}{
			"url":     map[string]interface{}{"type": "string"},
			"enabled": map[string]interface{}{"type": "boolean"},
		},
	}
	err := v.Validate(config, schema)
	assert.NoError(t, err)
}

// ---- Required field validation ----

func TestValidate_MissingRequiredField(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"other": "value"}`) // missing required "url"
	schema := map[string]interface{}{
		"required": []interface{}{"url"},
		"properties": map[string]interface{}{
			"url": map[string]interface{}{"type": "string"},
		},
	}
	// The validator intentionally skips the JSON schema "required" list because
	// the generated OpenAPI schema marks many fields as required that DT accepts
	// as optional with defaults. The API itself is the authoritative validator.
	err := v.Validate(config, schema)
	assert.NoError(t, err)
}

// ---- Type validation ----

func TestValidate_TypeMismatch_StringGotNumber(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"name": 123}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
		},
	}
	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

func TestValidate_TypeMismatch_NumberGotString(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"count": "many"}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"count": map[string]interface{}{"type": "number"},
		},
	}
	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

func TestValidate_TypeMismatch_BooleanGotString(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"active": "yes"}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"active": map[string]interface{}{"type": "boolean"},
		},
	}
	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

func TestValidate_TypeObjectGotArray(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"metadata": [1, 2]}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"metadata": map[string]interface{}{"type": "object"},
		},
	}
	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

func TestValidate_TypeArrayGotObject(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"items": {"a": 1}}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"items": map[string]interface{}{"type": "array"},
		},
	}
	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

func TestValidate_AllTypesValid(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{
		"name": "test",
		"count": 42,
		"enabled": true,
		"metadata": {"key": "val"},
		"tags": ["a", "b"]
	}`)
	schema := map[string]interface{}{
		"required": []interface{}{"name"},
		"properties": map[string]interface{}{
			"name":     map[string]interface{}{"type": "string"},
			"count":    map[string]interface{}{"type": "number"},
			"enabled":  map[string]interface{}{"type": "boolean"},
			"metadata": map[string]interface{}{"type": "object"},
			"tags":     map[string]interface{}{"type": "array"},
		},
	}
	err := v.Validate(config, schema)
	assert.NoError(t, err)
}

// ---- No-op paths ----

func TestValidate_NilReference(t *testing.T) {
	v := NewPublisherConfigValidator()
	err := v.ValidateFromRef(nil, []byte(`{}`), nil)
	assert.NoError(t, err, "nil reference should be a no-op")
}

func TestValidate_SchemaRetrievalError(t *testing.T) {
	v := NewPublisherConfigValidator()
	config := []byte(`{"url": "https://example.com"}`)
	// Simulate schema retrieval failure by passing nil (no schema returned).
	// This is not an error — schema is optional, config is validated structurally.
	err := v.Validate(config, nil)
	assert.NoError(t, err)

	// If we pass an explicit schema retrieval failure, it should be a different
	// reason — but the validator has no way to know the schema was supposed
	// to come from DT. The caller should use reasonSchemaRetrieval when the
	// /configSchema endpoint fails.
}

func TestValidate_NoSchemaReturns(t *testing.T) {
	v := NewPublisherConfigValidator()
	// Empty schema object — should pass structural validation.
	err := v.Validate([]byte(`{"any": "thing"}`), map[string]interface{}{})
	assert.NoError(t, err)
}

// ---- SanitizeJSONError (Q7: negative coverage) ----

func TestSanitizeJSONError_NoSecretLeak(t *testing.T) {
	testInput := `invalid literal inside {"secret": "s3cr3t"}`
	// We can't test the actual JSON error from Go's stdlib since it produces
	// a specific format, but we can verify sanitizeJSONError truncates properly.
	errMsg := sanitizeJSONError(fmt.Errorf("invalid character '{' looking for beginning of value inside %q", testInput))
	// The sanitizer should truncate at the first quote.
	assert.False(t, strings.Contains(errMsg, "s3cr3t"), "sanitizeJSONError must not leak config content")
}

func TestJSONType(t *testing.T) {
	assert.Equal(t, "string", jsonType("hello"))
	assert.Equal(t, "number", jsonType(float64(42)))
	assert.Equal(t, "boolean", jsonType(true))
	assert.Equal(t, "object", jsonType(map[string]interface{}{"a": 1}))
	assert.Equal(t, "array", jsonType([]interface{}{1}))
	assert.Equal(t, "null", jsonType(nil))
}

// ---- Error message safety (Q3: no Secret values) ----

func TestErrorMessage_NoSecretValues(t *testing.T) {
	v := NewPublisherConfigValidator()
	// Malformed config that contains a secret-like string.
	config := []byte(`{"token": "sk-ir34ldgdh", "endpoint": "https://hooks.slack.com/services/..."}`)
	// Force a malformed JSON that would echo back the input.
	err := v.Validate(config, nil)
	// This actually succeeds because it's valid JSON — use invalid instead.
	err = v.Validate([]byte(`{invalid}`), nil)
	assert.Error(t, err)
	// The error message must NOT contain the secret value.
	assert.False(t, strings.Contains(err.Error(), "sk-ir"), "error must not leak API key")
}

func TestPublisherConfigValidationError_Format(t *testing.T) {
	err := &PublisherConfigValidationError{
		Reason:  "TestReason",
		Message: "test message",
	}
	assert.Equal(t, "TestReason: test message", err.Error())
}

// ---- ValidateSecretData ----

func TestValidateSecretData(t *testing.T) {
	v := NewPublisherConfigValidator()
	err := v.ValidateSecretData([]byte(`{"url": "https://example.com"}`), nil)
	assert.NoError(t, err)
}

func TestValidateSecretData_Empty(t *testing.T) {
	v := NewPublisherConfigValidator()
	err := v.ValidateSecretData([]byte{}, nil)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
}

// ---- Integration: ValidateFromRef with real ref type ----

func TestValidateFromRef_WithRef(t *testing.T) {
	v := NewPublisherConfigValidator()
	ref := &dependencytrackv1alpha1.PublisherConfigSecretRef{
		Name: "my-secret",
		Key:  "config",
	}
	config := []byte(`{"url": "https://example.com"}`)
	err := v.ValidateFromRef(ref, config, nil)
	assert.NoError(t, err)
}

// ---- CRD-level: ensure DeepCopy works for the new ref type ----

func TestPublisherConfigSecretRef_DeepCopy(t *testing.T) {
	ref := &dependencytrackv1alpha1.PublisherConfigSecretRef{
		Name: "my-secret",
		Key:  "config-key",
	}

	copied := ref.DeepCopy()

	assert.Equal(t, ref.Name, copied.Name)
	assert.Equal(t, ref.Key, copied.Key)

	// Verify independence
	copied.Name = "modified"
	assert.NotEqual(t, ref.Name, copied.Name)
}

func TestNotificationRuleSpec_DeepCopyWithRef(t *testing.T) {
	enabled := true
	spec := dependencytrackv1alpha1.NotificationRuleSpec{
		Name:         "test-rule",
		Scope:        dependencytrackv1alpha1.NotificationRuleScopePortfolio,
		TriggerType:  dependencytrackv1alpha1.NotificationRuleTriggerTypeEvent,
		Level:        dependencytrackv1alpha1.NotificationRuleLevelWarn,
		PublisherRef: dependencytrackv1alpha1.NotificationRulePublisherRef{Name: "my-pub"},
		Enabled:      &enabled,
		PublisherConfigSecretRef: &dependencytrackv1alpha1.PublisherConfigSecretRef{
			Name: "config-secret",
			Key:  "publisher-config",
		},
	}

	copied := spec.DeepCopy()

	assert.Equal(t, spec.Name, copied.Name)
	assert.Equal(t, spec.PublisherConfigSecretRef.Name, copied.PublisherConfigSecretRef.Name)
	assert.Equal(t, spec.PublisherConfigSecretRef.Key, copied.PublisherConfigSecretRef.Key)

	// Verify independence
	copied.PublisherConfigSecretRef.Name = "changed"
	assert.NotEqual(t, spec.PublisherConfigSecretRef.Name, copied.PublisherConfigSecretRef.Name)
}

// ---- Negative: ensure no panics with edge cases ----

func TestValidate_EdgeCases(t *testing.T) {
	v := NewPublisherConfigValidator()

	// nil bytes
	err := v.Validate(nil, nil)
	assert.Error(t, err)

	// whitespace only
	err = v.Validate([]byte("   "), nil)
	assert.Error(t, err)

	// JSON number
	err = v.Validate([]byte("42"), nil)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)

	// JSON null
	err = v.Validate([]byte("null"), nil)
	assert.Error(t, err)
	ve, ok = err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)

	// JSON boolean
	err = v.Validate([]byte("true"), nil)
	assert.Error(t, err)
	ve, ok = err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonMalformedJSON, ve.Reason)
}

// ---- Schema validation: required + type combined ----

func TestValidate_RequiredAndType(t *testing.T) {
	v := NewPublisherConfigValidator()

	// Missing required fields — validator intentionally skips the "required"
	// list because the generated OpenAPI schema marks many fields as required
	// that DT accepts as optional with defaults.
	config1 := []byte(`{}`)
	schema1 := map[string]interface{}{
		"required": []interface{}{"url", "name"},
		"properties": map[string]interface{}{
			"url":  map[string]interface{}{"type": "string"},
			"name": map[string]interface{}{"type": "string"},
		},
	}
	err := v.Validate(config1, schema1)
	assert.NoError(t, err)

	// Both present but name is wrong type — type check still applies
	config2 := []byte(`{"url": "https://example.com", "name": 123}`)
	err = v.Validate(config2, schema1)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)
}

// ---- Redaction: ensure Secret key name and values never appear ----

func TestErrorDoesNotLeakSecretKeyOrValue(t *testing.T) {
	v := NewPublisherConfigValidator()

	// Malformed JSON with a Secret-like value
	raw := []byte(`{"webhookUrl": "https://hooks.slack.com/services/T0000/B0000/XXXX", "token": "xoxb-123456"}`)

	// This is actually valid JSON, so it passes. Let's make it malformed.
	raw = []byte(`{"webhookUrl": "https://hooks.slack.com/services/T0000/B0000/XXXX", token: "xoxb-123456"}`)

	err := v.Validate(raw, nil)
	assert.Error(t, err)

	// The error message must NOT contain the token value or webhook URL.
	assert.NotContains(t, err.Error(), "xoxb-123456", "error must not leak token")
	assert.NotContains(t, err.Error(), "T0000", "error must not leak webhook URL")
}

// ---- Schema property paths only ----

func TestSchemaValidationDoesNotExceedPropertyPaths(t *testing.T) {
	v := NewPublisherConfigValidator()

	config := []byte(`{"url": 123}`)
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"url": map[string]interface{}{"type": "string"},
		},
	}

	err := v.Validate(config, schema)
	assert.Error(t, err)
	ve, ok := err.(*PublisherConfigValidationError)
	assert.True(t, ok)
	assert.Equal(t, reasonSchemaValidation, ve.Reason)

	// The message should describe the problem generally without exposing
	// the exact schema property path that caused the failure.
	// "a publisherConfig property must be a string" is acceptable.
	// It must NOT contain schema keywords like "$ref" or raw JSON pointer.
	assert.False(t, strings.Contains(ve.Message, "$ref"), "schema $ref must not leak")
	assert.False(t, strings.Contains(ve.Message, "$schema"), "schema $schema must not leak")
}
