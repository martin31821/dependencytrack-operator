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

// checks if the ExtensionTestCheck type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExtensionTestCheck{}

// ExtensionTestCheck struct for ExtensionTestCheck
type ExtensionTestCheck struct {
	Name    string                   `json:"name"`
	Status  ExtensionTestCheckStatus `json:"status"`
	Message *string                  `json:"message,omitempty"`
}

type _ExtensionTestCheck ExtensionTestCheck

// NewExtensionTestCheck instantiates a new ExtensionTestCheck object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExtensionTestCheck(name string, status ExtensionTestCheckStatus) *ExtensionTestCheck {
	this := ExtensionTestCheck{}
	this.Name = name
	this.Status = status
	return &this
}

// NewExtensionTestCheckWithDefaults instantiates a new ExtensionTestCheck object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExtensionTestCheckWithDefaults() *ExtensionTestCheck {
	this := ExtensionTestCheck{}
	return &this
}

// GetName returns the Name field value
func (o *ExtensionTestCheck) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ExtensionTestCheck) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ExtensionTestCheck) SetName(v string) {
	o.Name = v
}

// GetStatus returns the Status field value
func (o *ExtensionTestCheck) GetStatus() ExtensionTestCheckStatus {
	if o == nil {
		var ret ExtensionTestCheckStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *ExtensionTestCheck) GetStatusOk() (*ExtensionTestCheckStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *ExtensionTestCheck) SetStatus(v ExtensionTestCheckStatus) {
	o.Status = v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ExtensionTestCheck) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtensionTestCheck) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ExtensionTestCheck) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ExtensionTestCheck) SetMessage(v string) {
	o.Message = &v
}

func (o ExtensionTestCheck) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExtensionTestCheck) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["status"] = o.Status
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	return toSerialize, nil
}

func (o *ExtensionTestCheck) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
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

	varExtensionTestCheck := _ExtensionTestCheck{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varExtensionTestCheck)

	if err != nil {
		return err
	}

	*o = ExtensionTestCheck(varExtensionTestCheck)

	return err
}

type NullableExtensionTestCheck struct {
	value *ExtensionTestCheck
	isSet bool
}

func (v NullableExtensionTestCheck) Get() *ExtensionTestCheck {
	return v.value
}

func (v *NullableExtensionTestCheck) Set(val *ExtensionTestCheck) {
	v.value = val
	v.isSet = true
}

func (v NullableExtensionTestCheck) IsSet() bool {
	return v.isSet
}

func (v *NullableExtensionTestCheck) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExtensionTestCheck(val *ExtensionTestCheck) *NullableExtensionTestCheck {
	return &NullableExtensionTestCheck{value: val, isSet: true}
}

func (v NullableExtensionTestCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExtensionTestCheck) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
