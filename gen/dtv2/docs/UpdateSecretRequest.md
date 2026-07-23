# UpdateSecretRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | The new description. Omit this field to retain the current description. | [optional] 
**Value** | Pointer to **string** | The new value. Omit this field to retain the current value. | [optional] 

## Methods

### NewUpdateSecretRequest

`func NewUpdateSecretRequest() *UpdateSecretRequest`

NewUpdateSecretRequest instantiates a new UpdateSecretRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateSecretRequestWithDefaults

`func NewUpdateSecretRequestWithDefaults() *UpdateSecretRequest`

NewUpdateSecretRequestWithDefaults instantiates a new UpdateSecretRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *UpdateSecretRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdateSecretRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdateSecretRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdateSecretRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetValue

`func (o *UpdateSecretRequest) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *UpdateSecretRequest) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *UpdateSecretRequest) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *UpdateSecretRequest) HasValue() bool`

HasValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


