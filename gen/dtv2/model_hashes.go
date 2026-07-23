/*
OWASP Dependency-Track

# REST API of OWASP Dependency-Track  ## Pagination  This API implements token-based pagination. Collection responses have the following structure:  ```json {   \"items\": [...],   \"next_page_token\": \"abcdefg\",   \"total\": {     \"count\": 100,     \"type\": \"EXACT\"   } } ```  `next_page_token` is present when more items exist, and absent otherwise. To fetch the next page, pass it as the `page_token` query parameter.  To navigate backwards, clients should keep track of previous page tokens as they paginate through collections. The API does *not* provide backward navigation!  Collections that support sorting will only consider the `sort_by` and `sort_direction` query parameters for the request of the first page. For subsequent pages, sorting preferences are bound to the page token.  Page tokens are opaque strings. Clients should not try to interpret or generate them. Their format may change without notice.  The `total` object discloses how many items exist in the collection *across all pages*. Because counting is expensive, some collections that hold *a lot* of items may return partial counts (type `AT_LEAST`) instead of exact counts (type `EXACT`). Which type to expect is usually documented in the operation's description.  ## Sorting  Items in a collection can be sorted using the `sort_by` and `sort_direction` query parameters. Which fields are sortable is documented in the respective operation's description.  Note that if no sortable fields are documented for an operation, sorting is not supported *at all*.  ## Field expansion  Some collection endpoints support an `expand` query parameter. Passing an expand value includes optional fields in each response item that are omitted by default, typically because they are expensive to compute and only needed in specific contexts.  Valid `expand` values for an endpoint are listed in its operation description. Unknown values are silently ignored.  ## Errors  All error responses use the `application/problem+json` media type as defined in [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html).  Example:  ```json {   \"type\": \"about:blank\",   \"status\": 404,   \"title\": \"Resource Not Found\",   \"detail\": \"No project with UUID 0976ef46-e4a0-4be4-8b0e-00e9d3625c1f exists\" } ```  ## Authentication  Two credential types are accepted:  - **API key** via the `X-Api-Key` request header. API keys are   long-lived and intended for machine-to-machine integrations. - **Bearer token** via the `Authorization: Bearer <token>` request   header. Bearer tokens are short-lived, user-bound, and opaque   server-issued session tokens.  Bearer tokens are obtained from one of the following endpoints:  - `POST /api/v1/user/login` - `POST /api/v1/user/oidc/login`  API v2 does not yet provide its own login endpoints; clients use the v1 endpoints above to acquire tokens and then call v2 with them.  Tokens are valid for 8 hours by default and **cannot be refreshed**. Clients must re-authenticate once a token expires.  Requests with missing or invalid credentials are rejected with `401 Unauthorized`.  ## Authorization  Access is gated by named permissions. Operations document the permission(s) they require; operations without a documented permission requirement only require authentication.  When the *Portfolio Access Control* feature is enabled (disabled by default), project-scoped operations additionally enforce per-project access via team membership. The `PORTFOLIO_ACCESS_CONTROL_BYPASS` permission grants access to all projects regardless of team mappings. When the feature is disabled, all authenticated callers holding the required permission can access all projects.  Authenticated requests that lack the required permission, or that target a project the caller cannot access, are rejected with `403 Forbidden`.  ## HTTP Methods  | Method   | Semantics                  | |----------|----------------------------| | `GET`    | Retrieve a resource        | | `POST`   | Create a new resource      | | `PUT`    | Update a resource          | | `PATCH`  | Partially update a resource| | `DELETE` | Delete a resource          |  ## Response Conventions  Create and update operations (`POST`, `PUT`, `PATCH`) do not return the full resource in the response. They return either no body, or only server-generated identifiers (e.g. a UUID). `POST` responses may include a `Location` header linking to the created resource.  Delete operations return `204 No Content` with no body.  ## Deprecations  Operations may be removed or replaced over time. When a response carries the `X-API-Deprecated: true` header, the operation that produced it is deprecated and may be removed in a future release. Clients should check for this header on every response and surface it (e.g. via a log warning) so that operators are aware of upcoming breakages. The respective operation's description points out which alternative operation(s) to use.  ## Internal operations  Operations under the `/internal` path prefix expose system internals and are reserved for first-party use. They are **not** part of the stable v2 API contract and may change or be removed without notice. Third-party clients should not depend on them.

API version: 2.0.0
Contact: dependencytrack@owasp.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dtv2

import (
	"encoding/json"
)

// checks if the Hashes type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Hashes{}

// Hashes struct for Hashes
type Hashes struct {
	Sha1       *string `json:"sha1,omitempty"`
	Sha256     *string `json:"sha256,omitempty"`
	Sha384     *string `json:"sha384,omitempty"`
	Sha512     *string `json:"sha512,omitempty"`
	Sha3256    *string `json:"sha3_256,omitempty"`
	Sha3384    *string `json:"sha3_384,omitempty"`
	Sha3512    *string `json:"sha3_512,omitempty"`
	Blake2b256 *string `json:"blake2b_256,omitempty"`
	Blake2b384 *string `json:"blake2b_384,omitempty"`
	Blake2b512 *string `json:"blake2b_512,omitempty"`
	Blake3     *string `json:"blake3,omitempty"`
	Md5        *string `json:"md5,omitempty"`
}

// NewHashes instantiates a new Hashes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHashes() *Hashes {
	this := Hashes{}
	return &this
}

// NewHashesWithDefaults instantiates a new Hashes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHashesWithDefaults() *Hashes {
	this := Hashes{}
	return &this
}

// GetSha1 returns the Sha1 field value if set, zero value otherwise.
func (o *Hashes) GetSha1() string {
	if o == nil || IsNil(o.Sha1) {
		var ret string
		return ret
	}
	return *o.Sha1
}

// GetSha1Ok returns a tuple with the Sha1 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha1Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha1) {
		return nil, false
	}
	return o.Sha1, true
}

// HasSha1 returns a boolean if a field has been set.
func (o *Hashes) HasSha1() bool {
	if o != nil && !IsNil(o.Sha1) {
		return true
	}

	return false
}

// SetSha1 gets a reference to the given string and assigns it to the Sha1 field.
func (o *Hashes) SetSha1(v string) {
	o.Sha1 = &v
}

// GetSha256 returns the Sha256 field value if set, zero value otherwise.
func (o *Hashes) GetSha256() string {
	if o == nil || IsNil(o.Sha256) {
		var ret string
		return ret
	}
	return *o.Sha256
}

// GetSha256Ok returns a tuple with the Sha256 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha256Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha256) {
		return nil, false
	}
	return o.Sha256, true
}

// HasSha256 returns a boolean if a field has been set.
func (o *Hashes) HasSha256() bool {
	if o != nil && !IsNil(o.Sha256) {
		return true
	}

	return false
}

// SetSha256 gets a reference to the given string and assigns it to the Sha256 field.
func (o *Hashes) SetSha256(v string) {
	o.Sha256 = &v
}

// GetSha384 returns the Sha384 field value if set, zero value otherwise.
func (o *Hashes) GetSha384() string {
	if o == nil || IsNil(o.Sha384) {
		var ret string
		return ret
	}
	return *o.Sha384
}

// GetSha384Ok returns a tuple with the Sha384 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha384Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha384) {
		return nil, false
	}
	return o.Sha384, true
}

// HasSha384 returns a boolean if a field has been set.
func (o *Hashes) HasSha384() bool {
	if o != nil && !IsNil(o.Sha384) {
		return true
	}

	return false
}

// SetSha384 gets a reference to the given string and assigns it to the Sha384 field.
func (o *Hashes) SetSha384(v string) {
	o.Sha384 = &v
}

// GetSha512 returns the Sha512 field value if set, zero value otherwise.
func (o *Hashes) GetSha512() string {
	if o == nil || IsNil(o.Sha512) {
		var ret string
		return ret
	}
	return *o.Sha512
}

// GetSha512Ok returns a tuple with the Sha512 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha512Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha512) {
		return nil, false
	}
	return o.Sha512, true
}

// HasSha512 returns a boolean if a field has been set.
func (o *Hashes) HasSha512() bool {
	if o != nil && !IsNil(o.Sha512) {
		return true
	}

	return false
}

// SetSha512 gets a reference to the given string and assigns it to the Sha512 field.
func (o *Hashes) SetSha512(v string) {
	o.Sha512 = &v
}

// GetSha3256 returns the Sha3256 field value if set, zero value otherwise.
func (o *Hashes) GetSha3256() string {
	if o == nil || IsNil(o.Sha3256) {
		var ret string
		return ret
	}
	return *o.Sha3256
}

// GetSha3256Ok returns a tuple with the Sha3256 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha3256Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha3256) {
		return nil, false
	}
	return o.Sha3256, true
}

// HasSha3256 returns a boolean if a field has been set.
func (o *Hashes) HasSha3256() bool {
	if o != nil && !IsNil(o.Sha3256) {
		return true
	}

	return false
}

// SetSha3256 gets a reference to the given string and assigns it to the Sha3256 field.
func (o *Hashes) SetSha3256(v string) {
	o.Sha3256 = &v
}

// GetSha3384 returns the Sha3384 field value if set, zero value otherwise.
func (o *Hashes) GetSha3384() string {
	if o == nil || IsNil(o.Sha3384) {
		var ret string
		return ret
	}
	return *o.Sha3384
}

// GetSha3384Ok returns a tuple with the Sha3384 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha3384Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha3384) {
		return nil, false
	}
	return o.Sha3384, true
}

// HasSha3384 returns a boolean if a field has been set.
func (o *Hashes) HasSha3384() bool {
	if o != nil && !IsNil(o.Sha3384) {
		return true
	}

	return false
}

// SetSha3384 gets a reference to the given string and assigns it to the Sha3384 field.
func (o *Hashes) SetSha3384(v string) {
	o.Sha3384 = &v
}

// GetSha3512 returns the Sha3512 field value if set, zero value otherwise.
func (o *Hashes) GetSha3512() string {
	if o == nil || IsNil(o.Sha3512) {
		var ret string
		return ret
	}
	return *o.Sha3512
}

// GetSha3512Ok returns a tuple with the Sha3512 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetSha3512Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha3512) {
		return nil, false
	}
	return o.Sha3512, true
}

// HasSha3512 returns a boolean if a field has been set.
func (o *Hashes) HasSha3512() bool {
	if o != nil && !IsNil(o.Sha3512) {
		return true
	}

	return false
}

// SetSha3512 gets a reference to the given string and assigns it to the Sha3512 field.
func (o *Hashes) SetSha3512(v string) {
	o.Sha3512 = &v
}

// GetBlake2b256 returns the Blake2b256 field value if set, zero value otherwise.
func (o *Hashes) GetBlake2b256() string {
	if o == nil || IsNil(o.Blake2b256) {
		var ret string
		return ret
	}
	return *o.Blake2b256
}

// GetBlake2b256Ok returns a tuple with the Blake2b256 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetBlake2b256Ok() (*string, bool) {
	if o == nil || IsNil(o.Blake2b256) {
		return nil, false
	}
	return o.Blake2b256, true
}

// HasBlake2b256 returns a boolean if a field has been set.
func (o *Hashes) HasBlake2b256() bool {
	if o != nil && !IsNil(o.Blake2b256) {
		return true
	}

	return false
}

// SetBlake2b256 gets a reference to the given string and assigns it to the Blake2b256 field.
func (o *Hashes) SetBlake2b256(v string) {
	o.Blake2b256 = &v
}

// GetBlake2b384 returns the Blake2b384 field value if set, zero value otherwise.
func (o *Hashes) GetBlake2b384() string {
	if o == nil || IsNil(o.Blake2b384) {
		var ret string
		return ret
	}
	return *o.Blake2b384
}

// GetBlake2b384Ok returns a tuple with the Blake2b384 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetBlake2b384Ok() (*string, bool) {
	if o == nil || IsNil(o.Blake2b384) {
		return nil, false
	}
	return o.Blake2b384, true
}

// HasBlake2b384 returns a boolean if a field has been set.
func (o *Hashes) HasBlake2b384() bool {
	if o != nil && !IsNil(o.Blake2b384) {
		return true
	}

	return false
}

// SetBlake2b384 gets a reference to the given string and assigns it to the Blake2b384 field.
func (o *Hashes) SetBlake2b384(v string) {
	o.Blake2b384 = &v
}

// GetBlake2b512 returns the Blake2b512 field value if set, zero value otherwise.
func (o *Hashes) GetBlake2b512() string {
	if o == nil || IsNil(o.Blake2b512) {
		var ret string
		return ret
	}
	return *o.Blake2b512
}

// GetBlake2b512Ok returns a tuple with the Blake2b512 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetBlake2b512Ok() (*string, bool) {
	if o == nil || IsNil(o.Blake2b512) {
		return nil, false
	}
	return o.Blake2b512, true
}

// HasBlake2b512 returns a boolean if a field has been set.
func (o *Hashes) HasBlake2b512() bool {
	if o != nil && !IsNil(o.Blake2b512) {
		return true
	}

	return false
}

// SetBlake2b512 gets a reference to the given string and assigns it to the Blake2b512 field.
func (o *Hashes) SetBlake2b512(v string) {
	o.Blake2b512 = &v
}

// GetBlake3 returns the Blake3 field value if set, zero value otherwise.
func (o *Hashes) GetBlake3() string {
	if o == nil || IsNil(o.Blake3) {
		var ret string
		return ret
	}
	return *o.Blake3
}

// GetBlake3Ok returns a tuple with the Blake3 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetBlake3Ok() (*string, bool) {
	if o == nil || IsNil(o.Blake3) {
		return nil, false
	}
	return o.Blake3, true
}

// HasBlake3 returns a boolean if a field has been set.
func (o *Hashes) HasBlake3() bool {
	if o != nil && !IsNil(o.Blake3) {
		return true
	}

	return false
}

// SetBlake3 gets a reference to the given string and assigns it to the Blake3 field.
func (o *Hashes) SetBlake3(v string) {
	o.Blake3 = &v
}

// GetMd5 returns the Md5 field value if set, zero value otherwise.
func (o *Hashes) GetMd5() string {
	if o == nil || IsNil(o.Md5) {
		var ret string
		return ret
	}
	return *o.Md5
}

// GetMd5Ok returns a tuple with the Md5 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Hashes) GetMd5Ok() (*string, bool) {
	if o == nil || IsNil(o.Md5) {
		return nil, false
	}
	return o.Md5, true
}

// HasMd5 returns a boolean if a field has been set.
func (o *Hashes) HasMd5() bool {
	if o != nil && !IsNil(o.Md5) {
		return true
	}

	return false
}

// SetMd5 gets a reference to the given string and assigns it to the Md5 field.
func (o *Hashes) SetMd5(v string) {
	o.Md5 = &v
}

func (o Hashes) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Hashes) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Sha1) {
		toSerialize["sha1"] = o.Sha1
	}
	if !IsNil(o.Sha256) {
		toSerialize["sha256"] = o.Sha256
	}
	if !IsNil(o.Sha384) {
		toSerialize["sha384"] = o.Sha384
	}
	if !IsNil(o.Sha512) {
		toSerialize["sha512"] = o.Sha512
	}
	if !IsNil(o.Sha3256) {
		toSerialize["sha3_256"] = o.Sha3256
	}
	if !IsNil(o.Sha3384) {
		toSerialize["sha3_384"] = o.Sha3384
	}
	if !IsNil(o.Sha3512) {
		toSerialize["sha3_512"] = o.Sha3512
	}
	if !IsNil(o.Blake2b256) {
		toSerialize["blake2b_256"] = o.Blake2b256
	}
	if !IsNil(o.Blake2b384) {
		toSerialize["blake2b_384"] = o.Blake2b384
	}
	if !IsNil(o.Blake2b512) {
		toSerialize["blake2b_512"] = o.Blake2b512
	}
	if !IsNil(o.Blake3) {
		toSerialize["blake3"] = o.Blake3
	}
	if !IsNil(o.Md5) {
		toSerialize["md5"] = o.Md5
	}
	return toSerialize, nil
}

type NullableHashes struct {
	value *Hashes
	isSet bool
}

func (v NullableHashes) Get() *Hashes {
	return v.value
}

func (v *NullableHashes) Set(val *Hashes) {
	v.value = val
	v.isSet = true
}

func (v NullableHashes) IsSet() bool {
	return v.isSet
}

func (v *NullableHashes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHashes(val *Hashes) *NullableHashes {
	return &NullableHashes{value: val, isSet: true}
}

func (v NullableHashes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHashes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
