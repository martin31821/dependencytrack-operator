# CreateNotificationRuleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Level** | **string** |  | 
**Name** | **string** |  | 
**Publisher** | [**Publisher**](Publisher.md) |  | 
**Scope** | **string** |  | 

## Methods

### NewCreateNotificationRuleRequest

`func NewCreateNotificationRuleRequest(level string, name string, publisher Publisher, scope string, ) *CreateNotificationRuleRequest`

NewCreateNotificationRuleRequest instantiates a new CreateNotificationRuleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateNotificationRuleRequestWithDefaults

`func NewCreateNotificationRuleRequestWithDefaults() *CreateNotificationRuleRequest`

NewCreateNotificationRuleRequestWithDefaults instantiates a new CreateNotificationRuleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLevel

`func (o *CreateNotificationRuleRequest) GetLevel() string`

GetLevel returns the Level field if non-nil, zero value otherwise.

### GetLevelOk

`func (o *CreateNotificationRuleRequest) GetLevelOk() (*string, bool)`

GetLevelOk returns a tuple with the Level field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLevel

`func (o *CreateNotificationRuleRequest) SetLevel(v string)`

SetLevel sets Level field to given value.


### GetName

`func (o *CreateNotificationRuleRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateNotificationRuleRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateNotificationRuleRequest) SetName(v string)`

SetName sets Name field to given value.


### GetPublisher

`func (o *CreateNotificationRuleRequest) GetPublisher() Publisher`

GetPublisher returns the Publisher field if non-nil, zero value otherwise.

### GetPublisherOk

`func (o *CreateNotificationRuleRequest) GetPublisherOk() (*Publisher, bool)`

GetPublisherOk returns a tuple with the Publisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublisher

`func (o *CreateNotificationRuleRequest) SetPublisher(v Publisher)`

SetPublisher sets Publisher field to given value.


### GetScope

`func (o *CreateNotificationRuleRequest) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *CreateNotificationRuleRequest) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *CreateNotificationRuleRequest) SetScope(v string)`

SetScope sets Scope field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


