/*
OWASP Dependency-Track

# REST API of OWASP Dependency-Track  ## Pagination  This API implements token-based pagination. Collection responses have the following structure:  ```json {   \"items\": [...],   \"next_page_token\": \"abcdefg\",   \"total\": {     \"count\": 100,     \"type\": \"EXACT\"   } } ```  `next_page_token` is present when more items exist, and absent otherwise. To fetch the next page, pass it as the `page_token` query parameter.  To navigate backwards, clients should keep track of previous page tokens as they paginate through collections. The API does *not* provide backward navigation!  Collections that support sorting will only consider the `sort_by` and `sort_direction` query parameters for the request of the first page. For subsequent pages, sorting preferences are bound to the page token.  Page tokens are opaque strings. Clients should not try to interpret or generate them. Their format may change without notice.  The `total` object discloses how many items exist in the collection *across all pages*. Because counting is expensive, some collections that hold *a lot* of items may return partial counts (type `AT_LEAST`) instead of exact counts (type `EXACT`). Which type to expect is usually documented in the operation's description.  ## Sorting  Items in a collection can be sorted using the `sort_by` and `sort_direction` query parameters. Which fields are sortable is documented in the respective operation's description.  Note that if no sortable fields are documented for an operation, sorting is not supported *at all*.  ## Field expansion  Some collection endpoints support an `expand` query parameter. Passing an expand value includes optional fields in each response item that are omitted by default, typically because they are expensive to compute and only needed in specific contexts.  Valid `expand` values for an endpoint are listed in its operation description. Unknown values are silently ignored.  ## Errors  All error responses use the `application/problem+json` media type as defined in [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html).  Example:  ```json {   \"type\": \"about:blank\",   \"status\": 404,   \"title\": \"Resource Not Found\",   \"detail\": \"No project with UUID 0976ef46-e4a0-4be4-8b0e-00e9d3625c1f exists\" } ```  ## Authentication  Two credential types are accepted:  - **API key** via the `X-Api-Key` request header. API keys are   long-lived and intended for machine-to-machine integrations. - **Bearer token** via the `Authorization: Bearer <token>` request   header. Bearer tokens are short-lived, user-bound, and opaque   server-issued session tokens.  Bearer tokens are obtained from one of the following endpoints:  - `POST /api/v1/user/login` - `POST /api/v1/user/oidc/login`  API v2 does not yet provide its own login endpoints; clients use the v1 endpoints above to acquire tokens and then call v2 with them.  Tokens are valid for 8 hours by default and **cannot be refreshed**. Clients must re-authenticate once a token expires.  Requests with missing or invalid credentials are rejected with `401 Unauthorized`.  ## Authorization  Access is gated by named permissions. Operations document the permission(s) they require; operations without a documented permission requirement only require authentication.  When the *Portfolio Access Control* feature is enabled (disabled by default), project-scoped operations additionally enforce per-project access via team membership. The `PORTFOLIO_ACCESS_CONTROL_BYPASS` permission grants access to all projects regardless of team mappings. When the feature is disabled, all authenticated callers holding the required permission can access all projects.  Authenticated requests that lack the required permission, or that target a project the caller cannot access, are rejected with `403 Forbidden`.  ## HTTP Methods  | Method   | Semantics                  | |----------|----------------------------| | `GET`    | Retrieve a resource        | | `POST`   | Create a new resource      | | `PUT`    | Update a resource          | | `PATCH`  | Partially update a resource| | `DELETE` | Delete a resource          |  ## Response Conventions  Create and update operations (`POST`, `PUT`, `PATCH`) do not return the full resource in the response. They return either no body, or only server-generated identifiers (e.g. a UUID). `POST` responses may include a `Location` header linking to the created resource.  Delete operations return `204 No Content` with no body.  ## Deprecations  Operations may be removed or replaced over time. When a response carries the `X-API-Deprecated: true` header, the operation that produced it is deprecated and may be removed in a future release. Clients should check for this header on every response and surface it (e.g. via a log warning) so that operators are aware of upcoming breakages. The respective operation's description points out which alternative operation(s) to use.  ## Internal operations  Operations under the `/internal` path prefix expose system internals and are reserved for first-party use. They are **not** part of the stable v2 API contract and may change or be removed without notice. Third-party clients should not depend on them.

API version: 2.0.0
Contact: dependencytrack@owasp.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dtv2

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the JsonSchemaValidationError type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JsonSchemaValidationError{}

// JsonSchemaValidationError A JSON Schema validation error as per <https://json-schema.org/draft/2020-12/json-schema-core.html#name-output-structure>.
type JsonSchemaValidationError struct {
	// JSON Pointer to the location in the instance that failed validation
	InstanceLocation string `json:"instance_location"`
	// JSON Pointer to the location in the schema during evaluation
	EvaluationPath *string `json:"evaluation_path,omitempty"`
	// Schema location that generated the error
	SchemaLocation *string `json:"schema_location,omitempty"`
	// The validation keyword that failed
	Keyword *string `json:"keyword,omitempty"`
	// Human-readable error message
	Message string `json:"message"`
}

type _JsonSchemaValidationError JsonSchemaValidationError

// NewJsonSchemaValidationError instantiates a new JsonSchemaValidationError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJsonSchemaValidationError(instanceLocation string, message string) *JsonSchemaValidationError {
	this := JsonSchemaValidationError{}
	this.InstanceLocation = instanceLocation
	this.Message = message
	return &this
}

// NewJsonSchemaValidationErrorWithDefaults instantiates a new JsonSchemaValidationError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJsonSchemaValidationErrorWithDefaults() *JsonSchemaValidationError {
	this := JsonSchemaValidationError{}
	return &this
}

// GetInstanceLocation returns the InstanceLocation field value
func (o *JsonSchemaValidationError) GetInstanceLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.InstanceLocation
}

// GetInstanceLocationOk returns a tuple with the InstanceLocation field value
// and a boolean to check if the value has been set.
func (o *JsonSchemaValidationError) GetInstanceLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InstanceLocation, true
}

// SetInstanceLocation sets field value
func (o *JsonSchemaValidationError) SetInstanceLocation(v string) {
	o.InstanceLocation = v
}

// GetEvaluationPath returns the EvaluationPath field value if set, zero value otherwise.
func (o *JsonSchemaValidationError) GetEvaluationPath() string {
	if o == nil || IsNil(o.EvaluationPath) {
		var ret string
		return ret
	}
	return *o.EvaluationPath
}

// GetEvaluationPathOk returns a tuple with the EvaluationPath field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JsonSchemaValidationError) GetEvaluationPathOk() (*string, bool) {
	if o == nil || IsNil(o.EvaluationPath) {
		return nil, false
	}
	return o.EvaluationPath, true
}

// HasEvaluationPath returns a boolean if a field has been set.
func (o *JsonSchemaValidationError) HasEvaluationPath() bool {
	if o != nil && !IsNil(o.EvaluationPath) {
		return true
	}

	return false
}

// SetEvaluationPath gets a reference to the given string and assigns it to the EvaluationPath field.
func (o *JsonSchemaValidationError) SetEvaluationPath(v string) {
	o.EvaluationPath = &v
}

// GetSchemaLocation returns the SchemaLocation field value if set, zero value otherwise.
func (o *JsonSchemaValidationError) GetSchemaLocation() string {
	if o == nil || IsNil(o.SchemaLocation) {
		var ret string
		return ret
	}
	return *o.SchemaLocation
}

// GetSchemaLocationOk returns a tuple with the SchemaLocation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JsonSchemaValidationError) GetSchemaLocationOk() (*string, bool) {
	if o == nil || IsNil(o.SchemaLocation) {
		return nil, false
	}
	return o.SchemaLocation, true
}

// HasSchemaLocation returns a boolean if a field has been set.
func (o *JsonSchemaValidationError) HasSchemaLocation() bool {
	if o != nil && !IsNil(o.SchemaLocation) {
		return true
	}

	return false
}

// SetSchemaLocation gets a reference to the given string and assigns it to the SchemaLocation field.
func (o *JsonSchemaValidationError) SetSchemaLocation(v string) {
	o.SchemaLocation = &v
}

// GetKeyword returns the Keyword field value if set, zero value otherwise.
func (o *JsonSchemaValidationError) GetKeyword() string {
	if o == nil || IsNil(o.Keyword) {
		var ret string
		return ret
	}
	return *o.Keyword
}

// GetKeywordOk returns a tuple with the Keyword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JsonSchemaValidationError) GetKeywordOk() (*string, bool) {
	if o == nil || IsNil(o.Keyword) {
		return nil, false
	}
	return o.Keyword, true
}

// HasKeyword returns a boolean if a field has been set.
func (o *JsonSchemaValidationError) HasKeyword() bool {
	if o != nil && !IsNil(o.Keyword) {
		return true
	}

	return false
}

// SetKeyword gets a reference to the given string and assigns it to the Keyword field.
func (o *JsonSchemaValidationError) SetKeyword(v string) {
	o.Keyword = &v
}

// GetMessage returns the Message field value
func (o *JsonSchemaValidationError) GetMessage() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Message
}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
func (o *JsonSchemaValidationError) GetMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Message, true
}

// SetMessage sets field value
func (o *JsonSchemaValidationError) SetMessage(v string) {
	o.Message = v
}

func (o JsonSchemaValidationError) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JsonSchemaValidationError) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["instance_location"] = o.InstanceLocation
	if !IsNil(o.EvaluationPath) {
		toSerialize["evaluation_path"] = o.EvaluationPath
	}
	if !IsNil(o.SchemaLocation) {
		toSerialize["schema_location"] = o.SchemaLocation
	}
	if !IsNil(o.Keyword) {
		toSerialize["keyword"] = o.Keyword
	}
	toSerialize["message"] = o.Message
	return toSerialize, nil
}

func (o *JsonSchemaValidationError) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"instance_location",
		"message",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varJsonSchemaValidationError := _JsonSchemaValidationError{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varJsonSchemaValidationError)

	if err != nil {
		return err
	}

	*o = JsonSchemaValidationError(varJsonSchemaValidationError)

	return err
}

type NullableJsonSchemaValidationError struct {
	value *JsonSchemaValidationError
	isSet bool
}

func (v NullableJsonSchemaValidationError) Get() *JsonSchemaValidationError {
	return v.value
}

func (v *NullableJsonSchemaValidationError) Set(val *JsonSchemaValidationError) {
	v.value = val
	v.isSet = true
}

func (v NullableJsonSchemaValidationError) IsSet() bool {
	return v.isSet
}

func (v *NullableJsonSchemaValidationError) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJsonSchemaValidationError(val *JsonSchemaValidationError) *NullableJsonSchemaValidationError {
	return &NullableJsonSchemaValidationError{value: val, isSet: true}
}

func (v NullableJsonSchemaValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJsonSchemaValidationError) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
