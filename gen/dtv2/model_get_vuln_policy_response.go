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

// checks if the GetVulnPolicyResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetVulnPolicyResponse{}

// GetVulnPolicyResponse struct for GetVulnPolicyResponse
type GetVulnPolicyResponse struct {
	Uuid          string                  `json:"uuid"`
	Name          string                  `json:"name"`
	Description   *string                 `json:"description,omitempty"`
	Author        *string                 `json:"author,omitempty"`
	Condition     string                  `json:"condition"`
	Analysis      VulnPolicyAnalysis      `json:"analysis"`
	Ratings       []VulnPolicyRating      `json:"ratings,omitempty"`
	OperationMode VulnPolicyOperationMode `json:"operation_mode"`
	Priority      int32                   `json:"priority"`
	Source        VulnPolicySource        `json:"source"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	ValidFrom *int64 `json:"valid_from,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	ValidUntil *int64 `json:"valid_until,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	Created *int64 `json:"created,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	Updated *int64 `json:"updated,omitempty"`
}

type _GetVulnPolicyResponse GetVulnPolicyResponse

// NewGetVulnPolicyResponse instantiates a new GetVulnPolicyResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetVulnPolicyResponse(uuid string, name string, condition string, analysis VulnPolicyAnalysis, operationMode VulnPolicyOperationMode, priority int32, source VulnPolicySource) *GetVulnPolicyResponse {
	this := GetVulnPolicyResponse{}
	this.Uuid = uuid
	this.Name = name
	this.Condition = condition
	this.Analysis = analysis
	this.OperationMode = operationMode
	this.Priority = priority
	this.Source = source
	return &this
}

// NewGetVulnPolicyResponseWithDefaults instantiates a new GetVulnPolicyResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetVulnPolicyResponseWithDefaults() *GetVulnPolicyResponse {
	this := GetVulnPolicyResponse{}
	return &this
}

// GetUuid returns the Uuid field value
func (o *GetVulnPolicyResponse) GetUuid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetUuidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uuid, true
}

// SetUuid sets field value
func (o *GetVulnPolicyResponse) SetUuid(v string) {
	o.Uuid = v
}

// GetName returns the Name field value
func (o *GetVulnPolicyResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GetVulnPolicyResponse) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *GetVulnPolicyResponse) SetDescription(v string) {
	o.Description = &v
}

// GetAuthor returns the Author field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetAuthor() string {
	if o == nil || IsNil(o.Author) {
		var ret string
		return ret
	}
	return *o.Author
}

// GetAuthorOk returns a tuple with the Author field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetAuthorOk() (*string, bool) {
	if o == nil || IsNil(o.Author) {
		return nil, false
	}
	return o.Author, true
}

// HasAuthor returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasAuthor() bool {
	if o != nil && !IsNil(o.Author) {
		return true
	}

	return false
}

// SetAuthor gets a reference to the given string and assigns it to the Author field.
func (o *GetVulnPolicyResponse) SetAuthor(v string) {
	o.Author = &v
}

// GetCondition returns the Condition field value
func (o *GetVulnPolicyResponse) GetCondition() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Condition
}

// GetConditionOk returns a tuple with the Condition field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetConditionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Condition, true
}

// SetCondition sets field value
func (o *GetVulnPolicyResponse) SetCondition(v string) {
	o.Condition = v
}

// GetAnalysis returns the Analysis field value
func (o *GetVulnPolicyResponse) GetAnalysis() VulnPolicyAnalysis {
	if o == nil {
		var ret VulnPolicyAnalysis
		return ret
	}

	return o.Analysis
}

// GetAnalysisOk returns a tuple with the Analysis field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetAnalysisOk() (*VulnPolicyAnalysis, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Analysis, true
}

// SetAnalysis sets field value
func (o *GetVulnPolicyResponse) SetAnalysis(v VulnPolicyAnalysis) {
	o.Analysis = v
}

// GetRatings returns the Ratings field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetRatings() []VulnPolicyRating {
	if o == nil || IsNil(o.Ratings) {
		var ret []VulnPolicyRating
		return ret
	}
	return o.Ratings
}

// GetRatingsOk returns a tuple with the Ratings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetRatingsOk() ([]VulnPolicyRating, bool) {
	if o == nil || IsNil(o.Ratings) {
		return nil, false
	}
	return o.Ratings, true
}

// HasRatings returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasRatings() bool {
	if o != nil && !IsNil(o.Ratings) {
		return true
	}

	return false
}

// SetRatings gets a reference to the given []VulnPolicyRating and assigns it to the Ratings field.
func (o *GetVulnPolicyResponse) SetRatings(v []VulnPolicyRating) {
	o.Ratings = v
}

// GetOperationMode returns the OperationMode field value
func (o *GetVulnPolicyResponse) GetOperationMode() VulnPolicyOperationMode {
	if o == nil {
		var ret VulnPolicyOperationMode
		return ret
	}

	return o.OperationMode
}

// GetOperationModeOk returns a tuple with the OperationMode field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetOperationModeOk() (*VulnPolicyOperationMode, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OperationMode, true
}

// SetOperationMode sets field value
func (o *GetVulnPolicyResponse) SetOperationMode(v VulnPolicyOperationMode) {
	o.OperationMode = v
}

// GetPriority returns the Priority field value
func (o *GetVulnPolicyResponse) GetPriority() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetPriorityOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Priority, true
}

// SetPriority sets field value
func (o *GetVulnPolicyResponse) SetPriority(v int32) {
	o.Priority = v
}

// GetSource returns the Source field value
func (o *GetVulnPolicyResponse) GetSource() VulnPolicySource {
	if o == nil {
		var ret VulnPolicySource
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetSourceOk() (*VulnPolicySource, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *GetVulnPolicyResponse) SetSource(v VulnPolicySource) {
	o.Source = v
}

// GetValidFrom returns the ValidFrom field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetValidFrom() int64 {
	if o == nil || IsNil(o.ValidFrom) {
		var ret int64
		return ret
	}
	return *o.ValidFrom
}

// GetValidFromOk returns a tuple with the ValidFrom field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetValidFromOk() (*int64, bool) {
	if o == nil || IsNil(o.ValidFrom) {
		return nil, false
	}
	return o.ValidFrom, true
}

// HasValidFrom returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasValidFrom() bool {
	if o != nil && !IsNil(o.ValidFrom) {
		return true
	}

	return false
}

// SetValidFrom gets a reference to the given int64 and assigns it to the ValidFrom field.
func (o *GetVulnPolicyResponse) SetValidFrom(v int64) {
	o.ValidFrom = &v
}

// GetValidUntil returns the ValidUntil field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetValidUntil() int64 {
	if o == nil || IsNil(o.ValidUntil) {
		var ret int64
		return ret
	}
	return *o.ValidUntil
}

// GetValidUntilOk returns a tuple with the ValidUntil field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetValidUntilOk() (*int64, bool) {
	if o == nil || IsNil(o.ValidUntil) {
		return nil, false
	}
	return o.ValidUntil, true
}

// HasValidUntil returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasValidUntil() bool {
	if o != nil && !IsNil(o.ValidUntil) {
		return true
	}

	return false
}

// SetValidUntil gets a reference to the given int64 and assigns it to the ValidUntil field.
func (o *GetVulnPolicyResponse) SetValidUntil(v int64) {
	o.ValidUntil = &v
}

// GetCreated returns the Created field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetCreated() int64 {
	if o == nil || IsNil(o.Created) {
		var ret int64
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetCreatedOk() (*int64, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}
	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given int64 and assigns it to the Created field.
func (o *GetVulnPolicyResponse) SetCreated(v int64) {
	o.Created = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise.
func (o *GetVulnPolicyResponse) GetUpdated() int64 {
	if o == nil || IsNil(o.Updated) {
		var ret int64
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetVulnPolicyResponse) GetUpdatedOk() (*int64, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}
	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *GetVulnPolicyResponse) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given int64 and assigns it to the Updated field.
func (o *GetVulnPolicyResponse) SetUpdated(v int64) {
	o.Updated = &v
}

func (o GetVulnPolicyResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetVulnPolicyResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uuid"] = o.Uuid
	toSerialize["name"] = o.Name
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Author) {
		toSerialize["author"] = o.Author
	}
	toSerialize["condition"] = o.Condition
	toSerialize["analysis"] = o.Analysis
	if !IsNil(o.Ratings) {
		toSerialize["ratings"] = o.Ratings
	}
	toSerialize["operation_mode"] = o.OperationMode
	toSerialize["priority"] = o.Priority
	toSerialize["source"] = o.Source
	if !IsNil(o.ValidFrom) {
		toSerialize["valid_from"] = o.ValidFrom
	}
	if !IsNil(o.ValidUntil) {
		toSerialize["valid_until"] = o.ValidUntil
	}
	if !IsNil(o.Created) {
		toSerialize["created"] = o.Created
	}
	if !IsNil(o.Updated) {
		toSerialize["updated"] = o.Updated
	}
	return toSerialize, nil
}

func (o *GetVulnPolicyResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"uuid",
		"name",
		"condition",
		"analysis",
		"operation_mode",
		"priority",
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

	varGetVulnPolicyResponse := _GetVulnPolicyResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varGetVulnPolicyResponse)

	if err != nil {
		return err
	}

	*o = GetVulnPolicyResponse(varGetVulnPolicyResponse)

	return err
}

type NullableGetVulnPolicyResponse struct {
	value *GetVulnPolicyResponse
	isSet bool
}

func (v NullableGetVulnPolicyResponse) Get() *GetVulnPolicyResponse {
	return v.value
}

func (v *NullableGetVulnPolicyResponse) Set(val *GetVulnPolicyResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetVulnPolicyResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetVulnPolicyResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetVulnPolicyResponse(val *GetVulnPolicyResponse) *NullableGetVulnPolicyResponse {
	return &NullableGetVulnPolicyResponse{value: val, isSet: true}
}

func (v NullableGetVulnPolicyResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetVulnPolicyResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
