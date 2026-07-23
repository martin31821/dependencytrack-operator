# JsonSchemaValidationError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceLocation** | **string** | JSON Pointer to the location in the instance that failed validation | 
**EvaluationPath** | Pointer to **string** | JSON Pointer to the location in the schema during evaluation | [optional] 
**SchemaLocation** | Pointer to **string** | Schema location that generated the error | [optional] 
**Keyword** | Pointer to **string** | The validation keyword that failed | [optional] 
**Message** | **string** | Human-readable error message | 

## Methods

### NewJsonSchemaValidationError

`func NewJsonSchemaValidationError(instanceLocation string, message string, ) *JsonSchemaValidationError`

NewJsonSchemaValidationError instantiates a new JsonSchemaValidationError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJsonSchemaValidationErrorWithDefaults

`func NewJsonSchemaValidationErrorWithDefaults() *JsonSchemaValidationError`

NewJsonSchemaValidationErrorWithDefaults instantiates a new JsonSchemaValidationError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInstanceLocation

`func (o *JsonSchemaValidationError) GetInstanceLocation() string`

GetInstanceLocation returns the InstanceLocation field if non-nil, zero value otherwise.

### GetInstanceLocationOk

`func (o *JsonSchemaValidationError) GetInstanceLocationOk() (*string, bool)`

GetInstanceLocationOk returns a tuple with the InstanceLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceLocation

`func (o *JsonSchemaValidationError) SetInstanceLocation(v string)`

SetInstanceLocation sets InstanceLocation field to given value.


### GetEvaluationPath

`func (o *JsonSchemaValidationError) GetEvaluationPath() string`

GetEvaluationPath returns the EvaluationPath field if non-nil, zero value otherwise.

### GetEvaluationPathOk

`func (o *JsonSchemaValidationError) GetEvaluationPathOk() (*string, bool)`

GetEvaluationPathOk returns a tuple with the EvaluationPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvaluationPath

`func (o *JsonSchemaValidationError) SetEvaluationPath(v string)`

SetEvaluationPath sets EvaluationPath field to given value.

### HasEvaluationPath

`func (o *JsonSchemaValidationError) HasEvaluationPath() bool`

HasEvaluationPath returns a boolean if a field has been set.

### GetSchemaLocation

`func (o *JsonSchemaValidationError) GetSchemaLocation() string`

GetSchemaLocation returns the SchemaLocation field if non-nil, zero value otherwise.

### GetSchemaLocationOk

`func (o *JsonSchemaValidationError) GetSchemaLocationOk() (*string, bool)`

GetSchemaLocationOk returns a tuple with the SchemaLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchemaLocation

`func (o *JsonSchemaValidationError) SetSchemaLocation(v string)`

SetSchemaLocation sets SchemaLocation field to given value.

### HasSchemaLocation

`func (o *JsonSchemaValidationError) HasSchemaLocation() bool`

HasSchemaLocation returns a boolean if a field has been set.

### GetKeyword

`func (o *JsonSchemaValidationError) GetKeyword() string`

GetKeyword returns the Keyword field if non-nil, zero value otherwise.

### GetKeywordOk

`func (o *JsonSchemaValidationError) GetKeywordOk() (*string, bool)`

GetKeywordOk returns a tuple with the Keyword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyword

`func (o *JsonSchemaValidationError) SetKeyword(v string)`

SetKeyword sets Keyword field to given value.

### HasKeyword

`func (o *JsonSchemaValidationError) HasKeyword() bool`

HasKeyword returns a boolean if a field has been set.

### GetMessage

`func (o *JsonSchemaValidationError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *JsonSchemaValidationError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *JsonSchemaValidationError) SetMessage(v string)`

SetMessage sets Message field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


