# UserPermissionsSetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Permissions** | **[]string** |  | 
**Username** | **string** |  | 

## Methods

### NewUserPermissionsSetRequest

`func NewUserPermissionsSetRequest(permissions []string, username string, ) *UserPermissionsSetRequest`

NewUserPermissionsSetRequest instantiates a new UserPermissionsSetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserPermissionsSetRequestWithDefaults

`func NewUserPermissionsSetRequestWithDefaults() *UserPermissionsSetRequest`

NewUserPermissionsSetRequestWithDefaults instantiates a new UserPermissionsSetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPermissions

`func (o *UserPermissionsSetRequest) GetPermissions() []string`

GetPermissions returns the Permissions field if non-nil, zero value otherwise.

### GetPermissionsOk

`func (o *UserPermissionsSetRequest) GetPermissionsOk() (*[]string, bool)`

GetPermissionsOk returns a tuple with the Permissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissions

`func (o *UserPermissionsSetRequest) SetPermissions(v []string)`

SetPermissions sets Permissions field to given value.


### GetUsername

`func (o *UserPermissionsSetRequest) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *UserPermissionsSetRequest) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *UserPermissionsSetRequest) SetUsername(v string)`

SetUsername sets Username field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


