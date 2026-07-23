# ListTaskQueuesResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | [**[]TaskQueue**](TaskQueue.md) |  | 

## Methods

### NewListTaskQueuesResponse

`func NewListTaskQueuesResponse(items []TaskQueue, ) *ListTaskQueuesResponse`

NewListTaskQueuesResponse instantiates a new ListTaskQueuesResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListTaskQueuesResponseWithDefaults

`func NewListTaskQueuesResponseWithDefaults() *ListTaskQueuesResponse`

NewListTaskQueuesResponseWithDefaults instantiates a new ListTaskQueuesResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *ListTaskQueuesResponse) GetItems() []TaskQueue`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ListTaskQueuesResponse) GetItemsOk() (*[]TaskQueue, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ListTaskQueuesResponse) SetItems(v []TaskQueue)`

SetItems sets Items field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


