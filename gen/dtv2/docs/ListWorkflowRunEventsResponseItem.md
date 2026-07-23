# ListWorkflowRunEventsResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SequenceNumber** | **int32** |  | 
**Event** | **map[string]interface{}** |  | 

## Methods

### NewListWorkflowRunEventsResponseItem

`func NewListWorkflowRunEventsResponseItem(sequenceNumber int32, event map[string]interface{}, ) *ListWorkflowRunEventsResponseItem`

NewListWorkflowRunEventsResponseItem instantiates a new ListWorkflowRunEventsResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListWorkflowRunEventsResponseItemWithDefaults

`func NewListWorkflowRunEventsResponseItemWithDefaults() *ListWorkflowRunEventsResponseItem`

NewListWorkflowRunEventsResponseItemWithDefaults instantiates a new ListWorkflowRunEventsResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSequenceNumber

`func (o *ListWorkflowRunEventsResponseItem) GetSequenceNumber() int32`

GetSequenceNumber returns the SequenceNumber field if non-nil, zero value otherwise.

### GetSequenceNumberOk

`func (o *ListWorkflowRunEventsResponseItem) GetSequenceNumberOk() (*int32, bool)`

GetSequenceNumberOk returns a tuple with the SequenceNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSequenceNumber

`func (o *ListWorkflowRunEventsResponseItem) SetSequenceNumber(v int32)`

SetSequenceNumber sets SequenceNumber field to given value.


### GetEvent

`func (o *ListWorkflowRunEventsResponseItem) GetEvent() map[string]interface{}`

GetEvent returns the Event field if non-nil, zero value otherwise.

### GetEventOk

`func (o *ListWorkflowRunEventsResponseItem) GetEventOk() (*map[string]interface{}, bool)`

GetEventOk returns a tuple with the Event field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvent

`func (o *ListWorkflowRunEventsResponseItem) SetEvent(v map[string]interface{})`

SetEvent sets Event field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


