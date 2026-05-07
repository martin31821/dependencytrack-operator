# BomUploadResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectUuid** | **string** | UUID of the project the BOM was uploaded for | 
**Token** | **string** | Token used to check task progress | 

## Methods

### NewBomUploadResponse

`func NewBomUploadResponse(projectUuid string, token string, ) *BomUploadResponse`

NewBomUploadResponse instantiates a new BomUploadResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBomUploadResponseWithDefaults

`func NewBomUploadResponseWithDefaults() *BomUploadResponse`

NewBomUploadResponseWithDefaults instantiates a new BomUploadResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProjectUuid

`func (o *BomUploadResponse) GetProjectUuid() string`

GetProjectUuid returns the ProjectUuid field if non-nil, zero value otherwise.

### GetProjectUuidOk

`func (o *BomUploadResponse) GetProjectUuidOk() (*string, bool)`

GetProjectUuidOk returns a tuple with the ProjectUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectUuid

`func (o *BomUploadResponse) SetProjectUuid(v string)`

SetProjectUuid sets ProjectUuid field to given value.


### GetToken

`func (o *BomUploadResponse) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *BomUploadResponse) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *BomUploadResponse) SetToken(v string)`

SetToken sets Token field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


