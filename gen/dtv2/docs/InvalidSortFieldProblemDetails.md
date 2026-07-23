# InvalidSortFieldProblemDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InvalidField** | **string** | Name of the field for which sorting is not supported. | 
**SupportedFields** | Pointer to **[]string** | Names of fields for which sorting is supported. When empty, sorting is explicitly *not* supported. When absent, sorting may be supported, but no definitive guarantees exist. Consult the operation&#39;s description. | [optional] 

## Methods

### NewInvalidSortFieldProblemDetails

`func NewInvalidSortFieldProblemDetails(invalidField string, ) *InvalidSortFieldProblemDetails`

NewInvalidSortFieldProblemDetails instantiates a new InvalidSortFieldProblemDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInvalidSortFieldProblemDetailsWithDefaults

`func NewInvalidSortFieldProblemDetailsWithDefaults() *InvalidSortFieldProblemDetails`

NewInvalidSortFieldProblemDetailsWithDefaults instantiates a new InvalidSortFieldProblemDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInvalidField

`func (o *InvalidSortFieldProblemDetails) GetInvalidField() string`

GetInvalidField returns the InvalidField field if non-nil, zero value otherwise.

### GetInvalidFieldOk

`func (o *InvalidSortFieldProblemDetails) GetInvalidFieldOk() (*string, bool)`

GetInvalidFieldOk returns a tuple with the InvalidField field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvalidField

`func (o *InvalidSortFieldProblemDetails) SetInvalidField(v string)`

SetInvalidField sets InvalidField field to given value.


### GetSupportedFields

`func (o *InvalidSortFieldProblemDetails) GetSupportedFields() []string`

GetSupportedFields returns the SupportedFields field if non-nil, zero value otherwise.

### GetSupportedFieldsOk

`func (o *InvalidSortFieldProblemDetails) GetSupportedFieldsOk() (*[]string, bool)`

GetSupportedFieldsOk returns a tuple with the SupportedFields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupportedFields

`func (o *InvalidSortFieldProblemDetails) SetSupportedFields(v []string)`

SetSupportedFields sets SupportedFields field to given value.

### HasSupportedFields

`func (o *InvalidSortFieldProblemDetails) HasSupportedFields() bool`

HasSupportedFields returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


