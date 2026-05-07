# CreateNotificationPublisherRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** |  | [optional] 
**ExtensionName** | **string** |  | 
**Name** | **string** |  | 
**Template** | Pointer to **string** |  | [optional] 
**TemplateMimeType** | Pointer to **string** |  | [optional] 

## Methods

### NewCreateNotificationPublisherRequest

`func NewCreateNotificationPublisherRequest(extensionName string, name string, ) *CreateNotificationPublisherRequest`

NewCreateNotificationPublisherRequest instantiates a new CreateNotificationPublisherRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateNotificationPublisherRequestWithDefaults

`func NewCreateNotificationPublisherRequestWithDefaults() *CreateNotificationPublisherRequest`

NewCreateNotificationPublisherRequestWithDefaults instantiates a new CreateNotificationPublisherRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *CreateNotificationPublisherRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateNotificationPublisherRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateNotificationPublisherRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreateNotificationPublisherRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetExtensionName

`func (o *CreateNotificationPublisherRequest) GetExtensionName() string`

GetExtensionName returns the ExtensionName field if non-nil, zero value otherwise.

### GetExtensionNameOk

`func (o *CreateNotificationPublisherRequest) GetExtensionNameOk() (*string, bool)`

GetExtensionNameOk returns a tuple with the ExtensionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtensionName

`func (o *CreateNotificationPublisherRequest) SetExtensionName(v string)`

SetExtensionName sets ExtensionName field to given value.


### GetName

`func (o *CreateNotificationPublisherRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateNotificationPublisherRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateNotificationPublisherRequest) SetName(v string)`

SetName sets Name field to given value.


### GetTemplate

`func (o *CreateNotificationPublisherRequest) GetTemplate() string`

GetTemplate returns the Template field if non-nil, zero value otherwise.

### GetTemplateOk

`func (o *CreateNotificationPublisherRequest) GetTemplateOk() (*string, bool)`

GetTemplateOk returns a tuple with the Template field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplate

`func (o *CreateNotificationPublisherRequest) SetTemplate(v string)`

SetTemplate sets Template field to given value.

### HasTemplate

`func (o *CreateNotificationPublisherRequest) HasTemplate() bool`

HasTemplate returns a boolean if a field has been set.

### GetTemplateMimeType

`func (o *CreateNotificationPublisherRequest) GetTemplateMimeType() string`

GetTemplateMimeType returns the TemplateMimeType field if non-nil, zero value otherwise.

### GetTemplateMimeTypeOk

`func (o *CreateNotificationPublisherRequest) GetTemplateMimeTypeOk() (*string, bool)`

GetTemplateMimeTypeOk returns a tuple with the TemplateMimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplateMimeType

`func (o *CreateNotificationPublisherRequest) SetTemplateMimeType(v string)`

SetTemplateMimeType sets TemplateMimeType field to given value.

### HasTemplateMimeType

`func (o *CreateNotificationPublisherRequest) HasTemplateMimeType() bool`

HasTemplateMimeType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


