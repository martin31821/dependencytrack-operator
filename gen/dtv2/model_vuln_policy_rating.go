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

// checks if the VulnPolicyRating type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VulnPolicyRating{}

// VulnPolicyRating struct for VulnPolicyRating
type VulnPolicyRating struct {
	Method   string   `json:"method"`
	Severity string   `json:"severity"`
	Vector   *string  `json:"vector,omitempty"`
	Score    *float64 `json:"score,omitempty"`
}

type _VulnPolicyRating VulnPolicyRating

// NewVulnPolicyRating instantiates a new VulnPolicyRating object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVulnPolicyRating(method string, severity string) *VulnPolicyRating {
	this := VulnPolicyRating{}
	this.Method = method
	this.Severity = severity
	return &this
}

// NewVulnPolicyRatingWithDefaults instantiates a new VulnPolicyRating object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVulnPolicyRatingWithDefaults() *VulnPolicyRating {
	this := VulnPolicyRating{}
	return &this
}

// GetMethod returns the Method field value
func (o *VulnPolicyRating) GetMethod() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Method
}

// GetMethodOk returns a tuple with the Method field value
// and a boolean to check if the value has been set.
func (o *VulnPolicyRating) GetMethodOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Method, true
}

// SetMethod sets field value
func (o *VulnPolicyRating) SetMethod(v string) {
	o.Method = v
}

// GetSeverity returns the Severity field value
func (o *VulnPolicyRating) GetSeverity() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Severity
}

// GetSeverityOk returns a tuple with the Severity field value
// and a boolean to check if the value has been set.
func (o *VulnPolicyRating) GetSeverityOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Severity, true
}

// SetSeverity sets field value
func (o *VulnPolicyRating) SetSeverity(v string) {
	o.Severity = v
}

// GetVector returns the Vector field value if set, zero value otherwise.
func (o *VulnPolicyRating) GetVector() string {
	if o == nil || IsNil(o.Vector) {
		var ret string
		return ret
	}
	return *o.Vector
}

// GetVectorOk returns a tuple with the Vector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VulnPolicyRating) GetVectorOk() (*string, bool) {
	if o == nil || IsNil(o.Vector) {
		return nil, false
	}
	return o.Vector, true
}

// HasVector returns a boolean if a field has been set.
func (o *VulnPolicyRating) HasVector() bool {
	if o != nil && !IsNil(o.Vector) {
		return true
	}

	return false
}

// SetVector gets a reference to the given string and assigns it to the Vector field.
func (o *VulnPolicyRating) SetVector(v string) {
	o.Vector = &v
}

// GetScore returns the Score field value if set, zero value otherwise.
func (o *VulnPolicyRating) GetScore() float64 {
	if o == nil || IsNil(o.Score) {
		var ret float64
		return ret
	}
	return *o.Score
}

// GetScoreOk returns a tuple with the Score field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VulnPolicyRating) GetScoreOk() (*float64, bool) {
	if o == nil || IsNil(o.Score) {
		return nil, false
	}
	return o.Score, true
}

// HasScore returns a boolean if a field has been set.
func (o *VulnPolicyRating) HasScore() bool {
	if o != nil && !IsNil(o.Score) {
		return true
	}

	return false
}

// SetScore gets a reference to the given float64 and assigns it to the Score field.
func (o *VulnPolicyRating) SetScore(v float64) {
	o.Score = &v
}

func (o VulnPolicyRating) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VulnPolicyRating) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["method"] = o.Method
	toSerialize["severity"] = o.Severity
	if !IsNil(o.Vector) {
		toSerialize["vector"] = o.Vector
	}
	if !IsNil(o.Score) {
		toSerialize["score"] = o.Score
	}
	return toSerialize, nil
}

func (o *VulnPolicyRating) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"method",
		"severity",
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

	varVulnPolicyRating := _VulnPolicyRating{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVulnPolicyRating)

	if err != nil {
		return err
	}

	*o = VulnPolicyRating(varVulnPolicyRating)

	return err
}

type NullableVulnPolicyRating struct {
	value *VulnPolicyRating
	isSet bool
}

func (v NullableVulnPolicyRating) Get() *VulnPolicyRating {
	return v.value
}

func (v *NullableVulnPolicyRating) Set(val *VulnPolicyRating) {
	v.value = val
	v.isSet = true
}

func (v NullableVulnPolicyRating) IsSet() bool {
	return v.isSet
}

func (v *NullableVulnPolicyRating) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVulnPolicyRating(val *VulnPolicyRating) *NullableVulnPolicyRating {
	return &NullableVulnPolicyRating{value: val, isSet: true}
}

func (v NullableVulnPolicyRating) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVulnPolicyRating) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
