# TotalCount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Count** | **int64** | The total number of records across all pages. Might be an exact count, or a lower bound. Refer to the &#x60;type&#x60; field for the applicable semantics. | 
**Type** | [**TotalCountType**](TotalCountType.md) |  | 

## Methods

### NewTotalCount

`func NewTotalCount(count int64, type_ TotalCountType, ) *TotalCount`

NewTotalCount instantiates a new TotalCount object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTotalCountWithDefaults

`func NewTotalCountWithDefaults() *TotalCount`

NewTotalCountWithDefaults instantiates a new TotalCount object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCount

`func (o *TotalCount) GetCount() int64`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *TotalCount) GetCountOk() (*int64, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *TotalCount) SetCount(v int64)`

SetCount sets Count field to given value.


### GetType

`func (o *TotalCount) GetType() TotalCountType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *TotalCount) GetTypeOk() (*TotalCountType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *TotalCount) SetType(v TotalCountType)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


