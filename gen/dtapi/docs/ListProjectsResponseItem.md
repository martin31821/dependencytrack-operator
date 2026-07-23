# ListProjectsResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** |  | [optional] 
**Authors** | Pointer to [**[]OrganizationalContact**](OrganizationalContact.md) |  | [optional] 
**Classifier** | Pointer to **string** |  | [optional] 
**CollectionLogic** | Pointer to **string** |  | [optional] 
**CollectionTag** | Pointer to [**Tag**](Tag.md) |  | [optional] 
**Cpe** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**DirectDependencies** | Pointer to **string** |  | [optional] 
**ExternalReferences** | Pointer to [**[]ExternalReference**](ExternalReference.md) |  | [optional] 
**Group** | Pointer to **string** |  | [optional] 
**HasChildren** | **bool** | Whether the project has child projects | 
**InactiveSince** | Pointer to **int64** | UNIX epoch timestamp in milliseconds | [optional] [readonly] 
**IsLatest** | Pointer to **bool** |  | [optional] 
**LastBomImport** | **int64** | UNIX epoch timestamp in milliseconds | 
**LastBomImportFormat** | Pointer to **string** |  | [optional] 
**LastInheritedRiskScore** | Pointer to **float64** |  | [optional] 
**LastVulnerabilityAnalysis** | Pointer to **int64** | UNIX epoch timestamp in milliseconds | [optional] 
**Manufacturer** | Pointer to [**OrganizationalEntity**](OrganizationalEntity.md) |  | [optional] 
**Metadata** | Pointer to [**ProjectMetadata**](ProjectMetadata.md) |  | [optional] 
**Metrics** | Pointer to [**ProjectMetrics**](ProjectMetrics.md) |  | [optional] 
**Name** | **string** |  | 
**Parent** | Pointer to [**Parent**](Parent.md) |  | [optional] 
**Publisher** | Pointer to **string** |  | [optional] 
**Purl** | Pointer to **string** |  | [optional] 
**Supplier** | Pointer to [**OrganizationalEntity**](OrganizationalEntity.md) |  | [optional] 
**SwidTagId** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to [**[]Tag**](Tag.md) |  | [optional] 
**Uuid** | **string** |  | 
**Version** | Pointer to **string** |  | [optional] 

## Methods

### NewListProjectsResponseItem

`func NewListProjectsResponseItem(hasChildren bool, lastBomImport int64, name string, uuid string, ) *ListProjectsResponseItem`

NewListProjectsResponseItem instantiates a new ListProjectsResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListProjectsResponseItemWithDefaults

`func NewListProjectsResponseItemWithDefaults() *ListProjectsResponseItem`

NewListProjectsResponseItemWithDefaults instantiates a new ListProjectsResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *ListProjectsResponseItem) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *ListProjectsResponseItem) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *ListProjectsResponseItem) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *ListProjectsResponseItem) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetAuthors

`func (o *ListProjectsResponseItem) GetAuthors() []OrganizationalContact`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ListProjectsResponseItem) GetAuthorsOk() (*[]OrganizationalContact, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ListProjectsResponseItem) SetAuthors(v []OrganizationalContact)`

SetAuthors sets Authors field to given value.

### HasAuthors

`func (o *ListProjectsResponseItem) HasAuthors() bool`

HasAuthors returns a boolean if a field has been set.

### GetClassifier

`func (o *ListProjectsResponseItem) GetClassifier() string`

GetClassifier returns the Classifier field if non-nil, zero value otherwise.

### GetClassifierOk

`func (o *ListProjectsResponseItem) GetClassifierOk() (*string, bool)`

GetClassifierOk returns a tuple with the Classifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifier

`func (o *ListProjectsResponseItem) SetClassifier(v string)`

SetClassifier sets Classifier field to given value.

### HasClassifier

`func (o *ListProjectsResponseItem) HasClassifier() bool`

HasClassifier returns a boolean if a field has been set.

### GetCollectionLogic

`func (o *ListProjectsResponseItem) GetCollectionLogic() string`

GetCollectionLogic returns the CollectionLogic field if non-nil, zero value otherwise.

### GetCollectionLogicOk

`func (o *ListProjectsResponseItem) GetCollectionLogicOk() (*string, bool)`

GetCollectionLogicOk returns a tuple with the CollectionLogic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollectionLogic

`func (o *ListProjectsResponseItem) SetCollectionLogic(v string)`

SetCollectionLogic sets CollectionLogic field to given value.

### HasCollectionLogic

`func (o *ListProjectsResponseItem) HasCollectionLogic() bool`

HasCollectionLogic returns a boolean if a field has been set.

### GetCollectionTag

`func (o *ListProjectsResponseItem) GetCollectionTag() Tag`

GetCollectionTag returns the CollectionTag field if non-nil, zero value otherwise.

### GetCollectionTagOk

`func (o *ListProjectsResponseItem) GetCollectionTagOk() (*Tag, bool)`

GetCollectionTagOk returns a tuple with the CollectionTag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollectionTag

`func (o *ListProjectsResponseItem) SetCollectionTag(v Tag)`

SetCollectionTag sets CollectionTag field to given value.

### HasCollectionTag

`func (o *ListProjectsResponseItem) HasCollectionTag() bool`

HasCollectionTag returns a boolean if a field has been set.

### GetCpe

`func (o *ListProjectsResponseItem) GetCpe() string`

GetCpe returns the Cpe field if non-nil, zero value otherwise.

### GetCpeOk

`func (o *ListProjectsResponseItem) GetCpeOk() (*string, bool)`

GetCpeOk returns a tuple with the Cpe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpe

`func (o *ListProjectsResponseItem) SetCpe(v string)`

SetCpe sets Cpe field to given value.

### HasCpe

`func (o *ListProjectsResponseItem) HasCpe() bool`

HasCpe returns a boolean if a field has been set.

### GetDescription

`func (o *ListProjectsResponseItem) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ListProjectsResponseItem) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ListProjectsResponseItem) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ListProjectsResponseItem) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetDirectDependencies

`func (o *ListProjectsResponseItem) GetDirectDependencies() string`

GetDirectDependencies returns the DirectDependencies field if non-nil, zero value otherwise.

### GetDirectDependenciesOk

`func (o *ListProjectsResponseItem) GetDirectDependenciesOk() (*string, bool)`

GetDirectDependenciesOk returns a tuple with the DirectDependencies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirectDependencies

`func (o *ListProjectsResponseItem) SetDirectDependencies(v string)`

SetDirectDependencies sets DirectDependencies field to given value.

### HasDirectDependencies

`func (o *ListProjectsResponseItem) HasDirectDependencies() bool`

HasDirectDependencies returns a boolean if a field has been set.

### GetExternalReferences

`func (o *ListProjectsResponseItem) GetExternalReferences() []ExternalReference`

GetExternalReferences returns the ExternalReferences field if non-nil, zero value otherwise.

### GetExternalReferencesOk

`func (o *ListProjectsResponseItem) GetExternalReferencesOk() (*[]ExternalReference, bool)`

GetExternalReferencesOk returns a tuple with the ExternalReferences field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalReferences

`func (o *ListProjectsResponseItem) SetExternalReferences(v []ExternalReference)`

SetExternalReferences sets ExternalReferences field to given value.

### HasExternalReferences

`func (o *ListProjectsResponseItem) HasExternalReferences() bool`

HasExternalReferences returns a boolean if a field has been set.

### GetGroup

`func (o *ListProjectsResponseItem) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ListProjectsResponseItem) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ListProjectsResponseItem) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ListProjectsResponseItem) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetHasChildren

`func (o *ListProjectsResponseItem) GetHasChildren() bool`

GetHasChildren returns the HasChildren field if non-nil, zero value otherwise.

### GetHasChildrenOk

`func (o *ListProjectsResponseItem) GetHasChildrenOk() (*bool, bool)`

GetHasChildrenOk returns a tuple with the HasChildren field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasChildren

`func (o *ListProjectsResponseItem) SetHasChildren(v bool)`

SetHasChildren sets HasChildren field to given value.


### GetInactiveSince

`func (o *ListProjectsResponseItem) GetInactiveSince() int64`

GetInactiveSince returns the InactiveSince field if non-nil, zero value otherwise.

### GetInactiveSinceOk

`func (o *ListProjectsResponseItem) GetInactiveSinceOk() (*int64, bool)`

GetInactiveSinceOk returns a tuple with the InactiveSince field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInactiveSince

`func (o *ListProjectsResponseItem) SetInactiveSince(v int64)`

SetInactiveSince sets InactiveSince field to given value.

### HasInactiveSince

`func (o *ListProjectsResponseItem) HasInactiveSince() bool`

HasInactiveSince returns a boolean if a field has been set.

### GetIsLatest

`func (o *ListProjectsResponseItem) GetIsLatest() bool`

GetIsLatest returns the IsLatest field if non-nil, zero value otherwise.

### GetIsLatestOk

`func (o *ListProjectsResponseItem) GetIsLatestOk() (*bool, bool)`

GetIsLatestOk returns a tuple with the IsLatest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLatest

`func (o *ListProjectsResponseItem) SetIsLatest(v bool)`

SetIsLatest sets IsLatest field to given value.

### HasIsLatest

`func (o *ListProjectsResponseItem) HasIsLatest() bool`

HasIsLatest returns a boolean if a field has been set.

### GetLastBomImport

`func (o *ListProjectsResponseItem) GetLastBomImport() int64`

GetLastBomImport returns the LastBomImport field if non-nil, zero value otherwise.

### GetLastBomImportOk

`func (o *ListProjectsResponseItem) GetLastBomImportOk() (*int64, bool)`

GetLastBomImportOk returns a tuple with the LastBomImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastBomImport

`func (o *ListProjectsResponseItem) SetLastBomImport(v int64)`

SetLastBomImport sets LastBomImport field to given value.


### GetLastBomImportFormat

`func (o *ListProjectsResponseItem) GetLastBomImportFormat() string`

GetLastBomImportFormat returns the LastBomImportFormat field if non-nil, zero value otherwise.

### GetLastBomImportFormatOk

`func (o *ListProjectsResponseItem) GetLastBomImportFormatOk() (*string, bool)`

GetLastBomImportFormatOk returns a tuple with the LastBomImportFormat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastBomImportFormat

`func (o *ListProjectsResponseItem) SetLastBomImportFormat(v string)`

SetLastBomImportFormat sets LastBomImportFormat field to given value.

### HasLastBomImportFormat

`func (o *ListProjectsResponseItem) HasLastBomImportFormat() bool`

HasLastBomImportFormat returns a boolean if a field has been set.

### GetLastInheritedRiskScore

`func (o *ListProjectsResponseItem) GetLastInheritedRiskScore() float64`

GetLastInheritedRiskScore returns the LastInheritedRiskScore field if non-nil, zero value otherwise.

### GetLastInheritedRiskScoreOk

`func (o *ListProjectsResponseItem) GetLastInheritedRiskScoreOk() (*float64, bool)`

GetLastInheritedRiskScoreOk returns a tuple with the LastInheritedRiskScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastInheritedRiskScore

`func (o *ListProjectsResponseItem) SetLastInheritedRiskScore(v float64)`

SetLastInheritedRiskScore sets LastInheritedRiskScore field to given value.

### HasLastInheritedRiskScore

`func (o *ListProjectsResponseItem) HasLastInheritedRiskScore() bool`

HasLastInheritedRiskScore returns a boolean if a field has been set.

### GetLastVulnerabilityAnalysis

`func (o *ListProjectsResponseItem) GetLastVulnerabilityAnalysis() int64`

GetLastVulnerabilityAnalysis returns the LastVulnerabilityAnalysis field if non-nil, zero value otherwise.

### GetLastVulnerabilityAnalysisOk

`func (o *ListProjectsResponseItem) GetLastVulnerabilityAnalysisOk() (*int64, bool)`

GetLastVulnerabilityAnalysisOk returns a tuple with the LastVulnerabilityAnalysis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastVulnerabilityAnalysis

`func (o *ListProjectsResponseItem) SetLastVulnerabilityAnalysis(v int64)`

SetLastVulnerabilityAnalysis sets LastVulnerabilityAnalysis field to given value.

### HasLastVulnerabilityAnalysis

`func (o *ListProjectsResponseItem) HasLastVulnerabilityAnalysis() bool`

HasLastVulnerabilityAnalysis returns a boolean if a field has been set.

### GetManufacturer

`func (o *ListProjectsResponseItem) GetManufacturer() OrganizationalEntity`

GetManufacturer returns the Manufacturer field if non-nil, zero value otherwise.

### GetManufacturerOk

`func (o *ListProjectsResponseItem) GetManufacturerOk() (*OrganizationalEntity, bool)`

GetManufacturerOk returns a tuple with the Manufacturer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManufacturer

`func (o *ListProjectsResponseItem) SetManufacturer(v OrganizationalEntity)`

SetManufacturer sets Manufacturer field to given value.

### HasManufacturer

`func (o *ListProjectsResponseItem) HasManufacturer() bool`

HasManufacturer returns a boolean if a field has been set.

### GetMetadata

`func (o *ListProjectsResponseItem) GetMetadata() ProjectMetadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ListProjectsResponseItem) GetMetadataOk() (*ProjectMetadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ListProjectsResponseItem) SetMetadata(v ProjectMetadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ListProjectsResponseItem) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetMetrics

`func (o *ListProjectsResponseItem) GetMetrics() ProjectMetrics`

GetMetrics returns the Metrics field if non-nil, zero value otherwise.

### GetMetricsOk

`func (o *ListProjectsResponseItem) GetMetricsOk() (*ProjectMetrics, bool)`

GetMetricsOk returns a tuple with the Metrics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetrics

`func (o *ListProjectsResponseItem) SetMetrics(v ProjectMetrics)`

SetMetrics sets Metrics field to given value.

### HasMetrics

`func (o *ListProjectsResponseItem) HasMetrics() bool`

HasMetrics returns a boolean if a field has been set.

### GetName

`func (o *ListProjectsResponseItem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListProjectsResponseItem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListProjectsResponseItem) SetName(v string)`

SetName sets Name field to given value.


### GetParent

`func (o *ListProjectsResponseItem) GetParent() Parent`

GetParent returns the Parent field if non-nil, zero value otherwise.

### GetParentOk

`func (o *ListProjectsResponseItem) GetParentOk() (*Parent, bool)`

GetParentOk returns a tuple with the Parent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParent

`func (o *ListProjectsResponseItem) SetParent(v Parent)`

SetParent sets Parent field to given value.

### HasParent

`func (o *ListProjectsResponseItem) HasParent() bool`

HasParent returns a boolean if a field has been set.

### GetPublisher

`func (o *ListProjectsResponseItem) GetPublisher() string`

GetPublisher returns the Publisher field if non-nil, zero value otherwise.

### GetPublisherOk

`func (o *ListProjectsResponseItem) GetPublisherOk() (*string, bool)`

GetPublisherOk returns a tuple with the Publisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublisher

`func (o *ListProjectsResponseItem) SetPublisher(v string)`

SetPublisher sets Publisher field to given value.

### HasPublisher

`func (o *ListProjectsResponseItem) HasPublisher() bool`

HasPublisher returns a boolean if a field has been set.

### GetPurl

`func (o *ListProjectsResponseItem) GetPurl() string`

GetPurl returns the Purl field if non-nil, zero value otherwise.

### GetPurlOk

`func (o *ListProjectsResponseItem) GetPurlOk() (*string, bool)`

GetPurlOk returns a tuple with the Purl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurl

`func (o *ListProjectsResponseItem) SetPurl(v string)`

SetPurl sets Purl field to given value.

### HasPurl

`func (o *ListProjectsResponseItem) HasPurl() bool`

HasPurl returns a boolean if a field has been set.

### GetSupplier

`func (o *ListProjectsResponseItem) GetSupplier() OrganizationalEntity`

GetSupplier returns the Supplier field if non-nil, zero value otherwise.

### GetSupplierOk

`func (o *ListProjectsResponseItem) GetSupplierOk() (*OrganizationalEntity, bool)`

GetSupplierOk returns a tuple with the Supplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupplier

`func (o *ListProjectsResponseItem) SetSupplier(v OrganizationalEntity)`

SetSupplier sets Supplier field to given value.

### HasSupplier

`func (o *ListProjectsResponseItem) HasSupplier() bool`

HasSupplier returns a boolean if a field has been set.

### GetSwidTagId

`func (o *ListProjectsResponseItem) GetSwidTagId() string`

GetSwidTagId returns the SwidTagId field if non-nil, zero value otherwise.

### GetSwidTagIdOk

`func (o *ListProjectsResponseItem) GetSwidTagIdOk() (*string, bool)`

GetSwidTagIdOk returns a tuple with the SwidTagId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwidTagId

`func (o *ListProjectsResponseItem) SetSwidTagId(v string)`

SetSwidTagId sets SwidTagId field to given value.

### HasSwidTagId

`func (o *ListProjectsResponseItem) HasSwidTagId() bool`

HasSwidTagId returns a boolean if a field has been set.

### GetTags

`func (o *ListProjectsResponseItem) GetTags() []Tag`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ListProjectsResponseItem) GetTagsOk() (*[]Tag, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ListProjectsResponseItem) SetTags(v []Tag)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ListProjectsResponseItem) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetUuid

`func (o *ListProjectsResponseItem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ListProjectsResponseItem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ListProjectsResponseItem) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetVersion

`func (o *ListProjectsResponseItem) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ListProjectsResponseItem) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ListProjectsResponseItem) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ListProjectsResponseItem) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


