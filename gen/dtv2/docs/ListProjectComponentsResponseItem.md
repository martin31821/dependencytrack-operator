# ListProjectComponentsResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Version** | Pointer to **string** |  | [optional] 
**Group** | Pointer to **string** |  | [optional] 
**Classifier** | Pointer to **string** |  | [optional] 
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
**OccurrenceCount** | Pointer to **int64** |  | [optional] 
**Scope** | Pointer to [**Scope**](Scope.md) |  | [optional] 
**LastInheritedRiskScore** | Pointer to **float64** |  | [optional] 
**Uuid** | **string** |  | 
**Metrics** | Pointer to [**DependencyMetrics**](DependencyMetrics.md) |  | [optional] 
**PackageMetadata** | Pointer to [**PackageMetadata**](PackageMetadata.md) |  | [optional] 
**PackageArtifactMetadata** | Pointer to [**PackageArtifactMetadata**](PackageArtifactMetadata.md) |  | [optional] 

## Methods

### NewListProjectComponentsResponseItem

`func NewListProjectComponentsResponseItem(name string, uuid string, ) *ListProjectComponentsResponseItem`

NewListProjectComponentsResponseItem instantiates a new ListProjectComponentsResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListProjectComponentsResponseItemWithDefaults

`func NewListProjectComponentsResponseItemWithDefaults() *ListProjectComponentsResponseItem`

NewListProjectComponentsResponseItemWithDefaults instantiates a new ListProjectComponentsResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ListProjectComponentsResponseItem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListProjectComponentsResponseItem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListProjectComponentsResponseItem) SetName(v string)`

SetName sets Name field to given value.


### GetVersion

`func (o *ListProjectComponentsResponseItem) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ListProjectComponentsResponseItem) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ListProjectComponentsResponseItem) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ListProjectComponentsResponseItem) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetGroup

`func (o *ListProjectComponentsResponseItem) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ListProjectComponentsResponseItem) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ListProjectComponentsResponseItem) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ListProjectComponentsResponseItem) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetClassifier

`func (o *ListProjectComponentsResponseItem) GetClassifier() string`

GetClassifier returns the Classifier field if non-nil, zero value otherwise.

### GetClassifierOk

`func (o *ListProjectComponentsResponseItem) GetClassifierOk() (*string, bool)`

GetClassifierOk returns a tuple with the Classifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifier

`func (o *ListProjectComponentsResponseItem) SetClassifier(v string)`

SetClassifier sets Classifier field to given value.

### HasClassifier

`func (o *ListProjectComponentsResponseItem) HasClassifier() bool`

HasClassifier returns a boolean if a field has been set.

### GetHashes

`func (o *ListProjectComponentsResponseItem) GetHashes() Hashes`

GetHashes returns the Hashes field if non-nil, zero value otherwise.

### GetHashesOk

`func (o *ListProjectComponentsResponseItem) GetHashesOk() (*Hashes, bool)`

GetHashesOk returns a tuple with the Hashes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHashes

`func (o *ListProjectComponentsResponseItem) SetHashes(v Hashes)`

SetHashes sets Hashes field to given value.

### HasHashes

`func (o *ListProjectComponentsResponseItem) HasHashes() bool`

HasHashes returns a boolean if a field has been set.

### GetCpe

`func (o *ListProjectComponentsResponseItem) GetCpe() string`

GetCpe returns the Cpe field if non-nil, zero value otherwise.

### GetCpeOk

`func (o *ListProjectComponentsResponseItem) GetCpeOk() (*string, bool)`

GetCpeOk returns a tuple with the Cpe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpe

`func (o *ListProjectComponentsResponseItem) SetCpe(v string)`

SetCpe sets Cpe field to given value.

### HasCpe

`func (o *ListProjectComponentsResponseItem) HasCpe() bool`

HasCpe returns a boolean if a field has been set.

### GetPurl

`func (o *ListProjectComponentsResponseItem) GetPurl() string`

GetPurl returns the Purl field if non-nil, zero value otherwise.

### GetPurlOk

`func (o *ListProjectComponentsResponseItem) GetPurlOk() (*string, bool)`

GetPurlOk returns a tuple with the Purl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurl

`func (o *ListProjectComponentsResponseItem) SetPurl(v string)`

SetPurl sets Purl field to given value.

### HasPurl

`func (o *ListProjectComponentsResponseItem) HasPurl() bool`

HasPurl returns a boolean if a field has been set.

### GetSwidTagId

`func (o *ListProjectComponentsResponseItem) GetSwidTagId() string`

GetSwidTagId returns the SwidTagId field if non-nil, zero value otherwise.

### GetSwidTagIdOk

`func (o *ListProjectComponentsResponseItem) GetSwidTagIdOk() (*string, bool)`

GetSwidTagIdOk returns a tuple with the SwidTagId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwidTagId

`func (o *ListProjectComponentsResponseItem) SetSwidTagId(v string)`

SetSwidTagId sets SwidTagId field to given value.

### HasSwidTagId

`func (o *ListProjectComponentsResponseItem) HasSwidTagId() bool`

HasSwidTagId returns a boolean if a field has been set.

### GetInternal

`func (o *ListProjectComponentsResponseItem) GetInternal() bool`

GetInternal returns the Internal field if non-nil, zero value otherwise.

### GetInternalOk

`func (o *ListProjectComponentsResponseItem) GetInternalOk() (*bool, bool)`

GetInternalOk returns a tuple with the Internal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternal

`func (o *ListProjectComponentsResponseItem) SetInternal(v bool)`

SetInternal sets Internal field to given value.

### HasInternal

`func (o *ListProjectComponentsResponseItem) HasInternal() bool`

HasInternal returns a boolean if a field has been set.

### GetCopyright

`func (o *ListProjectComponentsResponseItem) GetCopyright() string`

GetCopyright returns the Copyright field if non-nil, zero value otherwise.

### GetCopyrightOk

`func (o *ListProjectComponentsResponseItem) GetCopyrightOk() (*string, bool)`

GetCopyrightOk returns a tuple with the Copyright field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCopyright

`func (o *ListProjectComponentsResponseItem) SetCopyright(v string)`

SetCopyright sets Copyright field to given value.

### HasCopyright

`func (o *ListProjectComponentsResponseItem) HasCopyright() bool`

HasCopyright returns a boolean if a field has been set.

### GetLicense

`func (o *ListProjectComponentsResponseItem) GetLicense() string`

GetLicense returns the License field if non-nil, zero value otherwise.

### GetLicenseOk

`func (o *ListProjectComponentsResponseItem) GetLicenseOk() (*string, bool)`

GetLicenseOk returns a tuple with the License field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicense

`func (o *ListProjectComponentsResponseItem) SetLicense(v string)`

SetLicense sets License field to given value.

### HasLicense

`func (o *ListProjectComponentsResponseItem) HasLicense() bool`

HasLicense returns a boolean if a field has been set.

### GetLicenseExpression

`func (o *ListProjectComponentsResponseItem) GetLicenseExpression() string`

GetLicenseExpression returns the LicenseExpression field if non-nil, zero value otherwise.

### GetLicenseExpressionOk

`func (o *ListProjectComponentsResponseItem) GetLicenseExpressionOk() (*string, bool)`

GetLicenseExpressionOk returns a tuple with the LicenseExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseExpression

`func (o *ListProjectComponentsResponseItem) SetLicenseExpression(v string)`

SetLicenseExpression sets LicenseExpression field to given value.

### HasLicenseExpression

`func (o *ListProjectComponentsResponseItem) HasLicenseExpression() bool`

HasLicenseExpression returns a boolean if a field has been set.

### GetLicenseUrl

`func (o *ListProjectComponentsResponseItem) GetLicenseUrl() string`

GetLicenseUrl returns the LicenseUrl field if non-nil, zero value otherwise.

### GetLicenseUrlOk

`func (o *ListProjectComponentsResponseItem) GetLicenseUrlOk() (*string, bool)`

GetLicenseUrlOk returns a tuple with the LicenseUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLicenseUrl

`func (o *ListProjectComponentsResponseItem) SetLicenseUrl(v string)`

SetLicenseUrl sets LicenseUrl field to given value.

### HasLicenseUrl

`func (o *ListProjectComponentsResponseItem) HasLicenseUrl() bool`

HasLicenseUrl returns a boolean if a field has been set.

### GetResolvedLicense

`func (o *ListProjectComponentsResponseItem) GetResolvedLicense() License`

GetResolvedLicense returns the ResolvedLicense field if non-nil, zero value otherwise.

### GetResolvedLicenseOk

`func (o *ListProjectComponentsResponseItem) GetResolvedLicenseOk() (*License, bool)`

GetResolvedLicenseOk returns a tuple with the ResolvedLicense field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolvedLicense

`func (o *ListProjectComponentsResponseItem) SetResolvedLicense(v License)`

SetResolvedLicense sets ResolvedLicense field to given value.

### HasResolvedLicense

`func (o *ListProjectComponentsResponseItem) HasResolvedLicense() bool`

HasResolvedLicense returns a boolean if a field has been set.

### GetOccurrenceCount

`func (o *ListProjectComponentsResponseItem) GetOccurrenceCount() int64`

GetOccurrenceCount returns the OccurrenceCount field if non-nil, zero value otherwise.

### GetOccurrenceCountOk

`func (o *ListProjectComponentsResponseItem) GetOccurrenceCountOk() (*int64, bool)`

GetOccurrenceCountOk returns a tuple with the OccurrenceCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOccurrenceCount

`func (o *ListProjectComponentsResponseItem) SetOccurrenceCount(v int64)`

SetOccurrenceCount sets OccurrenceCount field to given value.

### HasOccurrenceCount

`func (o *ListProjectComponentsResponseItem) HasOccurrenceCount() bool`

HasOccurrenceCount returns a boolean if a field has been set.

### GetScope

`func (o *ListProjectComponentsResponseItem) GetScope() Scope`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *ListProjectComponentsResponseItem) GetScopeOk() (*Scope, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *ListProjectComponentsResponseItem) SetScope(v Scope)`

SetScope sets Scope field to given value.

### HasScope

`func (o *ListProjectComponentsResponseItem) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetLastInheritedRiskScore

`func (o *ListProjectComponentsResponseItem) GetLastInheritedRiskScore() float64`

GetLastInheritedRiskScore returns the LastInheritedRiskScore field if non-nil, zero value otherwise.

### GetLastInheritedRiskScoreOk

`func (o *ListProjectComponentsResponseItem) GetLastInheritedRiskScoreOk() (*float64, bool)`

GetLastInheritedRiskScoreOk returns a tuple with the LastInheritedRiskScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastInheritedRiskScore

`func (o *ListProjectComponentsResponseItem) SetLastInheritedRiskScore(v float64)`

SetLastInheritedRiskScore sets LastInheritedRiskScore field to given value.

### HasLastInheritedRiskScore

`func (o *ListProjectComponentsResponseItem) HasLastInheritedRiskScore() bool`

HasLastInheritedRiskScore returns a boolean if a field has been set.

### GetUuid

`func (o *ListProjectComponentsResponseItem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ListProjectComponentsResponseItem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ListProjectComponentsResponseItem) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetMetrics

`func (o *ListProjectComponentsResponseItem) GetMetrics() DependencyMetrics`

GetMetrics returns the Metrics field if non-nil, zero value otherwise.

### GetMetricsOk

`func (o *ListProjectComponentsResponseItem) GetMetricsOk() (*DependencyMetrics, bool)`

GetMetricsOk returns a tuple with the Metrics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetrics

`func (o *ListProjectComponentsResponseItem) SetMetrics(v DependencyMetrics)`

SetMetrics sets Metrics field to given value.

### HasMetrics

`func (o *ListProjectComponentsResponseItem) HasMetrics() bool`

HasMetrics returns a boolean if a field has been set.

### GetPackageMetadata

`func (o *ListProjectComponentsResponseItem) GetPackageMetadata() PackageMetadata`

GetPackageMetadata returns the PackageMetadata field if non-nil, zero value otherwise.

### GetPackageMetadataOk

`func (o *ListProjectComponentsResponseItem) GetPackageMetadataOk() (*PackageMetadata, bool)`

GetPackageMetadataOk returns a tuple with the PackageMetadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageMetadata

`func (o *ListProjectComponentsResponseItem) SetPackageMetadata(v PackageMetadata)`

SetPackageMetadata sets PackageMetadata field to given value.

### HasPackageMetadata

`func (o *ListProjectComponentsResponseItem) HasPackageMetadata() bool`

HasPackageMetadata returns a boolean if a field has been set.

### GetPackageArtifactMetadata

`func (o *ListProjectComponentsResponseItem) GetPackageArtifactMetadata() PackageArtifactMetadata`

GetPackageArtifactMetadata returns the PackageArtifactMetadata field if non-nil, zero value otherwise.

### GetPackageArtifactMetadataOk

`func (o *ListProjectComponentsResponseItem) GetPackageArtifactMetadataOk() (*PackageArtifactMetadata, bool)`

GetPackageArtifactMetadataOk returns a tuple with the PackageArtifactMetadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageArtifactMetadata

`func (o *ListProjectComponentsResponseItem) SetPackageArtifactMetadata(v PackageArtifactMetadata)`

SetPackageArtifactMetadata sets PackageArtifactMetadata field to given value.

### HasPackageArtifactMetadata

`func (o *ListProjectComponentsResponseItem) HasPackageArtifactMetadata() bool`

HasPackageArtifactMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


