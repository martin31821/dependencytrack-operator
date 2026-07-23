# VulnPolicyConditionError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Line** | **int32** | Line number where the error occurred | 
**Column** | **int32** | Column number where the error occurred | 
**Message** | **string** | Description of the error | 

## Methods

### NewVulnPolicyConditionError

`func NewVulnPolicyConditionError(line int32, column int32, message string, ) *VulnPolicyConditionError`

NewVulnPolicyConditionError instantiates a new VulnPolicyConditionError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVulnPolicyConditionErrorWithDefaults

`func NewVulnPolicyConditionErrorWithDefaults() *VulnPolicyConditionError`

NewVulnPolicyConditionErrorWithDefaults instantiates a new VulnPolicyConditionError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLine

`func (o *VulnPolicyConditionError) GetLine() int32`

GetLine returns the Line field if non-nil, zero value otherwise.

### GetLineOk

`func (o *VulnPolicyConditionError) GetLineOk() (*int32, bool)`

GetLineOk returns a tuple with the Line field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLine

`func (o *VulnPolicyConditionError) SetLine(v int32)`

SetLine sets Line field to given value.


### GetColumn

`func (o *VulnPolicyConditionError) GetColumn() int32`

GetColumn returns the Column field if non-nil, zero value otherwise.

### GetColumnOk

`func (o *VulnPolicyConditionError) GetColumnOk() (*int32, bool)`

GetColumnOk returns a tuple with the Column field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumn

`func (o *VulnPolicyConditionError) SetColumn(v int32)`

SetColumn sets Column field to given value.


### GetMessage

`func (o *VulnPolicyConditionError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *VulnPolicyConditionError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *VulnPolicyConditionError) SetMessage(v string)`

SetMessage sets Message field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


