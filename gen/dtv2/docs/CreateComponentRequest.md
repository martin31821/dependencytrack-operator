# CreateComponentRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectUuid** | **string** |  | 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Group** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**Classifier** | Pointer to [**Classifier**](Classifier.md) |  | [optional] 
**Filename** | Pointer to **string** |  | [optional] 
**Extension** | Pointer to **string** |  | [optional] 
**Hashes** | Pointer to [**Hashes**](Hashes.md) |  | [optional] 
**Cpe** | Pointer to **string** |  | [optional] 
**Publisher** | Pointer to **string** |  | [optional] 
**Supplier** | Pointer to [**OrganizationalEntity**](OrganizationalEntity.md) |  | [optional] 
**Authors** | Pointer to [**[]OrganizationalContact**](OrganizationalContact.md) |  | [optional] 
**Purl** | Pointer to **string** |  | [optional] 
**SwidTagId** | Pointer to **string** |  | [optional] 
**Internal** | Pointer to **bool** |  | [optional] 
**Copyright** | Pointer to **string** |  | [optional] 
**License** | Pointer to **string** |  | [optional] 
**LicenseExpression** | Pointer to **string** |  | [optional] 
**LicenseUrl** | Pointer to **string** |  | [optional] 
**Notes** | Pointer to **string** |  | [optional] 

## Methods

### NewCreateComponentRequest

`func NewCreateComponentRequest(projectUuid string, name string, ) *CreateComponentRequest`

NewCreateComponentRequest instantiates a new CreateComponentRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateComponentRequestWithDefaults

`func NewCreateComponentRequestWithDefaults() *CreateComponentRequest`

NewCreateComponentRequestWithDefaults instantiates a new CreateComponentRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProjectUuid

`func (o *CreateComponentRequest) GetProjectUuid() string`

GetProjectUuid returns the ProjectUuid field if non-nil, zero value otherwise.

### GetProjectUuidOk

`func (o *CreateComponentRequest) GetProjectUuidOk() (*string, bool)`

GetProjectUuidOk returns a tuple with the ProjectUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectUuid

`func (o *CreateComponentRequest) SetProjectUuid(v string)`

SetProjectUuid sets ProjectUuid field to given value.


### GetName

`func (o *CreateComponentRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateComponentRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateComponentRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CreateComponentRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateComponentRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateComponentRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreateComponentRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetGroup

`func (o *CreateComponentRequest) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *CreateComponentRequest) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *CreateComponentRequest) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *CreateComponentRequest) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetVersion

`func (o *CreateComponentRequest) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *CreateComponentRequest) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *CreateComponentRequest) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *CreateComponentRequest) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetClassifier

`func (o *CreateComponentRequest) GetClassifier() Classifier`

GetClassifier returns the Classifier field if non-nil, zero value otherwise.

### GetClassifierOk

`func (o *CreateComponentRequest) GetClassifierOk() (*Classifier, bool)`

GetClassifierOk returns a tuple with the Classifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifier

`func (o *CreateComponentRequest) SetClassifier(v Classifier)`

SetClassifier sets Classifier field to given value.

### HasClassifier

`func (o *CreateComponentRequest) HasClassifier() bool`

HasClassifier returns a boolean if a field has been set.

### GetFilename

`func (o *CreateComponentRequest) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *CreateComponentRequest) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *CreateComponentRequest) SetFilename(v string)`

SetFilename sets Filename field to given value.

### HasFilename

`func (o *CreateComponentRequest) HasFilename() bool`

HasFilename returns a boolean if a field has been set.

### GetExtension

`func (o *CreateComponentRequest) GetExtension() string`

GetExtension returns the Extension field if non-nil, zero value otherwise.

### GetExtensionOk

`func (o *CreateComponentRequest) GetExtensionOk() (*string, bool)`

GetExtensionOk returns a tuple with the Extension field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtension

`func (o *CreateComponentRequest) SetExtension(v string)`

SetExtension sets Extension field to given value.

### HasExtension

`func (o *CreateComponentRequest) HasExtension() bool`

HasExtension returns a boolean if a field has been set.

### GetHashes

`func (o *CreateComponentRequest) GetHashes() Hashes`

GetHashes returns the Hashes field if non-nil, zero value otherwise.

### GetHashesOk

`func (o *CreateComponentRequest) GetHashesOk() (*Hashes, bool)`

GetHashesOk returns a tuple with the Hashes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHashes

`func (o *CreateComponentRequest) SetHashes(v Hashes)`

SetHashes sets Hashes field to given value.

### HasHashes

`func (o *CreateComponentRequest) HasHashes() bool`

HasHashes returns a boolean if a field has been set.

### GetCpe

`func (o *CreateComponentRequest) GetCpe() string`

GetCpe returns the Cpe field if non-nil, zero value otherwise.

### GetCpeOk

`func (o *CreateComponentRequest) GetCpeOk() (*string, bool)`

GetCpeOk returns a tuple with the Cpe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpe

`func (o *CreateComponentRequest) SetCpe(v string)`

SetCpe sets Cpe field to given value.

### HasCpe

`func (o *CreateComponentRequest) HasCpe() bool`

HasCpe returns a boolean if a field has been set.

### GetPublisher

`func (o *CreateComponentRequest) GetPublisher() string`

GetPublisher returns the Publisher field if non-nil, zero value otherwise.

### GetPublisherOk

`func (o *CreateComponentRequest) GetPublisherOk() (*string, bool)`

GetPublisherOk returns a tuple with the Publisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublisher

`func (o *CreateComponentRequest) SetPublisher(v string)`

SetPublisher sets Publisher field to given value.

### HasPublisher

`func (o *CreateComponentRequest) HasPublisher() bool`

HasPublisher returns a boolean if a field has been set.

### GetSupplier

`func (o *CreateComponentRequest) GetSupplier() OrganizationalEntity`

GetSupplier returns the Supplier field if non-nil, zero value otherwise.

### GetSupplierOk

`func (o *CreateComponentRequest) GetSupplierOk() (*OrganizationalEntity, bool)`

GetSupplierOk returns a tuple with the Supplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupplier

`func (o *CreateComponentRequest) SetSupplier(v OrganizationalEntity)`

SetSupplier sets Supplier field to given value.

### HasSupplier

`func (o *CreateComponentRequest) HasSupplier() bool`

HasSupplier returns a boolean if a field has been set.

### GetAuthors

`func (o *CreateComponentRequest) GetAuthors() []OrganizationalContact`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *CreateComponentRequest) GetAuthorsOk() (*[]OrganizationalContact, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *CreateComponentRequest) SetAuthors(v []OrganizationalContact)`

SetAuthors sets Authors field to given value.

### HasAuthors

`func (o *CreateComponentRequest) HasAuthors() bool`

HasAuthors returns a boolean if a field has been set.

### GetPurl

`func (o *CreateComponentRequest) GetPurl() string`

GetPurl returns the Purl field if non-nil, zero value otherwise.

### GetPurlOk

`func (o *CreateComponentRequest) GetPurlOk() (*string, bool)`

GetPurlOk returns a tuple with the Purl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurl

`func (o *CreateComponentRequest) SetPurl(v string)`

SetPurl sets Purl field to given value.

### HasPurl

`func (o *CreateComponentRequest) HasPurl() bool`

HasPurl returns a boolean if a field has been set.

### GetSwidTagId

`func (o *CreateComponentRequest) GetSwidTagId() string`

GetSwidTagId returns the SwidTagId field if non-nil, zero value otherwise.

### GetSwidTagIdOk

`func (o *CreateComponentRequest) GetSwidTagIdOk() (*string, bool)`

GetSwidTagIdOk returns a tuple with the SwidTagId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwidTagId

`func (o *CreateComponentRequest) SetSwidTagId(v string)`

SetSwidTagId sets SwidTagId field to given value.

### HasSwidTagId

`func (o *CreateComponentRequest) HasSwidTagId() bool`

HasSwidTagId returns a boolean if a field has been set.

### GetInternal

`func (o *CreateComponentRequest) GetInternal() bool`

GetInternal returns the Internal field if non-nil, zero value otherwise.

### GetInternalOk

`func (o *CreateComponentRequest) GetInternalOk() (*bool, bool)`

GetInternalOk returns a tuple with the Internal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternal

`func (o *CreateComponentRequest) SetInternal(v bool)`

SetInternal sets Internal field to given value.

### HasInternal

`func (o *CreateComponentRequest) HasInternal() bool`

HasInternal returns a boolean if a field has been set.

### GetCopyright

`func (o *CreateComponentRequest) GetCopyright() string`

GetCopyright returns the Copyright field if non-nil, zero value otherwise.

### GetCopyrightOk

`func (o *CreateComponentRequest) GetCopyrightOk() (*string, bool)`

GetCopyrightOk returns a tuple with the Copyright field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCopyright

`func (o *CreateComponentRequest) SetCopyright(v string)`

SetCopyright sets Copyright field to given value.

### HasCopyright

`func (o *CreateComponentRequest) HasCopyright() bool`

HasCopyright returns a boolean if a field has been set.

### GetLicense

`func (o *CreateComponentRequest) GetLicense() string`

GetLicense returns the License field if non-nil, zero value otherwise.

### GetLicenseOk

`func (o *CreateComponentRequest) GetLicenseOk() (*string, bool)`

GetLicenseOk returns a tuple with the License field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicense

`func (o *CreateComponentRequest) SetLicense(v string)`

SetLicense sets License field to given value.

### HasLicense

`func (o *CreateComponentRequest) HasLicense() bool`

HasLicense returns a boolean if a field has been set.

### GetLicenseExpression

`func (o *CreateComponentRequest) GetLicenseExpression() string`

GetLicenseExpression returns the LicenseExpression field if non-nil, zero value otherwise.

### GetLicenseExpressionOk

`func (o *CreateComponentRequest) GetLicenseExpressionOk() (*string, bool)`

GetLicenseExpressionOk returns a tuple with the LicenseExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseExpression

`func (o *CreateComponentRequest) SetLicenseExpression(v string)`

SetLicenseExpression sets LicenseExpression field to given value.

### HasLicenseExpression

`func (o *CreateComponentRequest) HasLicenseExpression() bool`

HasLicenseExpression returns a boolean if a field has been set.

### GetLicenseUrl

`func (o *CreateComponentRequest) GetLicenseUrl() string`

GetLicenseUrl returns the LicenseUrl field if non-nil, zero value otherwise.

### GetLicenseUrlOk

`func (o *CreateComponentRequest) GetLicenseUrlOk() (*string, bool)`

GetLicenseUrlOk returns a tuple with the LicenseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseUrl

`func (o *CreateComponentRequest) SetLicenseUrl(v string)`

SetLicenseUrl sets LicenseUrl field to given value.

### HasLicenseUrl

`func (o *CreateComponentRequest) HasLicenseUrl() bool`

HasLicenseUrl returns a boolean if a field has been set.

### GetNotes

`func (o *CreateComponentRequest) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *CreateComponentRequest) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *CreateComponentRequest) SetNotes(v string)`

SetNotes sets Notes field to given value.

### HasNotes

`func (o *CreateComponentRequest) HasNotes() bool`

HasNotes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


