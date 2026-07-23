# ExtensionTestCheck

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Status** | [**ExtensionTestCheckStatus**](ExtensionTestCheckStatus.md) |  | 
**Message** | Pointer to **string** |  | [optional] 

## Methods

### NewExtensionTestCheck

`func NewExtensionTestCheck(name string, status ExtensionTestCheckStatus, ) *ExtensionTestCheck`

NewExtensionTestCheck instantiates a new ExtensionTestCheck object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExtensionTestCheckWithDefaults

`func NewExtensionTestCheckWithDefaults() *ExtensionTestCheck`

NewExtensionTestCheckWithDefaults instantiates a new ExtensionTestCheck object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ExtensionTestCheck) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ExtensionTestCheck) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ExtensionTestCheck) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *ExtensionTestCheck) GetStatus() ExtensionTestCheckStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ExtensionTestCheck) GetStatusOk() (*ExtensionTestCheckStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ExtensionTestCheck) SetStatus(v ExtensionTestCheckStatus)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *ExtensionTestCheck) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ExtensionTestCheck) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ExtensionTestCheck) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ExtensionTestCheck) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


