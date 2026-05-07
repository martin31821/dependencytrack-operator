# OidcUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** |  | [optional] 
**SubjectIdentifier** | Pointer to **string** |  | [optional] 
**Username** | **string** |  | 

## Methods

### NewOidcUser

`func NewOidcUser(username string, ) *OidcUser`

NewOidcUser instantiates a new OidcUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOidcUserWithDefaults

`func NewOidcUserWithDefaults() *OidcUser`

NewOidcUserWithDefaults instantiates a new OidcUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *OidcUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *OidcUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *OidcUser) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *OidcUser) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetSubjectIdentifier

`func (o *OidcUser) GetSubjectIdentifier() string`

GetSubjectIdentifier returns the SubjectIdentifier field if non-nil, zero value otherwise.

### GetSubjectIdentifierOk

`func (o *OidcUser) GetSubjectIdentifierOk() (*string, bool)`

GetSubjectIdentifierOk returns a tuple with the SubjectIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjectIdentifier

`func (o *OidcUser) SetSubjectIdentifier(v string)`

SetSubjectIdentifier sets SubjectIdentifier field to given value.

### HasSubjectIdentifier

`func (o *OidcUser) HasSubjectIdentifier() bool`

HasSubjectIdentifier returns a boolean if a field has been set.

### GetUsername

`func (o *OidcUser) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *OidcUser) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *OidcUser) SetUsername(v string)`

SetUsername sets Username field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


