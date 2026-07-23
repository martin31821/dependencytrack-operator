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

// checks if the ListComponentsResponseItem type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListComponentsResponseItem{}

// ListComponentsResponseItem struct for ListComponentsResponseItem
type ListComponentsResponseItem struct {
	Name                    string                   `json:"name"`
	Version                 *string                  `json:"version,omitempty"`
	Group                   *string                  `json:"group,omitempty"`
	Classifier              *string                  `json:"classifier,omitempty"`
	Scope                   *Scope                   `json:"scope,omitempty"`
	Hashes                  *Hashes                  `json:"hashes,omitempty"`
	Cpe                     *string                  `json:"cpe,omitempty"`
	Purl                    *string                  `json:"purl,omitempty"`
	SwidTagId               *string                  `json:"swid_tag_id,omitempty"`
	Internal                *bool                    `json:"internal,omitempty"`
	Copyright               *string                  `json:"copyright,omitempty"`
	License                 *string                  `json:"license,omitempty"`
	LicenseExpression       *string                  `json:"license_expression,omitempty"`
	LicenseUrl              *string                  `json:"license_url,omitempty"`
	ResolvedLicense         *License                 `json:"resolved_license,omitempty"`
	LastInheritedRiskScore  *float64                 `json:"last_inherited_risk_score,omitempty"`
	Uuid                    string                   `json:"uuid"`
	Project                 *ComponentProject        `json:"project,omitempty"`
	Metrics                 *DependencyMetrics       `json:"metrics,omitempty"`
	PackageMetadata         *PackageMetadata         `json:"package_metadata,omitempty"`
	PackageArtifactMetadata *PackageArtifactMetadata `json:"package_artifact_metadata,omitempty"`
}

type _ListComponentsResponseItem ListComponentsResponseItem

// NewListComponentsResponseItem instantiates a new ListComponentsResponseItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListComponentsResponseItem(name string, uuid string) *ListComponentsResponseItem {
	this := ListComponentsResponseItem{}
	this.Name = name
	this.Uuid = uuid
	return &this
}

// NewListComponentsResponseItemWithDefaults instantiates a new ListComponentsResponseItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListComponentsResponseItemWithDefaults() *ListComponentsResponseItem {
	this := ListComponentsResponseItem{}
	return &this
}

// GetName returns the Name field value
func (o *ListComponentsResponseItem) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ListComponentsResponseItem) SetName(v string) {
	o.Name = v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *ListComponentsResponseItem) SetVersion(v string) {
	o.Version = &v
}

// GetGroup returns the Group field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetGroup() string {
	if o == nil || IsNil(o.Group) {
		var ret string
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetGroupOk() (*string, bool) {
	if o == nil || IsNil(o.Group) {
		return nil, false
	}
	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasGroup() bool {
	if o != nil && !IsNil(o.Group) {
		return true
	}

	return false
}

// SetGroup gets a reference to the given string and assigns it to the Group field.
func (o *ListComponentsResponseItem) SetGroup(v string) {
	o.Group = &v
}

// GetClassifier returns the Classifier field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetClassifier() string {
	if o == nil || IsNil(o.Classifier) {
		var ret string
		return ret
	}
	return *o.Classifier
}

// GetClassifierOk returns a tuple with the Classifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetClassifierOk() (*string, bool) {
	if o == nil || IsNil(o.Classifier) {
		return nil, false
	}
	return o.Classifier, true
}

// HasClassifier returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasClassifier() bool {
	if o != nil && !IsNil(o.Classifier) {
		return true
	}

	return false
}

// SetClassifier gets a reference to the given string and assigns it to the Classifier field.
func (o *ListComponentsResponseItem) SetClassifier(v string) {
	o.Classifier = &v
}

// GetScope returns the Scope field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetScope() Scope {
	if o == nil || IsNil(o.Scope) {
		var ret Scope
		return ret
	}
	return *o.Scope
}

// GetScopeOk returns a tuple with the Scope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetScopeOk() (*Scope, bool) {
	if o == nil || IsNil(o.Scope) {
		return nil, false
	}
	return o.Scope, true
}

// HasScope returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasScope() bool {
	if o != nil && !IsNil(o.Scope) {
		return true
	}

	return false
}

// SetScope gets a reference to the given Scope and assigns it to the Scope field.
func (o *ListComponentsResponseItem) SetScope(v Scope) {
	o.Scope = &v
}

// GetHashes returns the Hashes field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetHashes() Hashes {
	if o == nil || IsNil(o.Hashes) {
		var ret Hashes
		return ret
	}
	return *o.Hashes
}

// GetHashesOk returns a tuple with the Hashes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetHashesOk() (*Hashes, bool) {
	if o == nil || IsNil(o.Hashes) {
		return nil, false
	}
	return o.Hashes, true
}

// HasHashes returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasHashes() bool {
	if o != nil && !IsNil(o.Hashes) {
		return true
	}

	return false
}

// SetHashes gets a reference to the given Hashes and assigns it to the Hashes field.
func (o *ListComponentsResponseItem) SetHashes(v Hashes) {
	o.Hashes = &v
}

// GetCpe returns the Cpe field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetCpe() string {
	if o == nil || IsNil(o.Cpe) {
		var ret string
		return ret
	}
	return *o.Cpe
}

// GetCpeOk returns a tuple with the Cpe field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetCpeOk() (*string, bool) {
	if o == nil || IsNil(o.Cpe) {
		return nil, false
	}
	return o.Cpe, true
}

// HasCpe returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasCpe() bool {
	if o != nil && !IsNil(o.Cpe) {
		return true
	}

	return false
}

// SetCpe gets a reference to the given string and assigns it to the Cpe field.
func (o *ListComponentsResponseItem) SetCpe(v string) {
	o.Cpe = &v
}

// GetPurl returns the Purl field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetPurl() string {
	if o == nil || IsNil(o.Purl) {
		var ret string
		return ret
	}
	return *o.Purl
}

// GetPurlOk returns a tuple with the Purl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetPurlOk() (*string, bool) {
	if o == nil || IsNil(o.Purl) {
		return nil, false
	}
	return o.Purl, true
}

// HasPurl returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasPurl() bool {
	if o != nil && !IsNil(o.Purl) {
		return true
	}

	return false
}

// SetPurl gets a reference to the given string and assigns it to the Purl field.
func (o *ListComponentsResponseItem) SetPurl(v string) {
	o.Purl = &v
}

// GetSwidTagId returns the SwidTagId field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetSwidTagId() string {
	if o == nil || IsNil(o.SwidTagId) {
		var ret string
		return ret
	}
	return *o.SwidTagId
}

// GetSwidTagIdOk returns a tuple with the SwidTagId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetSwidTagIdOk() (*string, bool) {
	if o == nil || IsNil(o.SwidTagId) {
		return nil, false
	}
	return o.SwidTagId, true
}

// HasSwidTagId returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasSwidTagId() bool {
	if o != nil && !IsNil(o.SwidTagId) {
		return true
	}

	return false
}

// SetSwidTagId gets a reference to the given string and assigns it to the SwidTagId field.
func (o *ListComponentsResponseItem) SetSwidTagId(v string) {
	o.SwidTagId = &v
}

// GetInternal returns the Internal field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetInternal() bool {
	if o == nil || IsNil(o.Internal) {
		var ret bool
		return ret
	}
	return *o.Internal
}

// GetInternalOk returns a tuple with the Internal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetInternalOk() (*bool, bool) {
	if o == nil || IsNil(o.Internal) {
		return nil, false
	}
	return o.Internal, true
}

// HasInternal returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasInternal() bool {
	if o != nil && !IsNil(o.Internal) {
		return true
	}

	return false
}

// SetInternal gets a reference to the given bool and assigns it to the Internal field.
func (o *ListComponentsResponseItem) SetInternal(v bool) {
	o.Internal = &v
}

// GetCopyright returns the Copyright field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetCopyright() string {
	if o == nil || IsNil(o.Copyright) {
		var ret string
		return ret
	}
	return *o.Copyright
}

// GetCopyrightOk returns a tuple with the Copyright field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetCopyrightOk() (*string, bool) {
	if o == nil || IsNil(o.Copyright) {
		return nil, false
	}
	return o.Copyright, true
}

// HasCopyright returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasCopyright() bool {
	if o != nil && !IsNil(o.Copyright) {
		return true
	}

	return false
}

// SetCopyright gets a reference to the given string and assigns it to the Copyright field.
func (o *ListComponentsResponseItem) SetCopyright(v string) {
	o.Copyright = &v
}

// GetLicense returns the License field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetLicense() string {
	if o == nil || IsNil(o.License) {
		var ret string
		return ret
	}
	return *o.License
}

// GetLicenseOk returns a tuple with the License field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetLicenseOk() (*string, bool) {
	if o == nil || IsNil(o.License) {
		return nil, false
	}
	return o.License, true
}

// HasLicense returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasLicense() bool {
	if o != nil && !IsNil(o.License) {
		return true
	}

	return false
}

// SetLicense gets a reference to the given string and assigns it to the License field.
func (o *ListComponentsResponseItem) SetLicense(v string) {
	o.License = &v
}

// GetLicenseExpression returns the LicenseExpression field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetLicenseExpression() string {
	if o == nil || IsNil(o.LicenseExpression) {
		var ret string
		return ret
	}
	return *o.LicenseExpression
}

// GetLicenseExpressionOk returns a tuple with the LicenseExpression field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetLicenseExpressionOk() (*string, bool) {
	if o == nil || IsNil(o.LicenseExpression) {
		return nil, false
	}
	return o.LicenseExpression, true
}

// HasLicenseExpression returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasLicenseExpression() bool {
	if o != nil && !IsNil(o.LicenseExpression) {
		return true
	}

	return false
}

// SetLicenseExpression gets a reference to the given string and assigns it to the LicenseExpression field.
func (o *ListComponentsResponseItem) SetLicenseExpression(v string) {
	o.LicenseExpression = &v
}

// GetLicenseUrl returns the LicenseUrl field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetLicenseUrl() string {
	if o == nil || IsNil(o.LicenseUrl) {
		var ret string
		return ret
	}
	return *o.LicenseUrl
}

// GetLicenseUrlOk returns a tuple with the LicenseUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetLicenseUrlOk() (*string, bool) {
	if o == nil || IsNil(o.LicenseUrl) {
		return nil, false
	}
	return o.LicenseUrl, true
}

// HasLicenseUrl returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasLicenseUrl() bool {
	if o != nil && !IsNil(o.LicenseUrl) {
		return true
	}

	return false
}

// SetLicenseUrl gets a reference to the given string and assigns it to the LicenseUrl field.
func (o *ListComponentsResponseItem) SetLicenseUrl(v string) {
	o.LicenseUrl = &v
}

// GetResolvedLicense returns the ResolvedLicense field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetResolvedLicense() License {
	if o == nil || IsNil(o.ResolvedLicense) {
		var ret License
		return ret
	}
	return *o.ResolvedLicense
}

// GetResolvedLicenseOk returns a tuple with the ResolvedLicense field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetResolvedLicenseOk() (*License, bool) {
	if o == nil || IsNil(o.ResolvedLicense) {
		return nil, false
	}
	return o.ResolvedLicense, true
}

// HasResolvedLicense returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasResolvedLicense() bool {
	if o != nil && !IsNil(o.ResolvedLicense) {
		return true
	}

	return false
}

// SetResolvedLicense gets a reference to the given License and assigns it to the ResolvedLicense field.
func (o *ListComponentsResponseItem) SetResolvedLicense(v License) {
	o.ResolvedLicense = &v
}

// GetLastInheritedRiskScore returns the LastInheritedRiskScore field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetLastInheritedRiskScore() float64 {
	if o == nil || IsNil(o.LastInheritedRiskScore) {
		var ret float64
		return ret
	}
	return *o.LastInheritedRiskScore
}

// GetLastInheritedRiskScoreOk returns a tuple with the LastInheritedRiskScore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetLastInheritedRiskScoreOk() (*float64, bool) {
	if o == nil || IsNil(o.LastInheritedRiskScore) {
		return nil, false
	}
	return o.LastInheritedRiskScore, true
}

// HasLastInheritedRiskScore returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasLastInheritedRiskScore() bool {
	if o != nil && !IsNil(o.LastInheritedRiskScore) {
		return true
	}

	return false
}

// SetLastInheritedRiskScore gets a reference to the given float64 and assigns it to the LastInheritedRiskScore field.
func (o *ListComponentsResponseItem) SetLastInheritedRiskScore(v float64) {
	o.LastInheritedRiskScore = &v
}

// GetUuid returns the Uuid field value
func (o *ListComponentsResponseItem) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *ListComponentsResponseItem) SetUuid(v string) {
	o.Uuid = v
}

// GetProject returns the Project field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetProject() ComponentProject {
	if o == nil || IsNil(o.Project) {
		var ret ComponentProject
		return ret
	}
	return *o.Project
}

// GetProjectOk returns a tuple with the Project field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetProjectOk() (*ComponentProject, bool) {
	if o == nil || IsNil(o.Project) {
		return nil, false
	}
	return o.Project, true
}

// HasProject returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasProject() bool {
	if o != nil && !IsNil(o.Project) {
		return true
	}

	return false
}

// SetProject gets a reference to the given ComponentProject and assigns it to the Project field.
func (o *ListComponentsResponseItem) SetProject(v ComponentProject) {
	o.Project = &v
}

// GetMetrics returns the Metrics field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetMetrics() DependencyMetrics {
	if o == nil || IsNil(o.Metrics) {
		var ret DependencyMetrics
		return ret
	}
	return *o.Metrics
}

// GetMetricsOk returns a tuple with the Metrics field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetMetricsOk() (*DependencyMetrics, bool) {
	if o == nil || IsNil(o.Metrics) {
		return nil, false
	}
	return o.Metrics, true
}

// HasMetrics returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasMetrics() bool {
	if o != nil && !IsNil(o.Metrics) {
		return true
	}

	return false
}

// SetMetrics gets a reference to the given DependencyMetrics and assigns it to the Metrics field.
func (o *ListComponentsResponseItem) SetMetrics(v DependencyMetrics) {
	o.Metrics = &v
}

// GetPackageMetadata returns the PackageMetadata field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetPackageMetadata() PackageMetadata {
	if o == nil || IsNil(o.PackageMetadata) {
		var ret PackageMetadata
		return ret
	}
	return *o.PackageMetadata
}

// GetPackageMetadataOk returns a tuple with the PackageMetadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetPackageMetadataOk() (*PackageMetadata, bool) {
	if o == nil || IsNil(o.PackageMetadata) {
		return nil, false
	}
	return o.PackageMetadata, true
}

// HasPackageMetadata returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasPackageMetadata() bool {
	if o != nil && !IsNil(o.PackageMetadata) {
		return true
	}

	return false
}

// SetPackageMetadata gets a reference to the given PackageMetadata and assigns it to the PackageMetadata field.
func (o *ListComponentsResponseItem) SetPackageMetadata(v PackageMetadata) {
	o.PackageMetadata = &v
}

// GetPackageArtifactMetadata returns the PackageArtifactMetadata field value if set, zero value otherwise.
func (o *ListComponentsResponseItem) GetPackageArtifactMetadata() PackageArtifactMetadata {
	if o == nil || IsNil(o.PackageArtifactMetadata) {
		var ret PackageArtifactMetadata
		return ret
	}
	return *o.PackageArtifactMetadata
}

// GetPackageArtifactMetadataOk returns a tuple with the PackageArtifactMetadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListComponentsResponseItem) GetPackageArtifactMetadataOk() (*PackageArtifactMetadata, bool) {
	if o == nil || IsNil(o.PackageArtifactMetadata) {
		return nil, false
	}
	return o.PackageArtifactMetadata, true
}

// HasPackageArtifactMetadata returns a boolean if a field has been set.
func (o *ListComponentsResponseItem) HasPackageArtifactMetadata() bool {
	if o != nil && !IsNil(o.PackageArtifactMetadata) {
		return true
	}

	return false
}

// SetPackageArtifactMetadata gets a reference to the given PackageArtifactMetadata and assigns it to the PackageArtifactMetadata field.
func (o *ListComponentsResponseItem) SetPackageArtifactMetadata(v PackageArtifactMetadata) {
	o.PackageArtifactMetadata = &v
}

func (o ListComponentsResponseItem) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListComponentsResponseItem) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.Group) {
		toSerialize["group"] = o.Group
	}
	if !IsNil(o.Classifier) {
		toSerialize["classifier"] = o.Classifier
	}
	if !IsNil(o.Scope) {
		toSerialize["scope"] = o.Scope
	}
	if !IsNil(o.Hashes) {
		toSerialize["hashes"] = o.Hashes
	}
	if !IsNil(o.Cpe) {
		toSerialize["cpe"] = o.Cpe
	}
	if !IsNil(o.Purl) {
		toSerialize["purl"] = o.Purl
	}
	if !IsNil(o.SwidTagId) {
		toSerialize["swid_tag_id"] = o.SwidTagId
	}
	if !IsNil(o.Internal) {
		toSerialize["internal"] = o.Internal
	}
	if !IsNil(o.Copyright) {
		toSerialize["copyright"] = o.Copyright
	}
	if !IsNil(o.License) {
		toSerialize["license"] = o.License
	}
	if !IsNil(o.LicenseExpression) {
		toSerialize["license_expression"] = o.LicenseExpression
	}
	if !IsNil(o.LicenseUrl) {
		toSerialize["license_url"] = o.LicenseUrl
	}
	if !IsNil(o.ResolvedLicense) {
		toSerialize["resolved_license"] = o.ResolvedLicense
	}
	if !IsNil(o.LastInheritedRiskScore) {
		toSerialize["last_inherited_risk_score"] = o.LastInheritedRiskScore
	}
	toSerialize["uuid"] = o.Uuid
	if !IsNil(o.Project) {
		toSerialize["project"] = o.Project
	}
	if !IsNil(o.Metrics) {
		toSerialize["metrics"] = o.Metrics
	}
	if !IsNil(o.PackageMetadata) {
		toSerialize["package_metadata"] = o.PackageMetadata
	}
	if !IsNil(o.PackageArtifactMetadata) {
		toSerialize["package_artifact_metadata"] = o.PackageArtifactMetadata
	}
	return toSerialize, nil
}

func (o *ListComponentsResponseItem) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"uuid",
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

	varListComponentsResponseItem := _ListComponentsResponseItem{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varListComponentsResponseItem)

	if err != nil {
		return err
	}

	*o = ListComponentsResponseItem(varListComponentsResponseItem)

	return err
}

type NullableListComponentsResponseItem struct {
	value *ListComponentsResponseItem
	isSet bool
}

func (v NullableListComponentsResponseItem) Get() *ListComponentsResponseItem {
	return v.value
}

func (v *NullableListComponentsResponseItem) Set(val *ListComponentsResponseItem) {
	v.value = val
	v.isSet = true
}

func (v NullableListComponentsResponseItem) IsSet() bool {
	return v.isSet
}

func (v *NullableListComponentsResponseItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListComponentsResponseItem(val *ListComponentsResponseItem) *NullableListComponentsResponseItem {
	return &NullableListComponentsResponseItem{value: val, isSet: true}
}

func (v NullableListComponentsResponseItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListComponentsResponseItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
