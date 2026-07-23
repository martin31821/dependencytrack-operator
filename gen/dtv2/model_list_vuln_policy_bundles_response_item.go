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

// checks if the ListVulnPolicyBundlesResponseItem type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListVulnPolicyBundlesResponseItem{}

// ListVulnPolicyBundlesResponseItem struct for ListVulnPolicyBundlesResponseItem
type ListVulnPolicyBundlesResponseItem struct {
	Uuid string  `json:"uuid"`
	Url  string  `json:"url"`
	Hash *string `json:"hash,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	LastSuccessfulSync *int64 `json:"last_successful_sync,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	Created *int64 `json:"created,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	Updated *int64 `json:"updated,omitempty"`
}

type _ListVulnPolicyBundlesResponseItem ListVulnPolicyBundlesResponseItem

// NewListVulnPolicyBundlesResponseItem instantiates a new ListVulnPolicyBundlesResponseItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListVulnPolicyBundlesResponseItem(uuid string, url string) *ListVulnPolicyBundlesResponseItem {
	this := ListVulnPolicyBundlesResponseItem{}
	this.Uuid = uuid
	this.Url = url
	return &this
}

// NewListVulnPolicyBundlesResponseItemWithDefaults instantiates a new ListVulnPolicyBundlesResponseItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListVulnPolicyBundlesResponseItemWithDefaults() *ListVulnPolicyBundlesResponseItem {
	this := ListVulnPolicyBundlesResponseItem{}
	return &this
}

// GetUuid returns the Uuid field value
func (o *ListVulnPolicyBundlesResponseItem) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *ListVulnPolicyBundlesResponseItem) SetUuid(v string) {
	o.Uuid = v
}

// GetUrl returns the Url field value
func (o *ListVulnPolicyBundlesResponseItem) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *ListVulnPolicyBundlesResponseItem) SetUrl(v string) {
	o.Url = v
}

// GetHash returns the Hash field value if set, zero value otherwise.
func (o *ListVulnPolicyBundlesResponseItem) GetHash() string {
	if o == nil || IsNil(o.Hash) {
		var ret string
		return ret
	}
	return *o.Hash
}

// GetHashOk returns a tuple with the Hash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetHashOk() (*string, bool) {
	if o == nil || IsNil(o.Hash) {
		return nil, false
	}
	return o.Hash, true
}

// HasHash returns a boolean if a field has been set.
func (o *ListVulnPolicyBundlesResponseItem) HasHash() bool {
	if o != nil && !IsNil(o.Hash) {
		return true
	}

	return false
}

// SetHash gets a reference to the given string and assigns it to the Hash field.
func (o *ListVulnPolicyBundlesResponseItem) SetHash(v string) {
	o.Hash = &v
}

// GetLastSuccessfulSync returns the LastSuccessfulSync field value if set, zero value otherwise.
func (o *ListVulnPolicyBundlesResponseItem) GetLastSuccessfulSync() int64 {
	if o == nil || IsNil(o.LastSuccessfulSync) {
		var ret int64
		return ret
	}
	return *o.LastSuccessfulSync
}

// GetLastSuccessfulSyncOk returns a tuple with the LastSuccessfulSync field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetLastSuccessfulSyncOk() (*int64, bool) {
	if o == nil || IsNil(o.LastSuccessfulSync) {
		return nil, false
	}
	return o.LastSuccessfulSync, true
}

// HasLastSuccessfulSync returns a boolean if a field has been set.
func (o *ListVulnPolicyBundlesResponseItem) HasLastSuccessfulSync() bool {
	if o != nil && !IsNil(o.LastSuccessfulSync) {
		return true
	}

	return false
}

// SetLastSuccessfulSync gets a reference to the given int64 and assigns it to the LastSuccessfulSync field.
func (o *ListVulnPolicyBundlesResponseItem) SetLastSuccessfulSync(v int64) {
	o.LastSuccessfulSync = &v
}

// GetCreated returns the Created field value if set, zero value otherwise.
func (o *ListVulnPolicyBundlesResponseItem) GetCreated() int64 {
	if o == nil || IsNil(o.Created) {
		var ret int64
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetCreatedOk() (*int64, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}
	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *ListVulnPolicyBundlesResponseItem) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given int64 and assigns it to the Created field.
func (o *ListVulnPolicyBundlesResponseItem) SetCreated(v int64) {
	o.Created = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise.
func (o *ListVulnPolicyBundlesResponseItem) GetUpdated() int64 {
	if o == nil || IsNil(o.Updated) {
		var ret int64
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListVulnPolicyBundlesResponseItem) GetUpdatedOk() (*int64, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}
	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *ListVulnPolicyBundlesResponseItem) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given int64 and assigns it to the Updated field.
func (o *ListVulnPolicyBundlesResponseItem) SetUpdated(v int64) {
	o.Updated = &v
}

func (o ListVulnPolicyBundlesResponseItem) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListVulnPolicyBundlesResponseItem) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uuid"] = o.Uuid
	toSerialize["url"] = o.Url
	if !IsNil(o.Hash) {
		toSerialize["hash"] = o.Hash
	}
	if !IsNil(o.LastSuccessfulSync) {
		toSerialize["last_successful_sync"] = o.LastSuccessfulSync
	}
	if !IsNil(o.Created) {
		toSerialize["created"] = o.Created
	}
	if !IsNil(o.Updated) {
		toSerialize["updated"] = o.Updated
	}
	return toSerialize, nil
}

func (o *ListVulnPolicyBundlesResponseItem) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"uuid",
		"url",
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

	varListVulnPolicyBundlesResponseItem := _ListVulnPolicyBundlesResponseItem{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varListVulnPolicyBundlesResponseItem)

	if err != nil {
		return err
	}

	*o = ListVulnPolicyBundlesResponseItem(varListVulnPolicyBundlesResponseItem)

	return err
}

type NullableListVulnPolicyBundlesResponseItem struct {
	value *ListVulnPolicyBundlesResponseItem
	isSet bool
}

func (v NullableListVulnPolicyBundlesResponseItem) Get() *ListVulnPolicyBundlesResponseItem {
	return v.value
}

func (v *NullableListVulnPolicyBundlesResponseItem) Set(val *ListVulnPolicyBundlesResponseItem) {
	v.value = val
	v.isSet = true
}

func (v NullableListVulnPolicyBundlesResponseItem) IsSet() bool {
	return v.isSet
}

func (v *NullableListVulnPolicyBundlesResponseItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListVulnPolicyBundlesResponseItem(val *ListVulnPolicyBundlesResponseItem) *NullableListVulnPolicyBundlesResponseItem {
	return &NullableListVulnPolicyBundlesResponseItem{value: val, isSet: true}
}

func (v NullableListVulnPolicyBundlesResponseItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListVulnPolicyBundlesResponseItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
