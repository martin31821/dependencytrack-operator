# Tools

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Components** | Pointer to [**[]Component**](Component.md) |  | [optional] 
**Services** | Pointer to [**[]ServiceComponent**](ServiceComponent.md) |  | [optional] 

## Methods

### NewTools

`func NewTools() *Tools`

NewTools instantiates a new Tools object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewToolsWithDefaults

`func NewToolsWithDefaults() *Tools`

NewToolsWithDefaults instantiates a new Tools object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetComponents

`func (o *Tools) GetComponents() []Component`

GetComponents returns the Components field if non-nil, zero value otherwise.

### GetComponentsOk

`func (o *Tools) GetComponentsOk() (*[]Component, bool)`

GetComponentsOk returns a tuple with the Components field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponents

`func (o *Tools) SetComponents(v []Component)`

SetComponents sets Components field to given value.

### HasComponents

`func (o *Tools) HasComponents() bool`

HasComponents returns a boolean if a field has been set.

### GetServices

`func (o *Tools) GetServices() []ServiceComponent`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *Tools) GetServicesOk() (*[]ServiceComponent, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *Tools) SetServices(v []ServiceComponent)`

SetServices sets Services field to given value.

### HasServices

`func (o *Tools) HasServices() bool`

HasServices returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


