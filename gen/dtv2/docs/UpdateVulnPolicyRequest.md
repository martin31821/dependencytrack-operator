# UpdateVulnPolicyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Author** | Pointer to **string** |  | [optional] 
**Condition** | **string** |  | 
**Analysis** | [**VulnPolicyAnalysis**](VulnPolicyAnalysis.md) |  | 
**Ratings** | Pointer to [**[]VulnPolicyRating**](VulnPolicyRating.md) |  | [optional] 
**OperationMode** | Pointer to [**VulnPolicyOperationMode**](VulnPolicyOperationMode.md) |  | [optional] 
**Priority** | Pointer to **int32** |  | [optional] [default to 0]
**ValidFrom** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**ValidUntil** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 

## Methods

### NewUpdateVulnPolicyRequest

`func NewUpdateVulnPolicyRequest(name string, condition string, analysis VulnPolicyAnalysis, ) *UpdateVulnPolicyRequest`

NewUpdateVulnPolicyRequest instantiates a new UpdateVulnPolicyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateVulnPolicyRequestWithDefaults

`func NewUpdateVulnPolicyRequestWithDefaults() *UpdateVulnPolicyRequest`

NewUpdateVulnPolicyRequestWithDefaults instantiates a new UpdateVulnPolicyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateVulnPolicyRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateVulnPolicyRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateVulnPolicyRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *UpdateVulnPolicyRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdateVulnPolicyRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdateVulnPolicyRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdateVulnPolicyRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAuthor

`func (o *UpdateVulnPolicyRequest) GetAuthor() string`

GetAuthor returns the Author field if non-nil, zero value otherwise.

### GetAuthorOk

`func (o *UpdateVulnPolicyRequest) GetAuthorOk() (*string, bool)`

GetAuthorOk returns a tuple with the Author field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthor

`func (o *UpdateVulnPolicyRequest) SetAuthor(v string)`

SetAuthor sets Author field to given value.

### HasAuthor

`func (o *UpdateVulnPolicyRequest) HasAuthor() bool`

HasAuthor returns a boolean if a field has been set.

### GetCondition

`func (o *UpdateVulnPolicyRequest) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *UpdateVulnPolicyRequest) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *UpdateVulnPolicyRequest) SetCondition(v string)`

SetCondition sets Condition field to given value.


### GetAnalysis

`func (o *UpdateVulnPolicyRequest) GetAnalysis() VulnPolicyAnalysis`

GetAnalysis returns the Analysis field if non-nil, zero value otherwise.

### GetAnalysisOk

`func (o *UpdateVulnPolicyRequest) GetAnalysisOk() (*VulnPolicyAnalysis, bool)`

GetAnalysisOk returns a tuple with the Analysis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysis

`func (o *UpdateVulnPolicyRequest) SetAnalysis(v VulnPolicyAnalysis)`

SetAnalysis sets Analysis field to given value.


### GetRatings

`func (o *UpdateVulnPolicyRequest) GetRatings() []VulnPolicyRating`

GetRatings returns the Ratings field if non-nil, zero value otherwise.

### GetRatingsOk

`func (o *UpdateVulnPolicyRequest) GetRatingsOk() (*[]VulnPolicyRating, bool)`

GetRatingsOk returns a tuple with the Ratings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRatings

`func (o *UpdateVulnPolicyRequest) SetRatings(v []VulnPolicyRating)`

SetRatings sets Ratings field to given value.

### HasRatings

`func (o *UpdateVulnPolicyRequest) HasRatings() bool`

HasRatings returns a boolean if a field has been set.

### GetOperationMode

`func (o *UpdateVulnPolicyRequest) GetOperationMode() VulnPolicyOperationMode`

GetOperationMode returns the OperationMode field if non-nil, zero value otherwise.

### GetOperationModeOk

`func (o *UpdateVulnPolicyRequest) GetOperationModeOk() (*VulnPolicyOperationMode, bool)`

GetOperationModeOk returns a tuple with the OperationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperationMode

`func (o *UpdateVulnPolicyRequest) SetOperationMode(v VulnPolicyOperationMode)`

SetOperationMode sets OperationMode field to given value.

### HasOperationMode

`func (o *UpdateVulnPolicyRequest) HasOperationMode() bool`

HasOperationMode returns a boolean if a field has been set.

### GetPriority

`func (o *UpdateVulnPolicyRequest) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *UpdateVulnPolicyRequest) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *UpdateVulnPolicyRequest) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *UpdateVulnPolicyRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetValidFrom

`func (o *UpdateVulnPolicyRequest) GetValidFrom() int64`

GetValidFrom returns the ValidFrom field if non-nil, zero value otherwise.

### GetValidFromOk

`func (o *UpdateVulnPolicyRequest) GetValidFromOk() (*int64, bool)`

GetValidFromOk returns a tuple with the ValidFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidFrom

`func (o *UpdateVulnPolicyRequest) SetValidFrom(v int64)`

SetValidFrom sets ValidFrom field to given value.

### HasValidFrom

`func (o *UpdateVulnPolicyRequest) HasValidFrom() bool`

HasValidFrom returns a boolean if a field has been set.

### GetValidUntil

`func (o *UpdateVulnPolicyRequest) GetValidUntil() int64`

GetValidUntil returns the ValidUntil field if non-nil, zero value otherwise.

### GetValidUntilOk

`func (o *UpdateVulnPolicyRequest) GetValidUntilOk() (*int64, bool)`

GetValidUntilOk returns a tuple with the ValidUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidUntil

`func (o *UpdateVulnPolicyRequest) SetValidUntil(v int64)`

SetValidUntil sets ValidUntil field to given value.

### HasValidUntil

`func (o *UpdateVulnPolicyRequest) HasValidUntil() bool`

HasValidUntil returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


