# Permission

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** |  | [optional] 
**LdapUsers** | Pointer to [**[]LdapUser**](LdapUser.md) |  | [optional] 
**ManagedUsers** | Pointer to [**[]ManagedUser**](ManagedUser.md) |  | [optional] 
**Name** | **string** |  | 
**OidcUsers** | Pointer to [**[]OidcUser**](OidcUser.md) |  | [optional] 

## Methods

### NewPermission

`func NewPermission(name string, ) *Permission`

NewPermission instantiates a new Permission object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPermissionWithDefaults

`func NewPermissionWithDefaults() *Permission`

NewPermissionWithDefaults instantiates a new Permission object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *Permission) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Permission) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Permission) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Permission) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetLdapUsers

`func (o *Permission) GetLdapUsers() []LdapUser`

GetLdapUsers returns the LdapUsers field if non-nil, zero value otherwise.

### GetLdapUsersOk

`func (o *Permission) GetLdapUsersOk() (*[]LdapUser, bool)`

GetLdapUsersOk returns a tuple with the LdapUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLdapUsers

`func (o *Permission) SetLdapUsers(v []LdapUser)`

SetLdapUsers sets LdapUsers field to given value.

### HasLdapUsers

`func (o *Permission) HasLdapUsers() bool`

HasLdapUsers returns a boolean if a field has been set.

### GetManagedUsers

`func (o *Permission) GetManagedUsers() []ManagedUser`

GetManagedUsers returns the ManagedUsers field if non-nil, zero value otherwise.

### GetManagedUsersOk

`func (o *Permission) GetManagedUsersOk() (*[]ManagedUser, bool)`

GetManagedUsersOk returns a tuple with the ManagedUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManagedUsers

`func (o *Permission) SetManagedUsers(v []ManagedUser)`

SetManagedUsers sets ManagedUsers field to given value.

### HasManagedUsers

`func (o *Permission) HasManagedUsers() bool`

HasManagedUsers returns a boolean if a field has been set.

### GetName

`func (o *Permission) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Permission) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Permission) SetName(v string)`

SetName sets Name field to given value.


### GetOidcUsers

`func (o *Permission) GetOidcUsers() []OidcUser`

GetOidcUsers returns the OidcUsers field if non-nil, zero value otherwise.

### GetOidcUsersOk

`func (o *Permission) GetOidcUsersOk() (*[]OidcUser, bool)`

GetOidcUsersOk returns a tuple with the OidcUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOidcUsers

`func (o *Permission) SetOidcUsers(v []OidcUser)`

SetOidcUsers sets OidcUsers field to given value.

### HasOidcUsers

`func (o *Permission) HasOidcUsers() bool`

HasOidcUsers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


