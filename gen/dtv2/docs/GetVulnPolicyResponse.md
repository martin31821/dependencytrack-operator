# GetVulnPolicyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Author** | Pointer to **string** |  | [optional] 
**Condition** | **string** |  | 
**Analysis** | [**VulnPolicyAnalysis**](VulnPolicyAnalysis.md) |  | 
**Ratings** | Pointer to [**[]VulnPolicyRating**](VulnPolicyRating.md) |  | [optional] 
**OperationMode** | [**VulnPolicyOperationMode**](VulnPolicyOperationMode.md) |  | 
**Priority** | **int32** |  | 
**Source** | [**VulnPolicySource**](VulnPolicySource.md) |  | 
**ValidFrom** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**ValidUntil** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**Created** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**Updated** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 

## Methods

### NewGetVulnPolicyResponse

`func NewGetVulnPolicyResponse(uuid string, name string, condition string, analysis VulnPolicyAnalysis, operationMode VulnPolicyOperationMode, priority int32, source VulnPolicySource, ) *GetVulnPolicyResponse`

NewGetVulnPolicyResponse instantiates a new GetVulnPolicyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetVulnPolicyResponseWithDefaults

`func NewGetVulnPolicyResponseWithDefaults() *GetVulnPolicyResponse`

NewGetVulnPolicyResponseWithDefaults instantiates a new GetVulnPolicyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *GetVulnPolicyResponse) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *GetVulnPolicyResponse) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *GetVulnPolicyResponse) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetName

`func (o *GetVulnPolicyResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GetVulnPolicyResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GetVulnPolicyResponse) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *GetVulnPolicyResponse) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *GetVulnPolicyResponse) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *GetVulnPolicyResponse) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *GetVulnPolicyResponse) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAuthor

`func (o *GetVulnPolicyResponse) GetAuthor() string`

GetAuthor returns the Author field if non-nil, zero value otherwise.

### GetAuthorOk

`func (o *GetVulnPolicyResponse) GetAuthorOk() (*string, bool)`

GetAuthorOk returns a tuple with the Author field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthor

`func (o *GetVulnPolicyResponse) SetAuthor(v string)`

SetAuthor sets Author field to given value.

### HasAuthor

`func (o *GetVulnPolicyResponse) HasAuthor() bool`

HasAuthor returns a boolean if a field has been set.

### GetCondition

`func (o *GetVulnPolicyResponse) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *GetVulnPolicyResponse) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *GetVulnPolicyResponse) SetCondition(v string)`

SetCondition sets Condition field to given value.


### GetAnalysis

`func (o *GetVulnPolicyResponse) GetAnalysis() VulnPolicyAnalysis`

GetAnalysis returns the Analysis field if non-nil, zero value otherwise.

### GetAnalysisOk

`func (o *GetVulnPolicyResponse) GetAnalysisOk() (*VulnPolicyAnalysis, bool)`

GetAnalysisOk returns a tuple with the Analysis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysis

`func (o *GetVulnPolicyResponse) SetAnalysis(v VulnPolicyAnalysis)`

SetAnalysis sets Analysis field to given value.


### GetRatings

`func (o *GetVulnPolicyResponse) GetRatings() []VulnPolicyRating`

GetRatings returns the Ratings field if non-nil, zero value otherwise.

### GetRatingsOk

`func (o *GetVulnPolicyResponse) GetRatingsOk() (*[]VulnPolicyRating, bool)`

GetRatingsOk returns a tuple with the Ratings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRatings

`func (o *GetVulnPolicyResponse) SetRatings(v []VulnPolicyRating)`

SetRatings sets Ratings field to given value.

### HasRatings

`func (o *GetVulnPolicyResponse) HasRatings() bool`

HasRatings returns a boolean if a field has been set.

### GetOperationMode

`func (o *GetVulnPolicyResponse) GetOperationMode() VulnPolicyOperationMode`

GetOperationMode returns the OperationMode field if non-nil, zero value otherwise.

### GetOperationModeOk

`func (o *GetVulnPolicyResponse) GetOperationModeOk() (*VulnPolicyOperationMode, bool)`

GetOperationModeOk returns a tuple with the OperationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperationMode

`func (o *GetVulnPolicyResponse) SetOperationMode(v VulnPolicyOperationMode)`

SetOperationMode sets OperationMode field to given value.


### GetPriority

`func (o *GetVulnPolicyResponse) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *GetVulnPolicyResponse) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *GetVulnPolicyResponse) SetPriority(v int32)`

SetPriority sets Priority field to given value.


### GetSource

`func (o *GetVulnPolicyResponse) GetSource() VulnPolicySource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *GetVulnPolicyResponse) GetSourceOk() (*VulnPolicySource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *GetVulnPolicyResponse) SetSource(v VulnPolicySource)`

SetSource sets Source field to given value.


### GetValidFrom

`func (o *GetVulnPolicyResponse) GetValidFrom() int64`

GetValidFrom returns the ValidFrom field if non-nil, zero value otherwise.

### GetValidFromOk

`func (o *GetVulnPolicyResponse) GetValidFromOk() (*int64, bool)`

GetValidFromOk returns a tuple with the ValidFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidFrom

`func (o *GetVulnPolicyResponse) SetValidFrom(v int64)`

SetValidFrom sets ValidFrom field to given value.

### HasValidFrom

`func (o *GetVulnPolicyResponse) HasValidFrom() bool`

HasValidFrom returns a boolean if a field has been set.

### GetValidUntil

`func (o *GetVulnPolicyResponse) GetValidUntil() int64`

GetValidUntil returns the ValidUntil field if non-nil, zero value otherwise.

### GetValidUntilOk

`func (o *GetVulnPolicyResponse) GetValidUntilOk() (*int64, bool)`

GetValidUntilOk returns a tuple with the ValidUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidUntil

`func (o *GetVulnPolicyResponse) SetValidUntil(v int64)`

SetValidUntil sets ValidUntil field to given value.

### HasValidUntil

`func (o *GetVulnPolicyResponse) HasValidUntil() bool`

HasValidUntil returns a boolean if a field has been set.

### GetCreated

`func (o *GetVulnPolicyResponse) GetCreated() int64`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *GetVulnPolicyResponse) GetCreatedOk() (*int64, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *GetVulnPolicyResponse) SetCreated(v int64)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *GetVulnPolicyResponse) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *GetVulnPolicyResponse) GetUpdated() int64`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *GetVulnPolicyResponse) GetUpdatedOk() (*int64, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *GetVulnPolicyResponse) SetUpdated(v int64)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *GetVulnPolicyResponse) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


