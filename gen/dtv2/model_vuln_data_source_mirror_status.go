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

// checks if the VulnDataSourceMirrorStatus type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VulnDataSourceMirrorStatus{}

// VulnDataSourceMirrorStatus struct for VulnDataSourceMirrorStatus
type VulnDataSourceMirrorStatus struct {
	// Status of the mirror run.
	Status string `json:"status"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	StartedAt *int64 `json:"started_at,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	CompletedAt *int64 `json:"completed_at,omitempty"`
	// Reason for why the mirror run failed.
	FailureReason *string `json:"failure_reason,omitempty"`
}

type _VulnDataSourceMirrorStatus VulnDataSourceMirrorStatus

// NewVulnDataSourceMirrorStatus instantiates a new VulnDataSourceMirrorStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVulnDataSourceMirrorStatus(status string) *VulnDataSourceMirrorStatus {
	this := VulnDataSourceMirrorStatus{}
	this.Status = status
	return &this
}

// NewVulnDataSourceMirrorStatusWithDefaults instantiates a new VulnDataSourceMirrorStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVulnDataSourceMirrorStatusWithDefaults() *VulnDataSourceMirrorStatus {
	this := VulnDataSourceMirrorStatus{}
	return &this
}

// GetStatus returns the Status field value
func (o *VulnDataSourceMirrorStatus) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *VulnDataSourceMirrorStatus) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *VulnDataSourceMirrorStatus) SetStatus(v string) {
	o.Status = v
}

// GetStartedAt returns the StartedAt field value if set, zero value otherwise.
func (o *VulnDataSourceMirrorStatus) GetStartedAt() int64 {
	if o == nil || IsNil(o.StartedAt) {
		var ret int64
		return ret
	}
	return *o.StartedAt
}

// GetStartedAtOk returns a tuple with the StartedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VulnDataSourceMirrorStatus) GetStartedAtOk() (*int64, bool) {
	if o == nil || IsNil(o.StartedAt) {
		return nil, false
	}
	return o.StartedAt, true
}

// HasStartedAt returns a boolean if a field has been set.
func (o *VulnDataSourceMirrorStatus) HasStartedAt() bool {
	if o != nil && !IsNil(o.StartedAt) {
		return true
	}

	return false
}

// SetStartedAt gets a reference to the given int64 and assigns it to the StartedAt field.
func (o *VulnDataSourceMirrorStatus) SetStartedAt(v int64) {
	o.StartedAt = &v
}

// GetCompletedAt returns the CompletedAt field value if set, zero value otherwise.
func (o *VulnDataSourceMirrorStatus) GetCompletedAt() int64 {
	if o == nil || IsNil(o.CompletedAt) {
		var ret int64
		return ret
	}
	return *o.CompletedAt
}

// GetCompletedAtOk returns a tuple with the CompletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VulnDataSourceMirrorStatus) GetCompletedAtOk() (*int64, bool) {
	if o == nil || IsNil(o.CompletedAt) {
		return nil, false
	}
	return o.CompletedAt, true
}

// HasCompletedAt returns a boolean if a field has been set.
func (o *VulnDataSourceMirrorStatus) HasCompletedAt() bool {
	if o != nil && !IsNil(o.CompletedAt) {
		return true
	}

	return false
}

// SetCompletedAt gets a reference to the given int64 and assigns it to the CompletedAt field.
func (o *VulnDataSourceMirrorStatus) SetCompletedAt(v int64) {
	o.CompletedAt = &v
}

// GetFailureReason returns the FailureReason field value if set, zero value otherwise.
func (o *VulnDataSourceMirrorStatus) GetFailureReason() string {
	if o == nil || IsNil(o.FailureReason) {
		var ret string
		return ret
	}
	return *o.FailureReason
}

// GetFailureReasonOk returns a tuple with the FailureReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VulnDataSourceMirrorStatus) GetFailureReasonOk() (*string, bool) {
	if o == nil || IsNil(o.FailureReason) {
		return nil, false
	}
	return o.FailureReason, true
}

// HasFailureReason returns a boolean if a field has been set.
func (o *VulnDataSourceMirrorStatus) HasFailureReason() bool {
	if o != nil && !IsNil(o.FailureReason) {
		return true
	}

	return false
}

// SetFailureReason gets a reference to the given string and assigns it to the FailureReason field.
func (o *VulnDataSourceMirrorStatus) SetFailureReason(v string) {
	o.FailureReason = &v
}

func (o VulnDataSourceMirrorStatus) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VulnDataSourceMirrorStatus) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["status"] = o.Status
	if !IsNil(o.StartedAt) {
		toSerialize["started_at"] = o.StartedAt
	}
	if !IsNil(o.CompletedAt) {
		toSerialize["completed_at"] = o.CompletedAt
	}
	if !IsNil(o.FailureReason) {
		toSerialize["failure_reason"] = o.FailureReason
	}
	return toSerialize, nil
}

func (o *VulnDataSourceMirrorStatus) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"status",
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

	varVulnDataSourceMirrorStatus := _VulnDataSourceMirrorStatus{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVulnDataSourceMirrorStatus)

	if err != nil {
		return err
	}

	*o = VulnDataSourceMirrorStatus(varVulnDataSourceMirrorStatus)

	return err
}

type NullableVulnDataSourceMirrorStatus struct {
	value *VulnDataSourceMirrorStatus
	isSet bool
}

func (v NullableVulnDataSourceMirrorStatus) Get() *VulnDataSourceMirrorStatus {
	return v.value
}

func (v *NullableVulnDataSourceMirrorStatus) Set(val *VulnDataSourceMirrorStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableVulnDataSourceMirrorStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableVulnDataSourceMirrorStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVulnDataSourceMirrorStatus(val *VulnDataSourceMirrorStatus) *NullableVulnDataSourceMirrorStatus {
	return &NullableVulnDataSourceMirrorStatus{value: val, isSet: true}
}

func (v NullableVulnDataSourceMirrorStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVulnDataSourceMirrorStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
