# UpdateTaskQueueRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to [**TaskQueueStatus**](TaskQueueStatus.md) |  | [optional] 
**Capacity** | Pointer to **int32** |  | [optional] 

## Methods

### NewUpdateTaskQueueRequest

`func NewUpdateTaskQueueRequest() *UpdateTaskQueueRequest`

NewUpdateTaskQueueRequest instantiates a new UpdateTaskQueueRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateTaskQueueRequestWithDefaults

`func NewUpdateTaskQueueRequestWithDefaults() *UpdateTaskQueueRequest`

NewUpdateTaskQueueRequestWithDefaults instantiates a new UpdateTaskQueueRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *UpdateTaskQueueRequest) GetStatus() TaskQueueStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdateTaskQueueRequest) GetStatusOk() (*TaskQueueStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdateTaskQueueRequest) SetStatus(v TaskQueueStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *UpdateTaskQueueRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCapacity

`func (o *UpdateTaskQueueRequest) GetCapacity() int32`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *UpdateTaskQueueRequest) GetCapacityOk() (*int32, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *UpdateTaskQueueRequest) SetCapacity(v int32)`

SetCapacity sets Capacity field to given value.

### HasCapacity

`func (o *UpdateTaskQueueRequest) HasCapacity() bool`

HasCapacity returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


