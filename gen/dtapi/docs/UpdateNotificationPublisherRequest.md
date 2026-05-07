# UpdateNotificationPublisherRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** |  | [optional] 
**ExtensionName** | **string** |  | 
**Name** | **string** |  | 
**Template** | Pointer to **string** |  | [optional] 
**TemplateMimeType** | Pointer to **string** |  | [optional] 
**Uuid** | **string** |  | 

## Methods

### NewUpdateNotificationPublisherRequest

`func NewUpdateNotificationPublisherRequest(extensionName string, name string, uuid string, ) *UpdateNotificationPublisherRequest`

NewUpdateNotificationPublisherRequest instantiates a new UpdateNotificationPublisherRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateNotificationPublisherRequestWithDefaults

`func NewUpdateNotificationPublisherRequestWithDefaults() *UpdateNotificationPublisherRequest`

NewUpdateNotificationPublisherRequestWithDefaults instantiates a new UpdateNotificationPublisherRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *UpdateNotificationPublisherRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdateNotificationPublisherRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdateNotificationPublisherRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdateNotificationPublisherRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetExtensionName

`func (o *UpdateNotificationPublisherRequest) GetExtensionName() string`

GetExtensionName returns the ExtensionName field if non-nil, zero value otherwise.

### GetExtensionNameOk

`func (o *UpdateNotificationPublisherRequest) GetExtensionNameOk() (*string, bool)`

GetExtensionNameOk returns a tuple with the ExtensionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtensionName

`func (o *UpdateNotificationPublisherRequest) SetExtensionName(v string)`

SetExtensionName sets ExtensionName field to given value.


### GetName

`func (o *UpdateNotificationPublisherRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateNotificationPublisherRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateNotificationPublisherRequest) SetName(v string)`

SetName sets Name field to given value.


### GetTemplate

`func (o *UpdateNotificationPublisherRequest) GetTemplate() string`

GetTemplate returns the Template field if non-nil, zero value otherwise.

### GetTemplateOk

`func (o *UpdateNotificationPublisherRequest) GetTemplateOk() (*string, bool)`

GetTemplateOk returns a tuple with the Template field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplate

`func (o *UpdateNotificationPublisherRequest) SetTemplate(v string)`

SetTemplate sets Template field to given value.

### HasTemplate

`func (o *UpdateNotificationPublisherRequest) HasTemplate() bool`

HasTemplate returns a boolean if a field has been set.

### GetTemplateMimeType

`func (o *UpdateNotificationPublisherRequest) GetTemplateMimeType() string`

GetTemplateMimeType returns the TemplateMimeType field if non-nil, zero value otherwise.

### GetTemplateMimeTypeOk

`func (o *UpdateNotificationPublisherRequest) GetTemplateMimeTypeOk() (*string, bool)`

GetTemplateMimeTypeOk returns a tuple with the TemplateMimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplateMimeType

`func (o *UpdateNotificationPublisherRequest) SetTemplateMimeType(v string)`

SetTemplateMimeType sets TemplateMimeType field to given value.

### HasTemplateMimeType

`func (o *UpdateNotificationPublisherRequest) HasTemplateMimeType() bool`

HasTemplateMimeType returns a boolean if a field has been set.

### GetUuid

`func (o *UpdateNotificationPublisherRequest) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *UpdateNotificationPublisherRequest) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *UpdateNotificationPublisherRequest) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


