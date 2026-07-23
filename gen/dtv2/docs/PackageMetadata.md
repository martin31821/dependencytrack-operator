# PackageMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LatestVersion** | Pointer to **NullableString** | Latest known version in the configured registries. Null when the resolver ran but returned no version.  | [optional] 
**LatestVersionPublishedAt** | Pointer to **NullableInt64** | When the latest version was published. May be null even when latest_version is non-null, as some registries do not report publication dates.  | [optional] 
**ResolvedAt** | **int64** | When package metadata was last resolved from the upstream repository. | 

## Methods

### NewPackageMetadata

`func NewPackageMetadata(resolvedAt int64, ) *PackageMetadata`

NewPackageMetadata instantiates a new PackageMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPackageMetadataWithDefaults

`func NewPackageMetadataWithDefaults() *PackageMetadata`

NewPackageMetadataWithDefaults instantiates a new PackageMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLatestVersion

`func (o *PackageMetadata) GetLatestVersion() string`

GetLatestVersion returns the LatestVersion field if non-nil, zero value otherwise.

### GetLatestVersionOk

`func (o *PackageMetadata) GetLatestVersionOk() (*string, bool)`

GetLatestVersionOk returns a tuple with the LatestVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLatestVersion

`func (o *PackageMetadata) SetLatestVersion(v string)`

SetLatestVersion sets LatestVersion field to given value.

### HasLatestVersion

`func (o *PackageMetadata) HasLatestVersion() bool`

HasLatestVersion returns a boolean if a field has been set.

### SetLatestVersionNil

`func (o *PackageMetadata) SetLatestVersionNil(b bool)`

 SetLatestVersionNil sets the value for LatestVersion to be an explicit nil

### UnsetLatestVersion
`func (o *PackageMetadata) UnsetLatestVersion()`

UnsetLatestVersion ensures that no value is present for LatestVersion, not even an explicit nil
### GetLatestVersionPublishedAt

`func (o *PackageMetadata) GetLatestVersionPublishedAt() int64`

GetLatestVersionPublishedAt returns the LatestVersionPublishedAt field if non-nil, zero value otherwise.

### GetLatestVersionPublishedAtOk

`func (o *PackageMetadata) GetLatestVersionPublishedAtOk() (*int64, bool)`

GetLatestVersionPublishedAtOk returns a tuple with the LatestVersionPublishedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLatestVersionPublishedAt

`func (o *PackageMetadata) SetLatestVersionPublishedAt(v int64)`

SetLatestVersionPublishedAt sets LatestVersionPublishedAt field to given value.

### HasLatestVersionPublishedAt

`func (o *PackageMetadata) HasLatestVersionPublishedAt() bool`

HasLatestVersionPublishedAt returns a boolean if a field has been set.

### SetLatestVersionPublishedAtNil

`func (o *PackageMetadata) SetLatestVersionPublishedAtNil(b bool)`

 SetLatestVersionPublishedAtNil sets the value for LatestVersionPublishedAt to be an explicit nil

### UnsetLatestVersionPublishedAt
`func (o *PackageMetadata) UnsetLatestVersionPublishedAt()`

UnsetLatestVersionPublishedAt ensures that no value is present for LatestVersionPublishedAt, not even an explicit nil
### GetResolvedAt

`func (o *PackageMetadata) GetResolvedAt() int64`

GetResolvedAt returns the ResolvedAt field if non-nil, zero value otherwise.

### GetResolvedAtOk

`func (o *PackageMetadata) GetResolvedAtOk() (*int64, bool)`

GetResolvedAtOk returns a tuple with the ResolvedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolvedAt

`func (o *PackageMetadata) SetResolvedAt(v int64)`

SetResolvedAt sets ResolvedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


