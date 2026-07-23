# PaginatedResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NextPageToken** | Pointer to **string** | Token to retrieve the next page. Absent when no more items exist. | [optional] 
**Total** | [**TotalCount**](TotalCount.md) |  | 

## Methods

### NewPaginatedResponse

`func NewPaginatedResponse(total TotalCount, ) *PaginatedResponse`

NewPaginatedResponse instantiates a new PaginatedResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaginatedResponseWithDefaults

`func NewPaginatedResponseWithDefaults() *PaginatedResponse`

NewPaginatedResponseWithDefaults instantiates a new PaginatedResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNextPageToken

`func (o *PaginatedResponse) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *PaginatedResponse) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *PaginatedResponse) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *PaginatedResponse) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.

### GetTotal

`func (o *PaginatedResponse) GetTotal() TotalCount`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *PaginatedResponse) GetTotalOk() (*TotalCount, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *PaginatedResponse) SetTotal(v TotalCount)`

SetTotal sets Total field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


