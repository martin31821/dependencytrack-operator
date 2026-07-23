# ConstraintViolationError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Path** | Pointer to **string** | Path to the invalid field in the request | [optional] 
**Value** | Pointer to **string** | The invalid value | [optional] 
**Message** | **string** | Message explaining the error | 

## Methods

### NewConstraintViolationError

`func NewConstraintViolationError(message string, ) *ConstraintViolationError`

NewConstraintViolationError instantiates a new ConstraintViolationError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConstraintViolationErrorWithDefaults

`func NewConstraintViolationErrorWithDefaults() *ConstraintViolationError`

NewConstraintViolationErrorWithDefaults instantiates a new ConstraintViolationError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPath

`func (o *ConstraintViolationError) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *ConstraintViolationError) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *ConstraintViolationError) SetPath(v string)`

SetPath sets Path field to given value.

### HasPath

`func (o *ConstraintViolationError) HasPath() bool`

HasPath returns a boolean if a field has been set.

### GetValue

`func (o *ConstraintViolationError) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *ConstraintViolationError) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *ConstraintViolationError) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *ConstraintViolationError) HasValue() bool`

HasValue returns a boolean if a field has been set.

### GetMessage

`func (o *ConstraintViolationError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ConstraintViolationError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ConstraintViolationError) SetMessage(v string)`

SetMessage sets Message field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


