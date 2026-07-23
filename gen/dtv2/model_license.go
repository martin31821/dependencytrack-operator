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

// checks if the License type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &License{}

// License struct for License
type License struct {
	Name          *string `json:"name,omitempty"`
	LicenseId     *string `json:"license_id,omitempty"`
	Uuid          *string `json:"uuid,omitempty"`
	OsiApproved   *bool   `json:"osi_approved,omitempty"`
	FsfLibre      *bool   `json:"fsf_libre,omitempty"`
	CustomLicense *bool   `json:"custom_license,omitempty"`
}

// NewLicense instantiates a new License object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLicense() *License {
	this := License{}
	return &this
}

// NewLicenseWithDefaults instantiates a new License object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLicenseWithDefaults() *License {
	this := License{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *License) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *License) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *License) SetName(v string) {
	o.Name = &v
}

// GetLicenseId returns the LicenseId field value if set, zero value otherwise.
func (o *License) GetLicenseId() string {
	if o == nil || IsNil(o.LicenseId) {
		var ret string
		return ret
	}
	return *o.LicenseId
}

// GetLicenseIdOk returns a tuple with the LicenseId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetLicenseIdOk() (*string, bool) {
	if o == nil || IsNil(o.LicenseId) {
		return nil, false
	}
	return o.LicenseId, true
}

// HasLicenseId returns a boolean if a field has been set.
func (o *License) HasLicenseId() bool {
	if o != nil && !IsNil(o.LicenseId) {
		return true
	}

	return false
}

// SetLicenseId gets a reference to the given string and assigns it to the LicenseId field.
func (o *License) SetLicenseId(v string) {
	o.LicenseId = &v
}

// GetUuid returns the Uuid field value if set, zero value otherwise.
func (o *License) GetUuid() string {
	if o == nil || IsNil(o.Uuid) {
		var ret string
		return ret
	}
	return *o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetUuidOk() (*string, bool) {
	if o == nil || IsNil(o.Uuid) {
		return nil, false
	}
	return o.Uuid, true
}

// HasUuid returns a boolean if a field has been set.
func (o *License) HasUuid() bool {
	if o != nil && !IsNil(o.Uuid) {
		return true
	}

	return false
}

// SetUuid gets a reference to the given string and assigns it to the Uuid field.
func (o *License) SetUuid(v string) {
	o.Uuid = &v
}

// GetOsiApproved returns the OsiApproved field value if set, zero value otherwise.
func (o *License) GetOsiApproved() bool {
	if o == nil || IsNil(o.OsiApproved) {
		var ret bool
		return ret
	}
	return *o.OsiApproved
}

// GetOsiApprovedOk returns a tuple with the OsiApproved field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetOsiApprovedOk() (*bool, bool) {
	if o == nil || IsNil(o.OsiApproved) {
		return nil, false
	}
	return o.OsiApproved, true
}

// HasOsiApproved returns a boolean if a field has been set.
func (o *License) HasOsiApproved() bool {
	if o != nil && !IsNil(o.OsiApproved) {
		return true
	}

	return false
}

// SetOsiApproved gets a reference to the given bool and assigns it to the OsiApproved field.
func (o *License) SetOsiApproved(v bool) {
	o.OsiApproved = &v
}

// GetFsfLibre returns the FsfLibre field value if set, zero value otherwise.
func (o *License) GetFsfLibre() bool {
	if o == nil || IsNil(o.FsfLibre) {
		var ret bool
		return ret
	}
	return *o.FsfLibre
}

// GetFsfLibreOk returns a tuple with the FsfLibre field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetFsfLibreOk() (*bool, bool) {
	if o == nil || IsNil(o.FsfLibre) {
		return nil, false
	}
	return o.FsfLibre, true
}

// HasFsfLibre returns a boolean if a field has been set.
func (o *License) HasFsfLibre() bool {
	if o != nil && !IsNil(o.FsfLibre) {
		return true
	}

	return false
}

// SetFsfLibre gets a reference to the given bool and assigns it to the FsfLibre field.
func (o *License) SetFsfLibre(v bool) {
	o.FsfLibre = &v
}

// GetCustomLicense returns the CustomLicense field value if set, zero value otherwise.
func (o *License) GetCustomLicense() bool {
	if o == nil || IsNil(o.CustomLicense) {
		var ret bool
		return ret
	}
	return *o.CustomLicense
}

// GetCustomLicenseOk returns a tuple with the CustomLicense field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *License) GetCustomLicenseOk() (*bool, bool) {
	if o == nil || IsNil(o.CustomLicense) {
		return nil, false
	}
	return o.CustomLicense, true
}

// HasCustomLicense returns a boolean if a field has been set.
func (o *License) HasCustomLicense() bool {
	if o != nil && !IsNil(o.CustomLicense) {
		return true
	}

	return false
}

// SetCustomLicense gets a reference to the given bool and assigns it to the CustomLicense field.
func (o *License) SetCustomLicense(v bool) {
	o.CustomLicense = &v
}

func (o License) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o License) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.LicenseId) {
		toSerialize["license_id"] = o.LicenseId
	}
	if !IsNil(o.Uuid) {
		toSerialize["uuid"] = o.Uuid
	}
	if !IsNil(o.OsiApproved) {
		toSerialize["osi_approved"] = o.OsiApproved
	}
	if !IsNil(o.FsfLibre) {
		toSerialize["fsf_libre"] = o.FsfLibre
	}
	if !IsNil(o.CustomLicense) {
		toSerialize["custom_license"] = o.CustomLicense
	}
	return toSerialize, nil
}

type NullableLicense struct {
	value *License
	isSet bool
}

func (v NullableLicense) Get() *License {
	return v.value
}

func (v *NullableLicense) Set(val *License) {
	v.value = val
	v.isSet = true
}

func (v NullableLicense) IsSet() bool {
	return v.isSet
}

func (v *NullableLicense) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLicense(val *License) *NullableLicense {
	return &NullableLicense{value: val, isSet: true}
}

func (v NullableLicense) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLicense) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
