# WorkflowState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FailureReason** | Pointer to **string** |  | [optional] 
**Parent** | Pointer to [**WorkflowState**](WorkflowState.md) |  | [optional] 
**StartedAt** | Pointer to **util.DTTime** |  | [optional] 
**Status** | **string** |  | 
**Step** | **string** |  | 
**Token** | **string** |  | 
**UpdatedAt** | Pointer to **util.DTTime** |  | [optional] 

## Methods

### NewWorkflowState

`func NewWorkflowState(status string, step string, token string, ) *WorkflowState`

NewWorkflowState instantiates a new WorkflowState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStateWithDefaults

`func NewWorkflowStateWithDefaults() *WorkflowState`

NewWorkflowStateWithDefaults instantiates a new WorkflowState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFailureReason

`func (o *WorkflowState) GetFailureReason() string`

GetFailureReason returns the FailureReason field if non-nil, zero value otherwise.

### GetFailureReasonOk

`func (o *WorkflowState) GetFailureReasonOk() (*string, bool)`

GetFailureReasonOk returns a tuple with the FailureReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailureReason

`func (o *WorkflowState) SetFailureReason(v string)`

SetFailureReason sets FailureReason field to given value.

### HasFailureReason

`func (o *WorkflowState) HasFailureReason() bool`

HasFailureReason returns a boolean if a field has been set.

### GetParent

`func (o *WorkflowState) GetParent() WorkflowState`

GetParent returns the Parent field if non-nil, zero value otherwise.

### GetParentOk

`func (o *WorkflowState) GetParentOk() (*WorkflowState, bool)`

GetParentOk returns a tuple with the Parent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParent

`func (o *WorkflowState) SetParent(v WorkflowState)`

SetParent sets Parent field to given value.

### HasParent

`func (o *WorkflowState) HasParent() bool`

HasParent returns a boolean if a field has been set.

### GetStartedAt

`func (o *WorkflowState) GetStartedAt() util.DTTime`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *WorkflowState) GetStartedAtOk() (*util.DTTime, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *WorkflowState) SetStartedAt(v util.DTTime)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *WorkflowState) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetStatus

`func (o *WorkflowState) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WorkflowState) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WorkflowState) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetStep

`func (o *WorkflowState) GetStep() string`

GetStep returns the Step field if non-nil, zero value otherwise.

### GetStepOk

`func (o *WorkflowState) GetStepOk() (*string, bool)`

GetStepOk returns a tuple with the Step field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStep

`func (o *WorkflowState) SetStep(v string)`

SetStep sets Step field to given value.


### GetToken

`func (o *WorkflowState) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *WorkflowState) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *WorkflowState) SetToken(v string)`

SetToken sets Token field to given value.


### GetUpdatedAt

`func (o *WorkflowState) GetUpdatedAt() util.DTTime`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *WorkflowState) GetUpdatedAtOk() (*util.DTTime, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *WorkflowState) SetUpdatedAt(v util.DTTime)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *WorkflowState) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


