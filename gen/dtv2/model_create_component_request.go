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

// checks if the CreateComponentRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateComponentRequest{}

// CreateComponentRequest struct for CreateComponentRequest
type CreateComponentRequest struct {
	ProjectUuid       string                  `json:"project_uuid"`
	Name              string                  `json:"name"`
	Description       *string                 `json:"description,omitempty"`
	Group             *string                 `json:"group,omitempty"`
	Version           *string                 `json:"version,omitempty"`
	Classifier        *Classifier             `json:"classifier,omitempty"`
	Filename          *string                 `json:"filename,omitempty"`
	Extension         *string                 `json:"extension,omitempty"`
	Hashes            *Hashes                 `json:"hashes,omitempty"`
	Cpe               *string                 `json:"cpe,omitempty"`
	Publisher         *string                 `json:"publisher,omitempty"`
	Supplier          *OrganizationalEntity   `json:"supplier,omitempty"`
	Authors           []OrganizationalContact `json:"authors,omitempty"`
	Purl              *string                 `json:"purl,omitempty"`
	SwidTagId         *string                 `json:"swid_tag_id,omitempty"`
	Internal          *bool                   `json:"internal,omitempty"`
	Copyright         *string                 `json:"copyright,omitempty"`
	License           *string                 `json:"license,omitempty"`
	LicenseExpression *string                 `json:"license_expression,omitempty"`
	LicenseUrl        *string                 `json:"license_url,omitempty"`
	Notes             *string                 `json:"notes,omitempty"`
}

type _CreateComponentRequest CreateComponentRequest

// NewCreateComponentRequest instantiates a new CreateComponentRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateComponentRequest(projectUuid string, name string) *CreateComponentRequest {
	this := CreateComponentRequest{}
	this.ProjectUuid = projectUuid
	this.Name = name
	return &this
}

// NewCreateComponentRequestWithDefaults instantiates a new CreateComponentRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateComponentRequestWithDefaults() *CreateComponentRequest {
	this := CreateComponentRequest{}
	return &this
}

// GetProjectUuid returns the ProjectUuid field value
func (o *CreateComponentRequest) GetProjectUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProjectUuid
}

// GetProjectUuidOk returns a tuple with the ProjectUuid field value
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetProjectUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProjectUuid, true
}

// SetProjectUuid sets field value
func (o *CreateComponentRequest) SetProjectUuid(v string) {
	o.ProjectUuid = v
}

// GetName returns the Name field value
func (o *CreateComponentRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CreateComponentRequest) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *CreateComponentRequest) SetDescription(v string) {
	o.Description = &v
}

// GetGroup returns the Group field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetGroup() string {
	if o == nil || IsNil(o.Group) {
		var ret string
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetGroupOk() (*string, bool) {
	if o == nil || IsNil(o.Group) {
		return nil, false
	}
	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasGroup() bool {
	if o != nil && !IsNil(o.Group) {
		return true
	}

	return false
}

// SetGroup gets a reference to the given string and assigns it to the Group field.
func (o *CreateComponentRequest) SetGroup(v string) {
	o.Group = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *CreateComponentRequest) SetVersion(v string) {
	o.Version = &v
}

// GetClassifier returns the Classifier field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetClassifier() Classifier {
	if o == nil || IsNil(o.Classifier) {
		var ret Classifier
		return ret
	}
	return *o.Classifier
}

// GetClassifierOk returns a tuple with the Classifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetClassifierOk() (*Classifier, bool) {
	if o == nil || IsNil(o.Classifier) {
		return nil, false
	}
	return o.Classifier, true
}

// HasClassifier returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasClassifier() bool {
	if o != nil && !IsNil(o.Classifier) {
		return true
	}

	return false
}

// SetClassifier gets a reference to the given Classifier and assigns it to the Classifier field.
func (o *CreateComponentRequest) SetClassifier(v Classifier) {
	o.Classifier = &v
}

// GetFilename returns the Filename field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetFilename() string {
	if o == nil || IsNil(o.Filename) {
		var ret string
		return ret
	}
	return *o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetFilenameOk() (*string, bool) {
	if o == nil || IsNil(o.Filename) {
		return nil, false
	}
	return o.Filename, true
}

// HasFilename returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasFilename() bool {
	if o != nil && !IsNil(o.Filename) {
		return true
	}

	return false
}

// SetFilename gets a reference to the given string and assigns it to the Filename field.
func (o *CreateComponentRequest) SetFilename(v string) {
	o.Filename = &v
}

// GetExtension returns the Extension field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetExtension() string {
	if o == nil || IsNil(o.Extension) {
		var ret string
		return ret
	}
	return *o.Extension
}

// GetExtensionOk returns a tuple with the Extension field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetExtensionOk() (*string, bool) {
	if o == nil || IsNil(o.Extension) {
		return nil, false
	}
	return o.Extension, true
}

// HasExtension returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasExtension() bool {
	if o != nil && !IsNil(o.Extension) {
		return true
	}

	return false
}

// SetExtension gets a reference to the given string and assigns it to the Extension field.
func (o *CreateComponentRequest) SetExtension(v string) {
	o.Extension = &v
}

// GetHashes returns the Hashes field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetHashes() Hashes {
	if o == nil || IsNil(o.Hashes) {
		var ret Hashes
		return ret
	}
	return *o.Hashes
}

// GetHashesOk returns a tuple with the Hashes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetHashesOk() (*Hashes, bool) {
	if o == nil || IsNil(o.Hashes) {
		return nil, false
	}
	return o.Hashes, true
}

// HasHashes returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasHashes() bool {
	if o != nil && !IsNil(o.Hashes) {
		return true
	}

	return false
}

// SetHashes gets a reference to the given Hashes and assigns it to the Hashes field.
func (o *CreateComponentRequest) SetHashes(v Hashes) {
	o.Hashes = &v
}

// GetCpe returns the Cpe field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetCpe() string {
	if o == nil || IsNil(o.Cpe) {
		var ret string
		return ret
	}
	return *o.Cpe
}

// GetCpeOk returns a tuple with the Cpe field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetCpeOk() (*string, bool) {
	if o == nil || IsNil(o.Cpe) {
		return nil, false
	}
	return o.Cpe, true
}

// HasCpe returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasCpe() bool {
	if o != nil && !IsNil(o.Cpe) {
		return true
	}

	return false
}

// SetCpe gets a reference to the given string and assigns it to the Cpe field.
func (o *CreateComponentRequest) SetCpe(v string) {
	o.Cpe = &v
}

// GetPublisher returns the Publisher field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetPublisher() string {
	if o == nil || IsNil(o.Publisher) {
		var ret string
		return ret
	}
	return *o.Publisher
}

// GetPublisherOk returns a tuple with the Publisher field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetPublisherOk() (*string, bool) {
	if o == nil || IsNil(o.Publisher) {
		return nil, false
	}
	return o.Publisher, true
}

// HasPublisher returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasPublisher() bool {
	if o != nil && !IsNil(o.Publisher) {
		return true
	}

	return false
}

// SetPublisher gets a reference to the given string and assigns it to the Publisher field.
func (o *CreateComponentRequest) SetPublisher(v string) {
	o.Publisher = &v
}

// GetSupplier returns the Supplier field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetSupplier() OrganizationalEntity {
	if o == nil || IsNil(o.Supplier) {
		var ret OrganizationalEntity
		return ret
	}
	return *o.Supplier
}

// GetSupplierOk returns a tuple with the Supplier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetSupplierOk() (*OrganizationalEntity, bool) {
	if o == nil || IsNil(o.Supplier) {
		return nil, false
	}
	return o.Supplier, true
}

// HasSupplier returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasSupplier() bool {
	if o != nil && !IsNil(o.Supplier) {
		return true
	}

	return false
}

// SetSupplier gets a reference to the given OrganizationalEntity and assigns it to the Supplier field.
func (o *CreateComponentRequest) SetSupplier(v OrganizationalEntity) {
	o.Supplier = &v
}

// GetAuthors returns the Authors field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetAuthors() []OrganizationalContact {
	if o == nil || IsNil(o.Authors) {
		var ret []OrganizationalContact
		return ret
	}
	return o.Authors
}

// GetAuthorsOk returns a tuple with the Authors field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetAuthorsOk() ([]OrganizationalContact, bool) {
	if o == nil || IsNil(o.Authors) {
		return nil, false
	}
	return o.Authors, true
}

// HasAuthors returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasAuthors() bool {
	if o != nil && !IsNil(o.Authors) {
		return true
	}

	return false
}

// SetAuthors gets a reference to the given []OrganizationalContact and assigns it to the Authors field.
func (o *CreateComponentRequest) SetAuthors(v []OrganizationalContact) {
	o.Authors = v
}

// GetPurl returns the Purl field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetPurl() string {
	if o == nil || IsNil(o.Purl) {
		var ret string
		return ret
	}
	return *o.Purl
}

// GetPurlOk returns a tuple with the Purl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetPurlOk() (*string, bool) {
	if o == nil || IsNil(o.Purl) {
		return nil, false
	}
	return o.Purl, true
}

// HasPurl returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasPurl() bool {
	if o != nil && !IsNil(o.Purl) {
		return true
	}

	return false
}

// SetPurl gets a reference to the given string and assigns it to the Purl field.
func (o *CreateComponentRequest) SetPurl(v string) {
	o.Purl = &v
}

// GetSwidTagId returns the SwidTagId field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetSwidTagId() string {
	if o == nil || IsNil(o.SwidTagId) {
		var ret string
		return ret
	}
	return *o.SwidTagId
}

// GetSwidTagIdOk returns a tuple with the SwidTagId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetSwidTagIdOk() (*string, bool) {
	if o == nil || IsNil(o.SwidTagId) {
		return nil, false
	}
	return o.SwidTagId, true
}

// HasSwidTagId returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasSwidTagId() bool {
	if o != nil && !IsNil(o.SwidTagId) {
		return true
	}

	return false
}

// SetSwidTagId gets a reference to the given string and assigns it to the SwidTagId field.
func (o *CreateComponentRequest) SetSwidTagId(v string) {
	o.SwidTagId = &v
}

// GetInternal returns the Internal field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetInternal() bool {
	if o == nil || IsNil(o.Internal) {
		var ret bool
		return ret
	}
	return *o.Internal
}

// GetInternalOk returns a tuple with the Internal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetInternalOk() (*bool, bool) {
	if o == nil || IsNil(o.Internal) {
		return nil, false
	}
	return o.Internal, true
}

// HasInternal returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasInternal() bool {
	if o != nil && !IsNil(o.Internal) {
		return true
	}

	return false
}

// SetInternal gets a reference to the given bool and assigns it to the Internal field.
func (o *CreateComponentRequest) SetInternal(v bool) {
	o.Internal = &v
}

// GetCopyright returns the Copyright field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetCopyright() string {
	if o == nil || IsNil(o.Copyright) {
		var ret string
		return ret
	}
	return *o.Copyright
}

// GetCopyrightOk returns a tuple with the Copyright field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetCopyrightOk() (*string, bool) {
	if o == nil || IsNil(o.Copyright) {
		return nil, false
	}
	return o.Copyright, true
}

// HasCopyright returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasCopyright() bool {
	if o != nil && !IsNil(o.Copyright) {
		return true
	}

	return false
}

// SetCopyright gets a reference to the given string and assigns it to the Copyright field.
func (o *CreateComponentRequest) SetCopyright(v string) {
	o.Copyright = &v
}

// GetLicense returns the License field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetLicense() string {
	if o == nil || IsNil(o.License) {
		var ret string
		return ret
	}
	return *o.License
}

// GetLicenseOk returns a tuple with the License field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetLicenseOk() (*string, bool) {
	if o == nil || IsNil(o.License) {
		return nil, false
	}
	return o.License, true
}

// HasLicense returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasLicense() bool {
	if o != nil && !IsNil(o.License) {
		return true
	}

	return false
}

// SetLicense gets a reference to the given string and assigns it to the License field.
func (o *CreateComponentRequest) SetLicense(v string) {
	o.License = &v
}

// GetLicenseExpression returns the LicenseExpression field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetLicenseExpression() string {
	if o == nil || IsNil(o.LicenseExpression) {
		var ret string
		return ret
	}
	return *o.LicenseExpression
}

// GetLicenseExpressionOk returns a tuple with the LicenseExpression field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetLicenseExpressionOk() (*string, bool) {
	if o == nil || IsNil(o.LicenseExpression) {
		return nil, false
	}
	return o.LicenseExpression, true
}

// HasLicenseExpression returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasLicenseExpression() bool {
	if o != nil && !IsNil(o.LicenseExpression) {
		return true
	}

	return false
}

// SetLicenseExpression gets a reference to the given string and assigns it to the LicenseExpression field.
func (o *CreateComponentRequest) SetLicenseExpression(v string) {
	o.LicenseExpression = &v
}

// GetLicenseUrl returns the LicenseUrl field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetLicenseUrl() string {
	if o == nil || IsNil(o.LicenseUrl) {
		var ret string
		return ret
	}
	return *o.LicenseUrl
}

// GetLicenseUrlOk returns a tuple with the LicenseUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetLicenseUrlOk() (*string, bool) {
	if o == nil || IsNil(o.LicenseUrl) {
		return nil, false
	}
	return o.LicenseUrl, true
}

// HasLicenseUrl returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasLicenseUrl() bool {
	if o != nil && !IsNil(o.LicenseUrl) {
		return true
	}

	return false
}

// SetLicenseUrl gets a reference to the given string and assigns it to the LicenseUrl field.
func (o *CreateComponentRequest) SetLicenseUrl(v string) {
	o.LicenseUrl = &v
}

// GetNotes returns the Notes field value if set, zero value otherwise.
func (o *CreateComponentRequest) GetNotes() string {
	if o == nil || IsNil(o.Notes) {
		var ret string
		return ret
	}
	return *o.Notes
}

// GetNotesOk returns a tuple with the Notes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateComponentRequest) GetNotesOk() (*string, bool) {
	if o == nil || IsNil(o.Notes) {
		return nil, false
	}
	return o.Notes, true
}

// HasNotes returns a boolean if a field has been set.
func (o *CreateComponentRequest) HasNotes() bool {
	if o != nil && !IsNil(o.Notes) {
		return true
	}

	return false
}

// SetNotes gets a reference to the given string and assigns it to the Notes field.
func (o *CreateComponentRequest) SetNotes(v string) {
	o.Notes = &v
}

func (o CreateComponentRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateComponentRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["project_uuid"] = o.ProjectUuid
	toSerialize["name"] = o.Name
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Group) {
		toSerialize["group"] = o.Group
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.Classifier) {
		toSerialize["classifier"] = o.Classifier
	}
	if !IsNil(o.Filename) {
		toSerialize["filename"] = o.Filename
	}
	if !IsNil(o.Extension) {
		toSerialize["extension"] = o.Extension
	}
	if !IsNil(o.Hashes) {
		toSerialize["hashes"] = o.Hashes
	}
	if !IsNil(o.Cpe) {
		toSerialize["cpe"] = o.Cpe
	}
	if !IsNil(o.Publisher) {
		toSerialize["publisher"] = o.Publisher
	}
	if !IsNil(o.Supplier) {
		toSerialize["supplier"] = o.Supplier
	}
	if !IsNil(o.Authors) {
		toSerialize["authors"] = o.Authors
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
	if !IsNil(o.Notes) {
		toSerialize["notes"] = o.Notes
	}
	return toSerialize, nil
}

func (o *CreateComponentRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"project_uuid",
		"name",
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

	varCreateComponentRequest := _CreateComponentRequest{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCreateComponentRequest)

	if err != nil {
		return err
	}

	*o = CreateComponentRequest(varCreateComponentRequest)

	return err
}

type NullableCreateComponentRequest struct {
	value *CreateComponentRequest
	isSet bool
}

func (v NullableCreateComponentRequest) Get() *CreateComponentRequest {
	return v.value
}

func (v *NullableCreateComponentRequest) Set(val *CreateComponentRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateComponentRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateComponentRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateComponentRequest(val *CreateComponentRequest) *NullableCreateComponentRequest {
	return &NullableCreateComponentRequest{value: val, isSet: true}
}

func (v NullableCreateComponentRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateComponentRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
