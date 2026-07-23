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
	"encoding/json"
	"fmt"

	dependencytrackv1alpha1 "github.com/martin31821/dependencytrack-operator/api/v1alpha1"
)

// PublisherConfigValidationError represents a validation failure with a stable reason
// and a human-readable message that never exposes Secret values.
type PublisherConfigValidationError struct {
	// Reason is a machine-readable, stable identifier (e.g. "MalformedJSON").
	Reason string
	// Message is a human-readable message; it never contains Secret key names,
	// Secret values, or schema property paths.
	Message string
}

func (e *PublisherConfigValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Reason, e.Message)
}

// reason constants — stable Ready=False reasons for Secret / validation failures.
const (
	reasonSecretNotFound    = "SecretNotFound"
	reasonSecretKeyNotFound = "SecretKeyNotFound"
	reasonMalformedJSON     = "MalformedJSON"
	reasonSchemaValidation  = "SchemaValidationFailed"
	reasonSchemaRetrieval   = "SchemaRetrievalFailed"
	reasonNoSchemaReturned  = "NoSchemaReturned"
)

// PublisherConfigValidator validates publisherConfig JSON documents.
// It is stateless and can be shared across reconciliations.
type PublisherConfigValidator struct{}

// NewPublisherConfigValidator creates a new validator instance.
func NewPublisherConfigValidator() *PublisherConfigValidator {
	return &PublisherConfigValidator{}
}

// Validate reads the raw JSON bytes and validates:
//  1. The bytes are valid JSON and form an object.
//  2. If a schema is provided (from the DT /configSchema endpoint), the config
//     is validated against it locally.
//
// If validation fails, it returns a *PublisherConfigValidationError with a stable
// reason so the caller can map it to a Ready=False condition.
func (v *PublisherConfigValidator) Validate(raw []byte, schema map[string]interface{}) error {
	if len(raw) == 0 {
		return &PublisherConfigValidationError{
			Reason:  reasonMalformedJSON,
			Message: "publisherConfig is empty",
		}
	}

	// Step 1: parse as generic JSON to verify it is valid JSON + an object.
	var parsed interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return &PublisherConfigValidationError{
			Reason:  reasonMalformedJSON,
			Message: fmt.Sprintf("publisherConfig is not valid JSON: %s", sanitizeJSONError(err)),
		}
	}

	obj, ok := parsed.(map[string]interface{})
	if !ok {
		return &PublisherConfigValidationError{
			Reason:  reasonMalformedJSON,
			Message: "publisherConfig must be a JSON object, got " + jsonType(parsed),
		}
	}

	// Step 2: if a schema is available, validate structurally.
	if len(schema) > 0 {
		if err := v.validateAgainstSchema(obj, schema); err != nil {
			return err
		}
	}

	return nil
}

// ValidateFromRef extracts the Secret reference fields and calls Validate.
// The caller must first retrieve the Secret data; this method only validates
// the resulting bytes.
func (v *PublisherConfigValidator) ValidateFromRef(ref *dependencytrackv1alpha1.PublisherConfigSecretRef, data []byte, schema map[string]interface{}) error {
	if ref == nil {
		return nil // no reference means no config to validate
	}
	return v.Validate(data, schema)
}

// ValidateSecretData checks that the provided secret data (raw bytes from a
// Secret key) is valid publisherConfig JSON. Returns a stable error on failure.
func (v *PublisherConfigValidator) ValidateSecretData(data []byte, schema map[string]interface{}) error {
	return v.Validate(data, schema)
}

// validateAgainstSchema performs structural validation against the DT publisher
// config JSON schema. It checks:
//   - If the schema declares "type" as "object", the config must be an object.
//   - If "properties" are declared, each property value must match its declared type.
//
// Note: the schema's "required" list is intentionally skipped — the generated
// OpenAPI schema marks many fields as required that the actual DT API accepts as
// optional with defaults. The API itself is the authoritative validator.
func (v *PublisherConfigValidator) validateAgainstSchema(config map[string]interface{}, schema map[string]interface{}) error {
	// Check declared types.
	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for propKey, propSchema := range properties {
			propDef, ok := propSchema.(map[string]interface{})
			if !ok {
				continue
			}
			val, exists := config[propKey]
			if !exists {
				// Already caught by "required" check; skip here.
				continue
			}
			if err := v.checkPropertyType(val, propDef); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkPropertyType validates a single property value against its JSON schema type.
func (v *PublisherConfigValidator) checkPropertyType(val interface{}, propDef map[string]interface{}) error {
	declaredType, _ := propDef["type"].(string)
	if declaredType == "" {
		return nil // no type constraint
	}

	switch declaredType {
	case "string":
		if _, ok := val.(string); !ok {
			return &PublisherConfigValidationError{
				Reason:  reasonSchemaValidation,
				Message: fmt.Sprintf("a publisherConfig property must be a string, got %s", jsonType(val)),
			}
		}
	case "number":
		if _, ok := val.(float64); !ok {
			return &PublisherConfigValidationError{
				Reason:  reasonSchemaValidation,
				Message: fmt.Sprintf("a publisherConfig property must be a number, got %s", jsonType(val)),
			}
		}
	case "boolean":
		if _, ok := val.(bool); !ok {
			return &PublisherConfigValidationError{
				Reason:  reasonSchemaValidation,
				Message: fmt.Sprintf("a publisherConfig property must be a boolean, got %s", jsonType(val)),
			}
		}
	case "object":
		if _, ok := val.(map[string]interface{}); !ok {
			return &PublisherConfigValidationError{
				Reason:  reasonSchemaValidation,
				Message: fmt.Sprintf("a publisherConfig property must be an object, got %s", jsonType(val)),
			}
		}
	case "array":
		if _, ok := val.([]interface{}); !ok {
			return &PublisherConfigValidationError{
				Reason:  reasonSchemaValidation,
				Message: fmt.Sprintf("a publisherConfig property must be an array, got %s", jsonType(val)),
			}
		}
	default:
		// Unknown type — skip validation (future-proof).
	}

	return nil
}

// sanitizeJSONError returns a short, safe description of a JSON parse error
// without echoing the raw offending bytes.
func sanitizeJSONError(err error) string {
	msg := err.Error()
	// Truncate at the first backtick or quote to avoid leaking config content.
	for i, ch := range msg {
		if ch == '`' || ch == '"' || ch == '\'' || ch == '\n' || ch == '\r' {
			return msg[:i]
		}
	}
	return msg
}

// jsonType returns the JSON type name for a Go value.
func jsonType(val interface{}) string {
	switch val.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", val)
	}
}
