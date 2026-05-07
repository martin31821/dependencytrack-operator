# CelExpressionError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Column** | Pointer to **int32** |  | [optional] 
**Line** | Pointer to **int32** |  | [optional] 
**Message** | Pointer to **string** |  | [optional] 

## Methods

### NewCelExpressionError

`func NewCelExpressionError() *CelExpressionError`

NewCelExpressionError instantiates a new CelExpressionError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCelExpressionErrorWithDefaults

`func NewCelExpressionErrorWithDefaults() *CelExpressionError`

NewCelExpressionErrorWithDefaults instantiates a new CelExpressionError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetColumn

`func (o *CelExpressionError) GetColumn() int32`

GetColumn returns the Column field if non-nil, zero value otherwise.

### GetColumnOk

`func (o *CelExpressionError) GetColumnOk() (*int32, bool)`

GetColumnOk returns a tuple with the Column field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumn

`func (o *CelExpressionError) SetColumn(v int32)`

SetColumn sets Column field to given value.

### HasColumn

`func (o *CelExpressionError) HasColumn() bool`

HasColumn returns a boolean if a field has been set.

### GetLine

`func (o *CelExpressionError) GetLine() int32`

GetLine returns the Line field if non-nil, zero value otherwise.

### GetLineOk

`func (o *CelExpressionError) GetLineOk() (*int32, bool)`

GetLineOk returns a tuple with the Line field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLine

`func (o *CelExpressionError) SetLine(v int32)`

SetLine sets Line field to given value.

### HasLine

`func (o *CelExpressionError) HasLine() bool`

HasLine returns a boolean if a field has been set.

### GetMessage

`func (o *CelExpressionError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *CelExpressionError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *CelExpressionError) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *CelExpressionError) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


