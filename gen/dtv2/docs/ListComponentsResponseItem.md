# ListComponentsResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Version** | Pointer to **string** |  | [optional] 
**Group** | Pointer to **string** |  | [optional] 
**Classifier** | Pointer to **string** |  | [optional] 
**Scope** | Pointer to [**Scope**](Scope.md) |  | [optional] 
**Hashes** | Pointer to [**Hashes**](Hashes.md) |  | [optional] 
**Cpe** | Pointer to **string** |  | [optional] 
**Purl** | Pointer to **string** |  | [optional] 
**SwidTagId** | Pointer to **string** |  | [optional] 
**Internal** | Pointer to **bool** |  | [optional] 
**Copyright** | Pointer to **string** |  | [optional] 
**License** | Pointer to **string** |  | [optional] 
**LicenseExpression** | Pointer to **string** |  | [optional] 
**LicenseUrl** | Pointer to **string** |  | [optional] 
**ResolvedLicense** | Pointer to [**License**](License.md) |  | [optional] 
**LastInheritedRiskScore** | Pointer to **float64** |  | [optional] 
**Uuid** | **string** |  | 
**Project** | Pointer to [**ComponentProject**](ComponentProject.md) |  | [optional] 
**Metrics** | Pointer to [**DependencyMetrics**](DependencyMetrics.md) |  | [optional] 
**PackageMetadata** | Pointer to [**PackageMetadata**](PackageMetadata.md) |  | [optional] 
**PackageArtifactMetadata** | Pointer to [**PackageArtifactMetadata**](PackageArtifactMetadata.md) |  | [optional] 

## Methods

### NewListComponentsResponseItem

`func NewListComponentsResponseItem(name string, uuid string, ) *ListComponentsResponseItem`

NewListComponentsResponseItem instantiates a new ListComponentsResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListComponentsResponseItemWithDefaults

`func NewListComponentsResponseItemWithDefaults() *ListComponentsResponseItem`

NewListComponentsResponseItemWithDefaults instantiates a new ListComponentsResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ListComponentsResponseItem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListComponentsResponseItem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListComponentsResponseItem) SetName(v string)`

SetName sets Name field to given value.


### GetVersion

`func (o *ListComponentsResponseItem) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ListComponentsResponseItem) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ListComponentsResponseItem) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ListComponentsResponseItem) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetGroup

`func (o *ListComponentsResponseItem) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ListComponentsResponseItem) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ListComponentsResponseItem) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ListComponentsResponseItem) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetClassifier

`func (o *ListComponentsResponseItem) GetClassifier() string`

GetClassifier returns the Classifier field if non-nil, zero value otherwise.

### GetClassifierOk

`func (o *ListComponentsResponseItem) GetClassifierOk() (*string, bool)`

GetClassifierOk returns a tuple with the Classifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifier

`func (o *ListComponentsResponseItem) SetClassifier(v string)`

SetClassifier sets Classifier field to given value.

### HasClassifier

`func (o *ListComponentsResponseItem) HasClassifier() bool`

HasClassifier returns a boolean if a field has been set.

### GetScope

`func (o *ListComponentsResponseItem) GetScope() Scope`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *ListComponentsResponseItem) GetScopeOk() (*Scope, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *ListComponentsResponseItem) SetScope(v Scope)`

SetScope sets Scope field to given value.

### HasScope

`func (o *ListComponentsResponseItem) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetHashes

`func (o *ListComponentsResponseItem) GetHashes() Hashes`

GetHashes returns the Hashes field if non-nil, zero value otherwise.

### GetHashesOk

`func (o *ListComponentsResponseItem) GetHashesOk() (*Hashes, bool)`

GetHashesOk returns a tuple with the Hashes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHashes

`func (o *ListComponentsResponseItem) SetHashes(v Hashes)`

SetHashes sets Hashes field to given value.

### HasHashes

`func (o *ListComponentsResponseItem) HasHashes() bool`

HasHashes returns a boolean if a field has been set.

### GetCpe

`func (o *ListComponentsResponseItem) GetCpe() string`

GetCpe returns the Cpe field if non-nil, zero value otherwise.

### GetCpeOk

`func (o *ListComponentsResponseItem) GetCpeOk() (*string, bool)`

GetCpeOk returns a tuple with the Cpe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpe

`func (o *ListComponentsResponseItem) SetCpe(v string)`

SetCpe sets Cpe field to given value.

### HasCpe

`func (o *ListComponentsResponseItem) HasCpe() bool`

HasCpe returns a boolean if a field has been set.

### GetPurl

`func (o *ListComponentsResponseItem) GetPurl() string`

GetPurl returns the Purl field if non-nil, zero value otherwise.

### GetPurlOk

`func (o *ListComponentsResponseItem) GetPurlOk() (*string, bool)`

GetPurlOk returns a tuple with the Purl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurl

`func (o *ListComponentsResponseItem) SetPurl(v string)`

SetPurl sets Purl field to given value.

### HasPurl

`func (o *ListComponentsResponseItem) HasPurl() bool`

HasPurl returns a boolean if a field has been set.

### GetSwidTagId

`func (o *ListComponentsResponseItem) GetSwidTagId() string`

GetSwidTagId returns the SwidTagId field if non-nil, zero value otherwise.

### GetSwidTagIdOk

`func (o *ListComponentsResponseItem) GetSwidTagIdOk() (*string, bool)`

GetSwidTagIdOk returns a tuple with the SwidTagId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwidTagId

`func (o *ListComponentsResponseItem) SetSwidTagId(v string)`

SetSwidTagId sets SwidTagId field to given value.

### HasSwidTagId

`func (o *ListComponentsResponseItem) HasSwidTagId() bool`

HasSwidTagId returns a boolean if a field has been set.

### GetInternal

`func (o *ListComponentsResponseItem) GetInternal() bool`

GetInternal returns the Internal field if non-nil, zero value otherwise.

### GetInternalOk

`func (o *ListComponentsResponseItem) GetInternalOk() (*bool, bool)`

GetInternalOk returns a tuple with the Internal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternal

`func (o *ListComponentsResponseItem) SetInternal(v bool)`

SetInternal sets Internal field to given value.

### HasInternal

`func (o *ListComponentsResponseItem) HasInternal() bool`

HasInternal returns a boolean if a field has been set.

### GetCopyright

`func (o *ListComponentsResponseItem) GetCopyright() string`

GetCopyright returns the Copyright field if non-nil, zero value otherwise.

### GetCopyrightOk

`func (o *ListComponentsResponseItem) GetCopyrightOk() (*string, bool)`

GetCopyrightOk returns a tuple with the Copyright field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCopyright

`func (o *ListComponentsResponseItem) SetCopyright(v string)`

SetCopyright sets Copyright field to given value.

### HasCopyright

`func (o *ListComponentsResponseItem) HasCopyright() bool`

HasCopyright returns a boolean if a field has been set.

### GetLicense

`func (o *ListComponentsResponseItem) GetLicense() string`

GetLicense returns the License field if non-nil, zero value otherwise.

### GetLicenseOk

`func (o *ListComponentsResponseItem) GetLicenseOk() (*string, bool)`

GetLicenseOk returns a tuple with the License field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicense

`func (o *ListComponentsResponseItem) SetLicense(v string)`

SetLicense sets License field to given value.

### HasLicense

`func (o *ListComponentsResponseItem) HasLicense() bool`

HasLicense returns a boolean if a field has been set.

### GetLicenseExpression

`func (o *ListComponentsResponseItem) GetLicenseExpression() string`

GetLicenseExpression returns the LicenseExpression field if non-nil, zero value otherwise.

### GetLicenseExpressionOk

`func (o *ListComponentsResponseItem) GetLicenseExpressionOk() (*string, bool)`

GetLicenseExpressionOk returns a tuple with the LicenseExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseExpression

`func (o *ListComponentsResponseItem) SetLicenseExpression(v string)`

SetLicenseExpression sets LicenseExpression field to given value.

### HasLicenseExpression

`func (o *ListComponentsResponseItem) HasLicenseExpression() bool`

HasLicenseExpression returns a boolean if a field has been set.

### GetLicenseUrl

`func (o *ListComponentsResponseItem) GetLicenseUrl() string`

GetLicenseUrl returns the LicenseUrl field if non-nil, zero value otherwise.

### GetLicenseUrlOk

`func (o *ListComponentsResponseItem) GetLicenseUrlOk() (*string, bool)`

GetLicenseUrlOk returns a tuple with the LicenseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseUrl

`func (o *ListComponentsResponseItem) SetLicenseUrl(v string)`

SetLicenseUrl sets LicenseUrl field to given value.

### HasLicenseUrl

`func (o *ListComponentsResponseItem) HasLicenseUrl() bool`

HasLicenseUrl returns a boolean if a field has been set.

### GetResolvedLicense

`func (o *ListComponentsResponseItem) GetResolvedLicense() License`

GetResolvedLicense returns the ResolvedLicense field if non-nil, zero value otherwise.

### GetResolvedLicenseOk

`func (o *ListComponentsResponseItem) GetResolvedLicenseOk() (*License, bool)`

GetResolvedLicenseOk returns a tuple with the ResolvedLicense field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolvedLicense

`func (o *ListComponentsResponseItem) SetResolvedLicense(v License)`

SetResolvedLicense sets ResolvedLicense field to given value.

### HasResolvedLicense

`func (o *ListComponentsResponseItem) HasResolvedLicense() bool`

HasResolvedLicense returns a boolean if a field has been set.

### GetLastInheritedRiskScore

`func (o *ListComponentsResponseItem) GetLastInheritedRiskScore() float64`

GetLastInheritedRiskScore returns the LastInheritedRiskScore field if non-nil, zero value otherwise.

### GetLastInheritedRiskScoreOk

`func (o *ListComponentsResponseItem) GetLastInheritedRiskScoreOk() (*float64, bool)`

GetLastInheritedRiskScoreOk returns a tuple with the LastInheritedRiskScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastInheritedRiskScore

`func (o *ListComponentsResponseItem) SetLastInheritedRiskScore(v float64)`

SetLastInheritedRiskScore sets LastInheritedRiskScore field to given value.

### HasLastInheritedRiskScore

`func (o *ListComponentsResponseItem) HasLastInheritedRiskScore() bool`

HasLastInheritedRiskScore returns a boolean if a field has been set.

### GetUuid

`func (o *ListComponentsResponseItem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ListComponentsResponseItem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ListComponentsResponseItem) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetProject

`func (o *ListComponentsResponseItem) GetProject() ComponentProject`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *ListComponentsResponseItem) GetProjectOk() (*ComponentProject, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *ListComponentsResponseItem) SetProject(v ComponentProject)`

SetProject sets Project field to given value.

### HasProject

`func (o *ListComponentsResponseItem) HasProject() bool`

HasProject returns a boolean if a field has been set.

### GetMetrics

`func (o *ListComponentsResponseItem) GetMetrics() DependencyMetrics`

GetMetrics returns the Metrics field if non-nil, zero value otherwise.

### GetMetricsOk

`func (o *ListComponentsResponseItem) GetMetricsOk() (*DependencyMetrics, bool)`

GetMetricsOk returns a tuple with the Metrics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetrics

`func (o *ListComponentsResponseItem) SetMetrics(v DependencyMetrics)`

SetMetrics sets Metrics field to given value.

### HasMetrics

`func (o *ListComponentsResponseItem) HasMetrics() bool`

HasMetrics returns a boolean if a field has been set.

### GetPackageMetadata

`func (o *ListComponentsResponseItem) GetPackageMetadata() PackageMetadata`

GetPackageMetadata returns the PackageMetadata field if non-nil, zero value otherwise.

### GetPackageMetadataOk

`func (o *ListComponentsResponseItem) GetPackageMetadataOk() (*PackageMetadata, bool)`

GetPackageMetadataOk returns a tuple with the PackageMetadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageMetadata

`func (o *ListComponentsResponseItem) SetPackageMetadata(v PackageMetadata)`

SetPackageMetadata sets PackageMetadata field to given value.

### HasPackageMetadata

`func (o *ListComponentsResponseItem) HasPackageMetadata() bool`

HasPackageMetadata returns a boolean if a field has been set.

### GetPackageArtifactMetadata

`func (o *ListComponentsResponseItem) GetPackageArtifactMetadata() PackageArtifactMetadata`

GetPackageArtifactMetadata returns the PackageArtifactMetadata field if non-nil, zero value otherwise.

### GetPackageArtifactMetadataOk

`func (o *ListComponentsResponseItem) GetPackageArtifactMetadataOk() (*PackageArtifactMetadata, bool)`

GetPackageArtifactMetadataOk returns a tuple with the PackageArtifactMetadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageArtifactMetadata

`func (o *ListComponentsResponseItem) SetPackageArtifactMetadata(v PackageArtifactMetadata)`

SetPackageArtifactMetadata sets PackageArtifactMetadata field to given value.

### HasPackageArtifactMetadata

`func (o *ListComponentsResponseItem) HasPackageArtifactMetadata() bool`

HasPackageArtifactMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


