# UpdateNotificationRuleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**FilterExpression** | Pointer to **string** |  | [optional] 
**Level** | **string** |  | 
**LogSuccessfulPublish** | Pointer to **bool** |  | [optional] 
**Name** | **string** |  | 
**NotifyChildren** | Pointer to **bool** |  | [optional] 
**NotifyOn** | Pointer to **[]string** |  | [optional] 
**PublisherConfig** | Pointer to **string** |  | [optional] 
**ScheduleCron** | Pointer to **string** |  | [optional] 
**ScheduleSkipUnchanged** | Pointer to **bool** |  | [optional] 
**Scope** | **string** |  | 
**Tags** | Pointer to [**[]Tag**](Tag.md) |  | [optional] 
**Uuid** | **string** |  | 

## Methods

### NewUpdateNotificationRuleRequest

`func NewUpdateNotificationRuleRequest(level string, name string, scope string, uuid string, ) *UpdateNotificationRuleRequest`

NewUpdateNotificationRuleRequest instantiates a new UpdateNotificationRuleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateNotificationRuleRequestWithDefaults

`func NewUpdateNotificationRuleRequestWithDefaults() *UpdateNotificationRuleRequest`

NewUpdateNotificationRuleRequestWithDefaults instantiates a new UpdateNotificationRuleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *UpdateNotificationRuleRequest) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *UpdateNotificationRuleRequest) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *UpdateNotificationRuleRequest) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *UpdateNotificationRuleRequest) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetFilterExpression

`func (o *UpdateNotificationRuleRequest) GetFilterExpression() string`

GetFilterExpression returns the FilterExpression field if non-nil, zero value otherwise.

### GetFilterExpressionOk

`func (o *UpdateNotificationRuleRequest) GetFilterExpressionOk() (*string, bool)`

GetFilterExpressionOk returns a tuple with the FilterExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilterExpression

`func (o *UpdateNotificationRuleRequest) SetFilterExpression(v string)`

SetFilterExpression sets FilterExpression field to given value.

### HasFilterExpression

`func (o *UpdateNotificationRuleRequest) HasFilterExpression() bool`

HasFilterExpression returns a boolean if a field has been set.

### GetLevel

`func (o *UpdateNotificationRuleRequest) GetLevel() string`

GetLevel returns the Level field if non-nil, zero value otherwise.

### GetLevelOk

`func (o *UpdateNotificationRuleRequest) GetLevelOk() (*string, bool)`

GetLevelOk returns a tuple with the Level field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLevel

`func (o *UpdateNotificationRuleRequest) SetLevel(v string)`

SetLevel sets Level field to given value.


### GetLogSuccessfulPublish

`func (o *UpdateNotificationRuleRequest) GetLogSuccessfulPublish() bool`

GetLogSuccessfulPublish returns the LogSuccessfulPublish field if non-nil, zero value otherwise.

### GetLogSuccessfulPublishOk

`func (o *UpdateNotificationRuleRequest) GetLogSuccessfulPublishOk() (*bool, bool)`

GetLogSuccessfulPublishOk returns a tuple with the LogSuccessfulPublish field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogSuccessfulPublish

`func (o *UpdateNotificationRuleRequest) SetLogSuccessfulPublish(v bool)`

SetLogSuccessfulPublish sets LogSuccessfulPublish field to given value.

### HasLogSuccessfulPublish

`func (o *UpdateNotificationRuleRequest) HasLogSuccessfulPublish() bool`

HasLogSuccessfulPublish returns a boolean if a field has been set.

### GetName

`func (o *UpdateNotificationRuleRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateNotificationRuleRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateNotificationRuleRequest) SetName(v string)`

SetName sets Name field to given value.


### GetNotifyChildren

`func (o *UpdateNotificationRuleRequest) GetNotifyChildren() bool`

GetNotifyChildren returns the NotifyChildren field if non-nil, zero value otherwise.

### GetNotifyChildrenOk

`func (o *UpdateNotificationRuleRequest) GetNotifyChildrenOk() (*bool, bool)`

GetNotifyChildrenOk returns a tuple with the NotifyChildren field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyChildren

`func (o *UpdateNotificationRuleRequest) SetNotifyChildren(v bool)`

SetNotifyChildren sets NotifyChildren field to given value.

### HasNotifyChildren

`func (o *UpdateNotificationRuleRequest) HasNotifyChildren() bool`

HasNotifyChildren returns a boolean if a field has been set.

### GetNotifyOn

`func (o *UpdateNotificationRuleRequest) GetNotifyOn() []string`

GetNotifyOn returns the NotifyOn field if non-nil, zero value otherwise.

### GetNotifyOnOk

`func (o *UpdateNotificationRuleRequest) GetNotifyOnOk() (*[]string, bool)`

GetNotifyOnOk returns a tuple with the NotifyOn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyOn

`func (o *UpdateNotificationRuleRequest) SetNotifyOn(v []string)`

SetNotifyOn sets NotifyOn field to given value.

### HasNotifyOn

`func (o *UpdateNotificationRuleRequest) HasNotifyOn() bool`

HasNotifyOn returns a boolean if a field has been set.

### GetPublisherConfig

`func (o *UpdateNotificationRuleRequest) GetPublisherConfig() string`

GetPublisherConfig returns the PublisherConfig field if non-nil, zero value otherwise.

### GetPublisherConfigOk

`func (o *UpdateNotificationRuleRequest) GetPublisherConfigOk() (*string, bool)`

GetPublisherConfigOk returns a tuple with the PublisherConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublisherConfig

`func (o *UpdateNotificationRuleRequest) SetPublisherConfig(v string)`

SetPublisherConfig sets PublisherConfig field to given value.

### HasPublisherConfig

`func (o *UpdateNotificationRuleRequest) HasPublisherConfig() bool`

HasPublisherConfig returns a boolean if a field has been set.

### GetScheduleCron

`func (o *UpdateNotificationRuleRequest) GetScheduleCron() string`

GetScheduleCron returns the ScheduleCron field if non-nil, zero value otherwise.

### GetScheduleCronOk

`func (o *UpdateNotificationRuleRequest) GetScheduleCronOk() (*string, bool)`

GetScheduleCronOk returns a tuple with the ScheduleCron field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduleCron

`func (o *UpdateNotificationRuleRequest) SetScheduleCron(v string)`

SetScheduleCron sets ScheduleCron field to given value.

### HasScheduleCron

`func (o *UpdateNotificationRuleRequest) HasScheduleCron() bool`

HasScheduleCron returns a boolean if a field has been set.

### GetScheduleSkipUnchanged

`func (o *UpdateNotificationRuleRequest) GetScheduleSkipUnchanged() bool`

GetScheduleSkipUnchanged returns the ScheduleSkipUnchanged field if non-nil, zero value otherwise.

### GetScheduleSkipUnchangedOk

`func (o *UpdateNotificationRuleRequest) GetScheduleSkipUnchangedOk() (*bool, bool)`

GetScheduleSkipUnchangedOk returns a tuple with the ScheduleSkipUnchanged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduleSkipUnchanged

`func (o *UpdateNotificationRuleRequest) SetScheduleSkipUnchanged(v bool)`

SetScheduleSkipUnchanged sets ScheduleSkipUnchanged field to given value.

### HasScheduleSkipUnchanged

`func (o *UpdateNotificationRuleRequest) HasScheduleSkipUnchanged() bool`

HasScheduleSkipUnchanged returns a boolean if a field has been set.

### GetScope

`func (o *UpdateNotificationRuleRequest) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *UpdateNotificationRuleRequest) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *UpdateNotificationRuleRequest) SetScope(v string)`

SetScope sets Scope field to given value.


### GetTags

`func (o *UpdateNotificationRuleRequest) GetTags() []Tag`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *UpdateNotificationRuleRequest) GetTagsOk() (*[]Tag, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *UpdateNotificationRuleRequest) SetTags(v []Tag)`

SetTags sets Tags field to given value.

### HasTags

`func (o *UpdateNotificationRuleRequest) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetUuid

`func (o *UpdateNotificationRuleRequest) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *UpdateNotificationRuleRequest) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *UpdateNotificationRuleRequest) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


