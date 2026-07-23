# WorkflowRunMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**ParentId** | Pointer to **string** |  | [optional] 
**WorkflowName** | **string** |  | 
**WorkflowVersion** | **int32** |  | 
**WorkflowInstanceId** | Pointer to **string** |  | [optional] 
**TaskQueueName** | **string** |  | 
**Status** | [**WorkflowRunStatus**](WorkflowRunStatus.md) |  | 
**Priority** | **int32** |  | 
**ConcurrencyKey** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**CreatedAt** | **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | 
**UpdatedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**StartedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**CompletedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 

## Methods

### NewWorkflowRunMetadata

`func NewWorkflowRunMetadata(id string, workflowName string, workflowVersion int32, taskQueueName string, status WorkflowRunStatus, priority int32, createdAt int64, ) *WorkflowRunMetadata`

NewWorkflowRunMetadata instantiates a new WorkflowRunMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowRunMetadataWithDefaults

`func NewWorkflowRunMetadataWithDefaults() *WorkflowRunMetadata`

NewWorkflowRunMetadataWithDefaults instantiates a new WorkflowRunMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *WorkflowRunMetadata) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WorkflowRunMetadata) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WorkflowRunMetadata) SetId(v string)`

SetId sets Id field to given value.


### GetParentId

`func (o *WorkflowRunMetadata) GetParentId() string`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *WorkflowRunMetadata) GetParentIdOk() (*string, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *WorkflowRunMetadata) SetParentId(v string)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *WorkflowRunMetadata) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetWorkflowName

`func (o *WorkflowRunMetadata) GetWorkflowName() string`

GetWorkflowName returns the WorkflowName field if non-nil, zero value otherwise.

### GetWorkflowNameOk

`func (o *WorkflowRunMetadata) GetWorkflowNameOk() (*string, bool)`

GetWorkflowNameOk returns a tuple with the WorkflowName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowName

`func (o *WorkflowRunMetadata) SetWorkflowName(v string)`

SetWorkflowName sets WorkflowName field to given value.


### GetWorkflowVersion

`func (o *WorkflowRunMetadata) GetWorkflowVersion() int32`

GetWorkflowVersion returns the WorkflowVersion field if non-nil, zero value otherwise.

### GetWorkflowVersionOk

`func (o *WorkflowRunMetadata) GetWorkflowVersionOk() (*int32, bool)`

GetWorkflowVersionOk returns a tuple with the WorkflowVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowVersion

`func (o *WorkflowRunMetadata) SetWorkflowVersion(v int32)`

SetWorkflowVersion sets WorkflowVersion field to given value.


### GetWorkflowInstanceId

`func (o *WorkflowRunMetadata) GetWorkflowInstanceId() string`

GetWorkflowInstanceId returns the WorkflowInstanceId field if non-nil, zero value otherwise.

### GetWorkflowInstanceIdOk

`func (o *WorkflowRunMetadata) GetWorkflowInstanceIdOk() (*string, bool)`

GetWorkflowInstanceIdOk returns a tuple with the WorkflowInstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowInstanceId

`func (o *WorkflowRunMetadata) SetWorkflowInstanceId(v string)`

SetWorkflowInstanceId sets WorkflowInstanceId field to given value.

### HasWorkflowInstanceId

`func (o *WorkflowRunMetadata) HasWorkflowInstanceId() bool`

HasWorkflowInstanceId returns a boolean if a field has been set.

### GetTaskQueueName

`func (o *WorkflowRunMetadata) GetTaskQueueName() string`

GetTaskQueueName returns the TaskQueueName field if non-nil, zero value otherwise.

### GetTaskQueueNameOk

`func (o *WorkflowRunMetadata) GetTaskQueueNameOk() (*string, bool)`

GetTaskQueueNameOk returns a tuple with the TaskQueueName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaskQueueName

`func (o *WorkflowRunMetadata) SetTaskQueueName(v string)`

SetTaskQueueName sets TaskQueueName field to given value.


### GetStatus

`func (o *WorkflowRunMetadata) GetStatus() WorkflowRunStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WorkflowRunMetadata) GetStatusOk() (*WorkflowRunStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WorkflowRunMetadata) SetStatus(v WorkflowRunStatus)`

SetStatus sets Status field to given value.


### GetPriority

`func (o *WorkflowRunMetadata) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *WorkflowRunMetadata) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *WorkflowRunMetadata) SetPriority(v int32)`

SetPriority sets Priority field to given value.


### GetConcurrencyKey

`func (o *WorkflowRunMetadata) GetConcurrencyKey() string`

GetConcurrencyKey returns the ConcurrencyKey field if non-nil, zero value otherwise.

### GetConcurrencyKeyOk

`func (o *WorkflowRunMetadata) GetConcurrencyKeyOk() (*string, bool)`

GetConcurrencyKeyOk returns a tuple with the ConcurrencyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConcurrencyKey

`func (o *WorkflowRunMetadata) SetConcurrencyKey(v string)`

SetConcurrencyKey sets ConcurrencyKey field to given value.

### HasConcurrencyKey

`func (o *WorkflowRunMetadata) HasConcurrencyKey() bool`

HasConcurrencyKey returns a boolean if a field has been set.

### GetLabels

`func (o *WorkflowRunMetadata) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *WorkflowRunMetadata) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *WorkflowRunMetadata) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *WorkflowRunMetadata) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetCreatedAt

`func (o *WorkflowRunMetadata) GetCreatedAt() int64`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *WorkflowRunMetadata) GetCreatedAtOk() (*int64, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *WorkflowRunMetadata) SetCreatedAt(v int64)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *WorkflowRunMetadata) GetUpdatedAt() int64`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *WorkflowRunMetadata) GetUpdatedAtOk() (*int64, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *WorkflowRunMetadata) SetUpdatedAt(v int64)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *WorkflowRunMetadata) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetStartedAt

`func (o *WorkflowRunMetadata) GetStartedAt() int64`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *WorkflowRunMetadata) GetStartedAtOk() (*int64, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *WorkflowRunMetadata) SetStartedAt(v int64)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *WorkflowRunMetadata) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetCompletedAt

`func (o *WorkflowRunMetadata) GetCompletedAt() int64`

GetCompletedAt returns the CompletedAt field if non-nil, zero value otherwise.

### GetCompletedAtOk

`func (o *WorkflowRunMetadata) GetCompletedAtOk() (*int64, bool)`

GetCompletedAtOk returns a tuple with the CompletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedAt

`func (o *WorkflowRunMetadata) SetCompletedAt(v int64)`

SetCompletedAt sets CompletedAt field to given value.

### HasCompletedAt

`func (o *WorkflowRunMetadata) HasCompletedAt() bool`

HasCompletedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


