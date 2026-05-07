# NotificationPublisher

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultPublisher** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**ExtensionName** | **string** |  | 
**Name** | **string** |  | 
**Template** | Pointer to **string** |  | [optional] 
**TemplateMimeType** | **string** |  | 
**Uuid** | **string** |  | 

## Methods

### NewNotificationPublisher

`func NewNotificationPublisher(extensionName string, name string, templateMimeType string, uuid string, ) *NotificationPublisher`

NewNotificationPublisher instantiates a new NotificationPublisher object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNotificationPublisherWithDefaults

`func NewNotificationPublisherWithDefaults() *NotificationPublisher`

NewNotificationPublisherWithDefaults instantiates a new NotificationPublisher object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefaultPublisher

`func (o *NotificationPublisher) GetDefaultPublisher() bool`

GetDefaultPublisher returns the DefaultPublisher field if non-nil, zero value otherwise.

### GetDefaultPublisherOk

`func (o *NotificationPublisher) GetDefaultPublisherOk() (*bool, bool)`

GetDefaultPublisherOk returns a tuple with the DefaultPublisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPublisher

`func (o *NotificationPublisher) SetDefaultPublisher(v bool)`

SetDefaultPublisher sets DefaultPublisher field to given value.

### HasDefaultPublisher

`func (o *NotificationPublisher) HasDefaultPublisher() bool`

HasDefaultPublisher returns a boolean if a field has been set.

### GetDescription

`func (o *NotificationPublisher) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NotificationPublisher) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NotificationPublisher) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NotificationPublisher) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetExtensionName

`func (o *NotificationPublisher) GetExtensionName() string`

GetExtensionName returns the ExtensionName field if non-nil, zero value otherwise.

### GetExtensionNameOk

`func (o *NotificationPublisher) GetExtensionNameOk() (*string, bool)`

GetExtensionNameOk returns a tuple with the ExtensionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtensionName

`func (o *NotificationPublisher) SetExtensionName(v string)`

SetExtensionName sets ExtensionName field to given value.


### GetName

`func (o *NotificationPublisher) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NotificationPublisher) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NotificationPublisher) SetName(v string)`

SetName sets Name field to given value.


### GetTemplate

`func (o *NotificationPublisher) GetTemplate() string`

GetTemplate returns the Template field if non-nil, zero value otherwise.

### GetTemplateOk

`func (o *NotificationPublisher) GetTemplateOk() (*string, bool)`

GetTemplateOk returns a tuple with the Template field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplate

`func (o *NotificationPublisher) SetTemplate(v string)`

SetTemplate sets Template field to given value.

### HasTemplate

`func (o *NotificationPublisher) HasTemplate() bool`

HasTemplate returns a boolean if a field has been set.

### GetTemplateMimeType

`func (o *NotificationPublisher) GetTemplateMimeType() string`

GetTemplateMimeType returns the TemplateMimeType field if non-nil, zero value otherwise.

### GetTemplateMimeTypeOk

`func (o *NotificationPublisher) GetTemplateMimeTypeOk() (*string, bool)`

GetTemplateMimeTypeOk returns a tuple with the TemplateMimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplateMimeType

`func (o *NotificationPublisher) SetTemplateMimeType(v string)`

SetTemplateMimeType sets TemplateMimeType field to given value.


### GetUuid

`func (o *NotificationPublisher) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *NotificationPublisher) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *NotificationPublisher) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


