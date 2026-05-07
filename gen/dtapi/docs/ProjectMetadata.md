# ProjectMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authors** | Pointer to [**[]OrganizationalContact**](OrganizationalContact.md) |  | [optional] 
**Supplier** | Pointer to [**OrganizationalEntity**](OrganizationalEntity.md) |  | [optional] 
**Tools** | Pointer to [**Tools**](Tools.md) |  | [optional] 

## Methods

### NewProjectMetadata

`func NewProjectMetadata() *ProjectMetadata`

NewProjectMetadata instantiates a new ProjectMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectMetadataWithDefaults

`func NewProjectMetadataWithDefaults() *ProjectMetadata`

NewProjectMetadataWithDefaults instantiates a new ProjectMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthors

`func (o *ProjectMetadata) GetAuthors() []OrganizationalContact`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ProjectMetadata) GetAuthorsOk() (*[]OrganizationalContact, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ProjectMetadata) SetAuthors(v []OrganizationalContact)`

SetAuthors sets Authors field to given value.

### HasAuthors

`func (o *ProjectMetadata) HasAuthors() bool`

HasAuthors returns a boolean if a field has been set.

### GetSupplier

`func (o *ProjectMetadata) GetSupplier() OrganizationalEntity`

GetSupplier returns the Supplier field if non-nil, zero value otherwise.

### GetSupplierOk

`func (o *ProjectMetadata) GetSupplierOk() (*OrganizationalEntity, bool)`

GetSupplierOk returns a tuple with the Supplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupplier

`func (o *ProjectMetadata) SetSupplier(v OrganizationalEntity)`

SetSupplier sets Supplier field to given value.

### HasSupplier

`func (o *ProjectMetadata) HasSupplier() bool`

HasSupplier returns a boolean if a field has been set.

### GetTools

`func (o *ProjectMetadata) GetTools() Tools`

GetTools returns the Tools field if non-nil, zero value otherwise.

### GetToolsOk

`func (o *ProjectMetadata) GetToolsOk() (*Tools, bool)`

GetToolsOk returns a tuple with the Tools field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTools

`func (o *ProjectMetadata) SetTools(v Tools)`

SetTools sets Tools field to given value.

### HasTools

`func (o *ProjectMetadata) HasTools() bool`

HasTools returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


