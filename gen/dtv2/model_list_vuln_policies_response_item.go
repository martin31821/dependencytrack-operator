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

// checks if the ListVulnPoliciesResponseItem type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListVulnPoliciesResponseItem{}

// ListVulnPoliciesResponseItem struct for ListVulnPoliciesResponseItem
type ListVulnPoliciesResponseItem struct {
	Uuid          string                  `json:"uuid"`
	Name          string                  `json:"name"`
	Description   *string                 `json:"description,omitempty"`
	Author        *string                 `json:"author,omitempty"`
	Priority      int32                   `json:"priority"`
	OperationMode VulnPolicyOperationMode `json:"operation_mode"`
	Source        VulnPolicySource        `json:"source"`
}

type _ListVulnPoliciesResponseItem ListVulnPoliciesResponseItem

// NewListVulnPoliciesResponseItem instantiates a new ListVulnPoliciesResponseItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListVulnPoliciesResponseItem(uuid string, name string, priority int32, operationMode VulnPolicyOperationMode, source VulnPolicySource) *ListVulnPoliciesResponseItem {
	this := ListVulnPoliciesResponseItem{}
	this.Uuid = uuid
	this.Name = name
	this.Priority = priority
	this.OperationMode = operationMode
	this.Source = source
	return &this
}

// NewListVulnPoliciesResponseItemWithDefaults instantiates a new ListVulnPoliciesResponseItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListVulnPoliciesResponseItemWithDefaults() *ListVulnPoliciesResponseItem {
	this := ListVulnPoliciesResponseItem{}
	return &this
}

// GetUuid returns the Uuid field value
func (o *ListVulnPoliciesResponseItem) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *ListVulnPoliciesResponseItem) SetUuid(v string) {
	o.Uuid = v
}

// GetName returns the Name field value
func (o *ListVulnPoliciesResponseItem) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ListVulnPoliciesResponseItem) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *ListVulnPoliciesResponseItem) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ListVulnPoliciesResponseItem) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ListVulnPoliciesResponseItem) SetDescription(v string) {
	o.Description = &v
}

// GetAuthor returns the Author field value if set, zero value otherwise.
func (o *ListVulnPoliciesResponseItem) GetAuthor() string {
	if o == nil || IsNil(o.Author) {
		var ret string
		return ret
	}
	return *o.Author
}

// GetAuthorOk returns a tuple with the Author field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetAuthorOk() (*string, bool) {
	if o == nil || IsNil(o.Author) {
		return nil, false
	}
	return o.Author, true
}

// HasAuthor returns a boolean if a field has been set.
func (o *ListVulnPoliciesResponseItem) HasAuthor() bool {
	if o != nil && !IsNil(o.Author) {
		return true
	}

	return false
}

// SetAuthor gets a reference to the given string and assigns it to the Author field.
func (o *ListVulnPoliciesResponseItem) SetAuthor(v string) {
	o.Author = &v
}

// GetPriority returns the Priority field value
func (o *ListVulnPoliciesResponseItem) GetPriority() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetPriorityOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Priority, true
}

// SetPriority sets field value
func (o *ListVulnPoliciesResponseItem) SetPriority(v int32) {
	o.Priority = v
}

// GetOperationMode returns the OperationMode field value
func (o *ListVulnPoliciesResponseItem) GetOperationMode() VulnPolicyOperationMode {
	if o == nil {
		var ret VulnPolicyOperationMode
		return ret
	}

	return o.OperationMode
}

// GetOperationModeOk returns a tuple with the OperationMode field value
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetOperationModeOk() (*VulnPolicyOperationMode, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OperationMode, true
}

// SetOperationMode sets field value
func (o *ListVulnPoliciesResponseItem) SetOperationMode(v VulnPolicyOperationMode) {
	o.OperationMode = v
}

// GetSource returns the Source field value
func (o *ListVulnPoliciesResponseItem) GetSource() VulnPolicySource {
	if o == nil {
		var ret VulnPolicySource
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *ListVulnPoliciesResponseItem) GetSourceOk() (*VulnPolicySource, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *ListVulnPoliciesResponseItem) SetSource(v VulnPolicySource) {
	o.Source = v
}

func (o ListVulnPoliciesResponseItem) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListVulnPoliciesResponseItem) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uuid"] = o.Uuid
	toSerialize["name"] = o.Name
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Author) {
		toSerialize["author"] = o.Author
	}
	toSerialize["priority"] = o.Priority
	toSerialize["operation_mode"] = o.OperationMode
	toSerialize["source"] = o.Source
	return toSerialize, nil
}

func (o *ListVulnPoliciesResponseItem) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"uuid",
		"name",
		"priority",
		"operation_mode",
		"source",
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

	varListVulnPoliciesResponseItem := _ListVulnPoliciesResponseItem{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varListVulnPoliciesResponseItem)

	if err != nil {
		return err
	}

	*o = ListVulnPoliciesResponseItem(varListVulnPoliciesResponseItem)

	return err
}

type NullableListVulnPoliciesResponseItem struct {
	value *ListVulnPoliciesResponseItem
	isSet bool
}

func (v NullableListVulnPoliciesResponseItem) Get() *ListVulnPoliciesResponseItem {
	return v.value
}

func (v *NullableListVulnPoliciesResponseItem) Set(val *ListVulnPoliciesResponseItem) {
	v.value = val
	v.isSet = true
}

func (v NullableListVulnPoliciesResponseItem) IsSet() bool {
	return v.isSet
}

func (v *NullableListVulnPoliciesResponseItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListVulnPoliciesResponseItem(val *ListVulnPoliciesResponseItem) *NullableListVulnPoliciesResponseItem {
	return &NullableListVulnPoliciesResponseItem{value: val, isSet: true}
}

func (v NullableListVulnPoliciesResponseItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListVulnPoliciesResponseItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
