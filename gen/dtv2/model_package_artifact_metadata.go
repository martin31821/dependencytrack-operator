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

// checks if the PackageArtifactMetadata type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PackageArtifactMetadata{}

// PackageArtifactMetadata Artifact-level metadata for the component's exact version from configured package repositories. Only present when the component has a PURL with a version, the PURL type is supported by at least one configured repository, and artifact metadata has been successfully resolved from an upstream repository. Metadata resolution is asynchronous and runs in the background after a component is created or updated. This field may be absent for recently created components until resolution completes.
type PackageArtifactMetadata struct {
	Hashes *Hashes `json:"hashes,omitempty"`
	// When this artifact was published to the repository.
	PublishedAt NullableInt64 `json:"published_at,omitempty"`
	// Identifier of the repository from which artifact metadata was fetched.
	ResolvedFrom NullableString `json:"resolved_from,omitempty"`
	// When artifact metadata was last resolved from the upstream repository.
	ResolvedAt NullableInt64 `json:"resolved_at,omitempty"`
}

// NewPackageArtifactMetadata instantiates a new PackageArtifactMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPackageArtifactMetadata() *PackageArtifactMetadata {
	this := PackageArtifactMetadata{}
	return &this
}

// NewPackageArtifactMetadataWithDefaults instantiates a new PackageArtifactMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPackageArtifactMetadataWithDefaults() *PackageArtifactMetadata {
	this := PackageArtifactMetadata{}
	return &this
}

// GetHashes returns the Hashes field value if set, zero value otherwise.
func (o *PackageArtifactMetadata) GetHashes() Hashes {
	if o == nil || IsNil(o.Hashes) {
		var ret Hashes
		return ret
	}
	return *o.Hashes
}

// GetHashesOk returns a tuple with the Hashes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageArtifactMetadata) GetHashesOk() (*Hashes, bool) {
	if o == nil || IsNil(o.Hashes) {
		return nil, false
	}
	return o.Hashes, true
}

// HasHashes returns a boolean if a field has been set.
func (o *PackageArtifactMetadata) HasHashes() bool {
	if o != nil && !IsNil(o.Hashes) {
		return true
	}

	return false
}

// SetHashes gets a reference to the given Hashes and assigns it to the Hashes field.
func (o *PackageArtifactMetadata) SetHashes(v Hashes) {
	o.Hashes = &v
}

// GetPublishedAt returns the PublishedAt field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PackageArtifactMetadata) GetPublishedAt() int64 {
	if o == nil || IsNil(o.PublishedAt.Get()) {
		var ret int64
		return ret
	}
	return *o.PublishedAt.Get()
}

// GetPublishedAtOk returns a tuple with the PublishedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PackageArtifactMetadata) GetPublishedAtOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return o.PublishedAt.Get(), o.PublishedAt.IsSet()
}

// HasPublishedAt returns a boolean if a field has been set.
func (o *PackageArtifactMetadata) HasPublishedAt() bool {
	if o != nil && o.PublishedAt.IsSet() {
		return true
	}

	return false
}

// SetPublishedAt gets a reference to the given NullableInt64 and assigns it to the PublishedAt field.
func (o *PackageArtifactMetadata) SetPublishedAt(v int64) {
	o.PublishedAt.Set(&v)
}

// SetPublishedAtNil sets the value for PublishedAt to be an explicit nil
func (o *PackageArtifactMetadata) SetPublishedAtNil() {
	o.PublishedAt.Set(nil)
}

// UnsetPublishedAt ensures that no value is present for PublishedAt, not even an explicit nil
func (o *PackageArtifactMetadata) UnsetPublishedAt() {
	o.PublishedAt.Unset()
}

// GetResolvedFrom returns the ResolvedFrom field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PackageArtifactMetadata) GetResolvedFrom() string {
	if o == nil || IsNil(o.ResolvedFrom.Get()) {
		var ret string
		return ret
	}
	return *o.ResolvedFrom.Get()
}

// GetResolvedFromOk returns a tuple with the ResolvedFrom field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PackageArtifactMetadata) GetResolvedFromOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResolvedFrom.Get(), o.ResolvedFrom.IsSet()
}

// HasResolvedFrom returns a boolean if a field has been set.
func (o *PackageArtifactMetadata) HasResolvedFrom() bool {
	if o != nil && o.ResolvedFrom.IsSet() {
		return true
	}

	return false
}

// SetResolvedFrom gets a reference to the given NullableString and assigns it to the ResolvedFrom field.
func (o *PackageArtifactMetadata) SetResolvedFrom(v string) {
	o.ResolvedFrom.Set(&v)
}

// SetResolvedFromNil sets the value for ResolvedFrom to be an explicit nil
func (o *PackageArtifactMetadata) SetResolvedFromNil() {
	o.ResolvedFrom.Set(nil)
}

// UnsetResolvedFrom ensures that no value is present for ResolvedFrom, not even an explicit nil
func (o *PackageArtifactMetadata) UnsetResolvedFrom() {
	o.ResolvedFrom.Unset()
}

// GetResolvedAt returns the ResolvedAt field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PackageArtifactMetadata) GetResolvedAt() int64 {
	if o == nil || IsNil(o.ResolvedAt.Get()) {
		var ret int64
		return ret
	}
	return *o.ResolvedAt.Get()
}

// GetResolvedAtOk returns a tuple with the ResolvedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PackageArtifactMetadata) GetResolvedAtOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResolvedAt.Get(), o.ResolvedAt.IsSet()
}

// HasResolvedAt returns a boolean if a field has been set.
func (o *PackageArtifactMetadata) HasResolvedAt() bool {
	if o != nil && o.ResolvedAt.IsSet() {
		return true
	}

	return false
}

// SetResolvedAt gets a reference to the given NullableInt64 and assigns it to the ResolvedAt field.
func (o *PackageArtifactMetadata) SetResolvedAt(v int64) {
	o.ResolvedAt.Set(&v)
}

// SetResolvedAtNil sets the value for ResolvedAt to be an explicit nil
func (o *PackageArtifactMetadata) SetResolvedAtNil() {
	o.ResolvedAt.Set(nil)
}

// UnsetResolvedAt ensures that no value is present for ResolvedAt, not even an explicit nil
func (o *PackageArtifactMetadata) UnsetResolvedAt() {
	o.ResolvedAt.Unset()
}

func (o PackageArtifactMetadata) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PackageArtifactMetadata) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Hashes) {
		toSerialize["hashes"] = o.Hashes
	}
	if o.PublishedAt.IsSet() {
		toSerialize["published_at"] = o.PublishedAt.Get()
	}
	if o.ResolvedFrom.IsSet() {
		toSerialize["resolved_from"] = o.ResolvedFrom.Get()
	}
	if o.ResolvedAt.IsSet() {
		toSerialize["resolved_at"] = o.ResolvedAt.Get()
	}
	return toSerialize, nil
}

type NullablePackageArtifactMetadata struct {
	value *PackageArtifactMetadata
	isSet bool
}

func (v NullablePackageArtifactMetadata) Get() *PackageArtifactMetadata {
	return v.value
}

func (v *NullablePackageArtifactMetadata) Set(val *PackageArtifactMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullablePackageArtifactMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullablePackageArtifactMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePackageArtifactMetadata(val *PackageArtifactMetadata) *NullablePackageArtifactMetadata {
	return &NullablePackageArtifactMetadata{value: val, isSet: true}
}

func (v NullablePackageArtifactMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePackageArtifactMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
