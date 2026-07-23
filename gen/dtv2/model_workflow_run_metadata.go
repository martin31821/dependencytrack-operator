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

// checks if the WorkflowRunMetadata type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkflowRunMetadata{}

// WorkflowRunMetadata struct for WorkflowRunMetadata
type WorkflowRunMetadata struct {
	Id                 string             `json:"id"`
	ParentId           *string            `json:"parent_id,omitempty"`
	WorkflowName       string             `json:"workflow_name"`
	WorkflowVersion    int32              `json:"workflow_version"`
	WorkflowInstanceId *string            `json:"workflow_instance_id,omitempty"`
	TaskQueueName      string             `json:"task_queue_name"`
	Status             WorkflowRunStatus  `json:"status"`
	Priority           int32              `json:"priority"`
	ConcurrencyKey     *string            `json:"concurrency_key,omitempty"`
	Labels             *map[string]string `json:"labels,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	CreatedAt int64 `json:"created_at"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	UpdatedAt *int64 `json:"updated_at,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	StartedAt *int64 `json:"started_at,omitempty"`
	// Epoch timestamp in milliseconds since January 1, 1970 UTC.
	CompletedAt *int64 `json:"completed_at,omitempty"`
}

type _WorkflowRunMetadata WorkflowRunMetadata

// NewWorkflowRunMetadata instantiates a new WorkflowRunMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowRunMetadata(id string, workflowName string, workflowVersion int32, taskQueueName string, status WorkflowRunStatus, priority int32, createdAt int64) *WorkflowRunMetadata {
	this := WorkflowRunMetadata{}
	this.Id = id
	this.WorkflowName = workflowName
	this.WorkflowVersion = workflowVersion
	this.TaskQueueName = taskQueueName
	this.Status = status
	this.Priority = priority
	this.CreatedAt = createdAt
	return &this
}

// NewWorkflowRunMetadataWithDefaults instantiates a new WorkflowRunMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowRunMetadataWithDefaults() *WorkflowRunMetadata {
	this := WorkflowRunMetadata{}
	return &this
}

// GetId returns the Id field value
func (o *WorkflowRunMetadata) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkflowRunMetadata) SetId(v string) {
	o.Id = v
}

// GetParentId returns the ParentId field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetParentId() string {
	if o == nil || IsNil(o.ParentId) {
		var ret string
		return ret
	}
	return *o.ParentId
}

// GetParentIdOk returns a tuple with the ParentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetParentIdOk() (*string, bool) {
	if o == nil || IsNil(o.ParentId) {
		return nil, false
	}
	return o.ParentId, true
}

// HasParentId returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasParentId() bool {
	if o != nil && !IsNil(o.ParentId) {
		return true
	}

	return false
}

// SetParentId gets a reference to the given string and assigns it to the ParentId field.
func (o *WorkflowRunMetadata) SetParentId(v string) {
	o.ParentId = &v
}

// GetWorkflowName returns the WorkflowName field value
func (o *WorkflowRunMetadata) GetWorkflowName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowName
}

// GetWorkflowNameOk returns a tuple with the WorkflowName field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetWorkflowNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowName, true
}

// SetWorkflowName sets field value
func (o *WorkflowRunMetadata) SetWorkflowName(v string) {
	o.WorkflowName = v
}

// GetWorkflowVersion returns the WorkflowVersion field value
func (o *WorkflowRunMetadata) GetWorkflowVersion() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.WorkflowVersion
}

// GetWorkflowVersionOk returns a tuple with the WorkflowVersion field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetWorkflowVersionOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowVersion, true
}

// SetWorkflowVersion sets field value
func (o *WorkflowRunMetadata) SetWorkflowVersion(v int32) {
	o.WorkflowVersion = v
}

// GetWorkflowInstanceId returns the WorkflowInstanceId field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetWorkflowInstanceId() string {
	if o == nil || IsNil(o.WorkflowInstanceId) {
		var ret string
		return ret
	}
	return *o.WorkflowInstanceId
}

// GetWorkflowInstanceIdOk returns a tuple with the WorkflowInstanceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetWorkflowInstanceIdOk() (*string, bool) {
	if o == nil || IsNil(o.WorkflowInstanceId) {
		return nil, false
	}
	return o.WorkflowInstanceId, true
}

// HasWorkflowInstanceId returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasWorkflowInstanceId() bool {
	if o != nil && !IsNil(o.WorkflowInstanceId) {
		return true
	}

	return false
}

// SetWorkflowInstanceId gets a reference to the given string and assigns it to the WorkflowInstanceId field.
func (o *WorkflowRunMetadata) SetWorkflowInstanceId(v string) {
	o.WorkflowInstanceId = &v
}

// GetTaskQueueName returns the TaskQueueName field value
func (o *WorkflowRunMetadata) GetTaskQueueName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TaskQueueName
}

// GetTaskQueueNameOk returns a tuple with the TaskQueueName field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetTaskQueueNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TaskQueueName, true
}

// SetTaskQueueName sets field value
func (o *WorkflowRunMetadata) SetTaskQueueName(v string) {
	o.TaskQueueName = v
}

// GetStatus returns the Status field value
func (o *WorkflowRunMetadata) GetStatus() WorkflowRunStatus {
	if o == nil {
		var ret WorkflowRunStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetStatusOk() (*WorkflowRunStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *WorkflowRunMetadata) SetStatus(v WorkflowRunStatus) {
	o.Status = v
}

// GetPriority returns the Priority field value
func (o *WorkflowRunMetadata) GetPriority() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetPriorityOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Priority, true
}

// SetPriority sets field value
func (o *WorkflowRunMetadata) SetPriority(v int32) {
	o.Priority = v
}

// GetConcurrencyKey returns the ConcurrencyKey field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetConcurrencyKey() string {
	if o == nil || IsNil(o.ConcurrencyKey) {
		var ret string
		return ret
	}
	return *o.ConcurrencyKey
}

// GetConcurrencyKeyOk returns a tuple with the ConcurrencyKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetConcurrencyKeyOk() (*string, bool) {
	if o == nil || IsNil(o.ConcurrencyKey) {
		return nil, false
	}
	return o.ConcurrencyKey, true
}

// HasConcurrencyKey returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasConcurrencyKey() bool {
	if o != nil && !IsNil(o.ConcurrencyKey) {
		return true
	}

	return false
}

// SetConcurrencyKey gets a reference to the given string and assigns it to the ConcurrencyKey field.
func (o *WorkflowRunMetadata) SetConcurrencyKey(v string) {
	o.ConcurrencyKey = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *WorkflowRunMetadata) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetCreatedAt returns the CreatedAt field value
func (o *WorkflowRunMetadata) GetCreatedAt() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetCreatedAtOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *WorkflowRunMetadata) SetCreatedAt(v int64) {
	o.CreatedAt = v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetUpdatedAt() int64 {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret int64
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetUpdatedAtOk() (*int64, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given int64 and assigns it to the UpdatedAt field.
func (o *WorkflowRunMetadata) SetUpdatedAt(v int64) {
	o.UpdatedAt = &v
}

// GetStartedAt returns the StartedAt field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetStartedAt() int64 {
	if o == nil || IsNil(o.StartedAt) {
		var ret int64
		return ret
	}
	return *o.StartedAt
}

// GetStartedAtOk returns a tuple with the StartedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetStartedAtOk() (*int64, bool) {
	if o == nil || IsNil(o.StartedAt) {
		return nil, false
	}
	return o.StartedAt, true
}

// HasStartedAt returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasStartedAt() bool {
	if o != nil && !IsNil(o.StartedAt) {
		return true
	}

	return false
}

// SetStartedAt gets a reference to the given int64 and assigns it to the StartedAt field.
func (o *WorkflowRunMetadata) SetStartedAt(v int64) {
	o.StartedAt = &v
}

// GetCompletedAt returns the CompletedAt field value if set, zero value otherwise.
func (o *WorkflowRunMetadata) GetCompletedAt() int64 {
	if o == nil || IsNil(o.CompletedAt) {
		var ret int64
		return ret
	}
	return *o.CompletedAt
}

// GetCompletedAtOk returns a tuple with the CompletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowRunMetadata) GetCompletedAtOk() (*int64, bool) {
	if o == nil || IsNil(o.CompletedAt) {
		return nil, false
	}
	return o.CompletedAt, true
}

// HasCompletedAt returns a boolean if a field has been set.
func (o *WorkflowRunMetadata) HasCompletedAt() bool {
	if o != nil && !IsNil(o.CompletedAt) {
		return true
	}

	return false
}

// SetCompletedAt gets a reference to the given int64 and assigns it to the CompletedAt field.
func (o *WorkflowRunMetadata) SetCompletedAt(v int64) {
	o.CompletedAt = &v
}

func (o WorkflowRunMetadata) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkflowRunMetadata) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	if !IsNil(o.ParentId) {
		toSerialize["parent_id"] = o.ParentId
	}
	toSerialize["workflow_name"] = o.WorkflowName
	toSerialize["workflow_version"] = o.WorkflowVersion
	if !IsNil(o.WorkflowInstanceId) {
		toSerialize["workflow_instance_id"] = o.WorkflowInstanceId
	}
	toSerialize["task_queue_name"] = o.TaskQueueName
	toSerialize["status"] = o.Status
	toSerialize["priority"] = o.Priority
	if !IsNil(o.ConcurrencyKey) {
		toSerialize["concurrency_key"] = o.ConcurrencyKey
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["created_at"] = o.CreatedAt
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !IsNil(o.StartedAt) {
		toSerialize["started_at"] = o.StartedAt
	}
	if !IsNil(o.CompletedAt) {
		toSerialize["completed_at"] = o.CompletedAt
	}
	return toSerialize, nil
}

func (o *WorkflowRunMetadata) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"workflow_name",
		"workflow_version",
		"task_queue_name",
		"status",
		"priority",
		"created_at",
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

	varWorkflowRunMetadata := _WorkflowRunMetadata{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varWorkflowRunMetadata)

	if err != nil {
		return err
	}

	*o = WorkflowRunMetadata(varWorkflowRunMetadata)

	return err
}

type NullableWorkflowRunMetadata struct {
	value *WorkflowRunMetadata
	isSet bool
}

func (v NullableWorkflowRunMetadata) Get() *WorkflowRunMetadata {
	return v.value
}

func (v *NullableWorkflowRunMetadata) Set(val *WorkflowRunMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowRunMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowRunMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowRunMetadata(val *WorkflowRunMetadata) *NullableWorkflowRunMetadata {
	return &NullableWorkflowRunMetadata{value: val, isSet: true}
}

func (v NullableWorkflowRunMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowRunMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
