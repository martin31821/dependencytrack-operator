# TeamPermissionsSetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Permissions** | **[]string** |  | 
**Team** | **string** |  | 

## Methods

### NewTeamPermissionsSetRequest

`func NewTeamPermissionsSetRequest(permissions []string, team string, ) *TeamPermissionsSetRequest`

NewTeamPermissionsSetRequest instantiates a new TeamPermissionsSetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTeamPermissionsSetRequestWithDefaults

`func NewTeamPermissionsSetRequestWithDefaults() *TeamPermissionsSetRequest`

NewTeamPermissionsSetRequestWithDefaults instantiates a new TeamPermissionsSetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPermissions

`func (o *TeamPermissionsSetRequest) GetPermissions() []string`

GetPermissions returns the Permissions field if non-nil, zero value otherwise.

### GetPermissionsOk

`func (o *TeamPermissionsSetRequest) GetPermissionsOk() (*[]string, bool)`

GetPermissionsOk returns a tuple with the Permissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissions

`func (o *TeamPermissionsSetRequest) SetPermissions(v []string)`

SetPermissions sets Permissions field to given value.


### GetTeam

`func (o *TeamPermissionsSetRequest) GetTeam() string`

GetTeam returns the Team field if non-nil, zero value otherwise.

### GetTeamOk

`func (o *TeamPermissionsSetRequest) GetTeamOk() (*string, bool)`

GetTeamOk returns a tuple with the Team field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTeam

`func (o *TeamPermissionsSetRequest) SetTeam(v string)`

SetTeam sets Team field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


