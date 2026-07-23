# ListExtensionsResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Configurable** | **bool** | Whether the extension supports runtime configuration. | 
**Testable** | **bool** | Whether the extension can be tested. | 

## Methods

### NewListExtensionsResponseItem

`func NewListExtensionsResponseItem(name string, configurable bool, testable bool, ) *ListExtensionsResponseItem`

NewListExtensionsResponseItem instantiates a new ListExtensionsResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListExtensionsResponseItemWithDefaults

`func NewListExtensionsResponseItemWithDefaults() *ListExtensionsResponseItem`

NewListExtensionsResponseItemWithDefaults instantiates a new ListExtensionsResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ListExtensionsResponseItem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListExtensionsResponseItem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListExtensionsResponseItem) SetName(v string)`

SetName sets Name field to given value.


### GetConfigurable

`func (o *ListExtensionsResponseItem) GetConfigurable() bool`

GetConfigurable returns the Configurable field if non-nil, zero value otherwise.

### GetConfigurableOk

`func (o *ListExtensionsResponseItem) GetConfigurableOk() (*bool, bool)`

GetConfigurableOk returns a tuple with the Configurable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigurable

`func (o *ListExtensionsResponseItem) SetConfigurable(v bool)`

SetConfigurable sets Configurable field to given value.


### GetTestable

`func (o *ListExtensionsResponseItem) GetTestable() bool`

GetTestable returns the Testable field if non-nil, zero value otherwise.

### GetTestableOk

`func (o *ListExtensionsResponseItem) GetTestableOk() (*bool, bool)`

GetTestableOk returns a tuple with the Testable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTestable

`func (o *ListExtensionsResponseItem) SetTestable(v bool)`

SetTestable sets Testable field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


