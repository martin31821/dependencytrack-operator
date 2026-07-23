# CreateVulnPolicyRequest

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

### NewCreateVulnPolicyRequest

`func NewCreateVulnPolicyRequest(name string, condition string, analysis VulnPolicyAnalysis, ) *CreateVulnPolicyRequest`

NewCreateVulnPolicyRequest instantiates a new CreateVulnPolicyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateVulnPolicyRequestWithDefaults

`func NewCreateVulnPolicyRequestWithDefaults() *CreateVulnPolicyRequest`

NewCreateVulnPolicyRequestWithDefaults instantiates a new CreateVulnPolicyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateVulnPolicyRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateVulnPolicyRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateVulnPolicyRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CreateVulnPolicyRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateVulnPolicyRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateVulnPolicyRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreateVulnPolicyRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAuthor

`func (o *CreateVulnPolicyRequest) GetAuthor() string`

GetAuthor returns the Author field if non-nil, zero value otherwise.

### GetAuthorOk

`func (o *CreateVulnPolicyRequest) GetAuthorOk() (*string, bool)`

GetAuthorOk returns a tuple with the Author field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthor

`func (o *CreateVulnPolicyRequest) SetAuthor(v string)`

SetAuthor sets Author field to given value.

### HasAuthor

`func (o *CreateVulnPolicyRequest) HasAuthor() bool`

HasAuthor returns a boolean if a field has been set.

### GetCondition

`func (o *CreateVulnPolicyRequest) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *CreateVulnPolicyRequest) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *CreateVulnPolicyRequest) SetCondition(v string)`

SetCondition sets Condition field to given value.


### GetAnalysis

`func (o *CreateVulnPolicyRequest) GetAnalysis() VulnPolicyAnalysis`

GetAnalysis returns the Analysis field if non-nil, zero value otherwise.

### GetAnalysisOk

`func (o *CreateVulnPolicyRequest) GetAnalysisOk() (*VulnPolicyAnalysis, bool)`

GetAnalysisOk returns a tuple with the Analysis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysis

`func (o *CreateVulnPolicyRequest) SetAnalysis(v VulnPolicyAnalysis)`

SetAnalysis sets Analysis field to given value.


### GetRatings

`func (o *CreateVulnPolicyRequest) GetRatings() []VulnPolicyRating`

GetRatings returns the Ratings field if non-nil, zero value otherwise.

### GetRatingsOk

`func (o *CreateVulnPolicyRequest) GetRatingsOk() (*[]VulnPolicyRating, bool)`

GetRatingsOk returns a tuple with the Ratings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRatings

`func (o *CreateVulnPolicyRequest) SetRatings(v []VulnPolicyRating)`

SetRatings sets Ratings field to given value.

### HasRatings

`func (o *CreateVulnPolicyRequest) HasRatings() bool`

HasRatings returns a boolean if a field has been set.

### GetOperationMode

`func (o *CreateVulnPolicyRequest) GetOperationMode() VulnPolicyOperationMode`

GetOperationMode returns the OperationMode field if non-nil, zero value otherwise.

### GetOperationModeOk

`func (o *CreateVulnPolicyRequest) GetOperationModeOk() (*VulnPolicyOperationMode, bool)`

GetOperationModeOk returns a tuple with the OperationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperationMode

`func (o *CreateVulnPolicyRequest) SetOperationMode(v VulnPolicyOperationMode)`

SetOperationMode sets OperationMode field to given value.

### HasOperationMode

`func (o *CreateVulnPolicyRequest) HasOperationMode() bool`

HasOperationMode returns a boolean if a field has been set.

### GetPriority

`func (o *CreateVulnPolicyRequest) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *CreateVulnPolicyRequest) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *CreateVulnPolicyRequest) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *CreateVulnPolicyRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetValidFrom

`func (o *CreateVulnPolicyRequest) GetValidFrom() int64`

GetValidFrom returns the ValidFrom field if non-nil, zero value otherwise.

### GetValidFromOk

`func (o *CreateVulnPolicyRequest) GetValidFromOk() (*int64, bool)`

GetValidFromOk returns a tuple with the ValidFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidFrom

`func (o *CreateVulnPolicyRequest) SetValidFrom(v int64)`

SetValidFrom sets ValidFrom field to given value.

### HasValidFrom

`func (o *CreateVulnPolicyRequest) HasValidFrom() bool`

HasValidFrom returns a boolean if a field has been set.

### GetValidUntil

`func (o *CreateVulnPolicyRequest) GetValidUntil() int64`

GetValidUntil returns the ValidUntil field if non-nil, zero value otherwise.

### GetValidUntilOk

`func (o *CreateVulnPolicyRequest) GetValidUntilOk() (*int64, bool)`

GetValidUntilOk returns a tuple with the ValidUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidUntil

`func (o *CreateVulnPolicyRequest) SetValidUntil(v int64)`

SetValidUntil sets ValidUntil field to given value.

### HasValidUntil

`func (o *CreateVulnPolicyRequest) HasValidUntil() bool`

HasValidUntil returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


