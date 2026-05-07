# CreateScheduledNotificationRuleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Level** | **string** |  | 
**Name** | **string** |  | 
**Publisher** | [**Publisher**](Publisher.md) |  | 
**Scope** | **string** |  | 

## Methods

### NewCreateScheduledNotificationRuleRequest

`func NewCreateScheduledNotificationRuleRequest(level string, name string, publisher Publisher, scope string, ) *CreateScheduledNotificationRuleRequest`

NewCreateScheduledNotificationRuleRequest instantiates a new CreateScheduledNotificationRuleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateScheduledNotificationRuleRequestWithDefaults

`func NewCreateScheduledNotificationRuleRequestWithDefaults() *CreateScheduledNotificationRuleRequest`

NewCreateScheduledNotificationRuleRequestWithDefaults instantiates a new CreateScheduledNotificationRuleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLevel

`func (o *CreateScheduledNotificationRuleRequest) GetLevel() string`

GetLevel returns the Level field if non-nil, zero value otherwise.

### GetLevelOk

`func (o *CreateScheduledNotificationRuleRequest) GetLevelOk() (*string, bool)`

GetLevelOk returns a tuple with the Level field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLevel

`func (o *CreateScheduledNotificationRuleRequest) SetLevel(v string)`

SetLevel sets Level field to given value.


### GetName

`func (o *CreateScheduledNotificationRuleRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateScheduledNotificationRuleRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateScheduledNotificationRuleRequest) SetName(v string)`

SetName sets Name field to given value.


### GetPublisher

`func (o *CreateScheduledNotificationRuleRequest) GetPublisher() Publisher`

GetPublisher returns the Publisher field if non-nil, zero value otherwise.

### GetPublisherOk

`func (o *CreateScheduledNotificationRuleRequest) GetPublisherOk() (*Publisher, bool)`

GetPublisherOk returns a tuple with the Publisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublisher

`func (o *CreateScheduledNotificationRuleRequest) SetPublisher(v Publisher)`

SetPublisher sets Publisher field to given value.


### GetScope

`func (o *CreateScheduledNotificationRuleRequest) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *CreateScheduledNotificationRuleRequest) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *CreateScheduledNotificationRuleRequest) SetScope(v string)`

SetScope sets Scope field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


