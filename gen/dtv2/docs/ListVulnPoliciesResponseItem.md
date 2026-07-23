# ListVulnPoliciesResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Author** | Pointer to **string** |  | [optional] 
**Priority** | **int32** |  | 
**OperationMode** | [**VulnPolicyOperationMode**](VulnPolicyOperationMode.md) |  | 
**Source** | [**VulnPolicySource**](VulnPolicySource.md) |  | 

## Methods

### NewListVulnPoliciesResponseItem

`func NewListVulnPoliciesResponseItem(uuid string, name string, priority int32, operationMode VulnPolicyOperationMode, source VulnPolicySource, ) *ListVulnPoliciesResponseItem`

NewListVulnPoliciesResponseItem instantiates a new ListVulnPoliciesResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListVulnPoliciesResponseItemWithDefaults

`func NewListVulnPoliciesResponseItemWithDefaults() *ListVulnPoliciesResponseItem`

NewListVulnPoliciesResponseItemWithDefaults instantiates a new ListVulnPoliciesResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *ListVulnPoliciesResponseItem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ListVulnPoliciesResponseItem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ListVulnPoliciesResponseItem) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetName

`func (o *ListVulnPoliciesResponseItem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListVulnPoliciesResponseItem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListVulnPoliciesResponseItem) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *ListVulnPoliciesResponseItem) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ListVulnPoliciesResponseItem) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ListVulnPoliciesResponseItem) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ListVulnPoliciesResponseItem) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAuthor

`func (o *ListVulnPoliciesResponseItem) GetAuthor() string`

GetAuthor returns the Author field if non-nil, zero value otherwise.

### GetAuthorOk

`func (o *ListVulnPoliciesResponseItem) GetAuthorOk() (*string, bool)`

GetAuthorOk returns a tuple with the Author field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthor

`func (o *ListVulnPoliciesResponseItem) SetAuthor(v string)`

SetAuthor sets Author field to given value.

### HasAuthor

`func (o *ListVulnPoliciesResponseItem) HasAuthor() bool`

HasAuthor returns a boolean if a field has been set.

### GetPriority

`func (o *ListVulnPoliciesResponseItem) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ListVulnPoliciesResponseItem) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ListVulnPoliciesResponseItem) SetPriority(v int32)`

SetPriority sets Priority field to given value.


### GetOperationMode

`func (o *ListVulnPoliciesResponseItem) GetOperationMode() VulnPolicyOperationMode`

GetOperationMode returns the OperationMode field if non-nil, zero value otherwise.

### GetOperationModeOk

`func (o *ListVulnPoliciesResponseItem) GetOperationModeOk() (*VulnPolicyOperationMode, bool)`

GetOperationModeOk returns a tuple with the OperationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperationMode

`func (o *ListVulnPoliciesResponseItem) SetOperationMode(v VulnPolicyOperationMode)`

SetOperationMode sets OperationMode field to given value.


### GetSource

`func (o *ListVulnPoliciesResponseItem) GetSource() VulnPolicySource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *ListVulnPoliciesResponseItem) GetSourceOk() (*VulnPolicySource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *ListVulnPoliciesResponseItem) SetSource(v VulnPolicySource)`

SetSource sets Source field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


