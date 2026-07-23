# JsonSchemaValidationProblemDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | [**[]JsonSchemaValidationError**](JsonSchemaValidationError.md) | List of JSON Schema validation errors | 

## Methods

### NewJsonSchemaValidationProblemDetails

`func NewJsonSchemaValidationProblemDetails(errors []JsonSchemaValidationError, ) *JsonSchemaValidationProblemDetails`

NewJsonSchemaValidationProblemDetails instantiates a new JsonSchemaValidationProblemDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJsonSchemaValidationProblemDetailsWithDefaults

`func NewJsonSchemaValidationProblemDetailsWithDefaults() *JsonSchemaValidationProblemDetails`

NewJsonSchemaValidationProblemDetailsWithDefaults instantiates a new JsonSchemaValidationProblemDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *JsonSchemaValidationProblemDetails) GetErrors() []JsonSchemaValidationError`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *JsonSchemaValidationProblemDetails) GetErrorsOk() (*[]JsonSchemaValidationError, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *JsonSchemaValidationProblemDetails) SetErrors(v []JsonSchemaValidationError)`

SetErrors sets Errors field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


