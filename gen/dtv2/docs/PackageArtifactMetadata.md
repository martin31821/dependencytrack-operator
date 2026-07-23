# PackageArtifactMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hashes** | Pointer to [**Hashes**](Hashes.md) |  | [optional] 
**PublishedAt** | Pointer to **NullableInt64** | When this artifact was published to the repository. | [optional] 
**ResolvedFrom** | Pointer to **NullableString** | Identifier of the repository from which artifact metadata was fetched. | [optional] 
**ResolvedAt** | Pointer to **NullableInt64** | When artifact metadata was last resolved from the upstream repository. | [optional] 

## Methods

### NewPackageArtifactMetadata

`func NewPackageArtifactMetadata() *PackageArtifactMetadata`

NewPackageArtifactMetadata instantiates a new PackageArtifactMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPackageArtifactMetadataWithDefaults

`func NewPackageArtifactMetadataWithDefaults() *PackageArtifactMetadata`

NewPackageArtifactMetadataWithDefaults instantiates a new PackageArtifactMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHashes

`func (o *PackageArtifactMetadata) GetHashes() Hashes`

GetHashes returns the Hashes field if non-nil, zero value otherwise.

### GetHashesOk

`func (o *PackageArtifactMetadata) GetHashesOk() (*Hashes, bool)`

GetHashesOk returns a tuple with the Hashes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHashes

`func (o *PackageArtifactMetadata) SetHashes(v Hashes)`

SetHashes sets Hashes field to given value.

### HasHashes

`func (o *PackageArtifactMetadata) HasHashes() bool`

HasHashes returns a boolean if a field has been set.

### GetPublishedAt

`func (o *PackageArtifactMetadata) GetPublishedAt() int64`

GetPublishedAt returns the PublishedAt field if non-nil, zero value otherwise.

### GetPublishedAtOk

`func (o *PackageArtifactMetadata) GetPublishedAtOk() (*int64, bool)`

GetPublishedAtOk returns a tuple with the PublishedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishedAt

`func (o *PackageArtifactMetadata) SetPublishedAt(v int64)`

SetPublishedAt sets PublishedAt field to given value.

### HasPublishedAt

`func (o *PackageArtifactMetadata) HasPublishedAt() bool`

HasPublishedAt returns a boolean if a field has been set.

### SetPublishedAtNil

`func (o *PackageArtifactMetadata) SetPublishedAtNil(b bool)`

 SetPublishedAtNil sets the value for PublishedAt to be an explicit nil

### UnsetPublishedAt
`func (o *PackageArtifactMetadata) UnsetPublishedAt()`

UnsetPublishedAt ensures that no value is present for PublishedAt, not even an explicit nil
### GetResolvedFrom

`func (o *PackageArtifactMetadata) GetResolvedFrom() string`

GetResolvedFrom returns the ResolvedFrom field if non-nil, zero value otherwise.

### GetResolvedFromOk

`func (o *PackageArtifactMetadata) GetResolvedFromOk() (*string, bool)`

GetResolvedFromOk returns a tuple with the ResolvedFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolvedFrom

`func (o *PackageArtifactMetadata) SetResolvedFrom(v string)`

SetResolvedFrom sets ResolvedFrom field to given value.

### HasResolvedFrom

`func (o *PackageArtifactMetadata) HasResolvedFrom() bool`

HasResolvedFrom returns a boolean if a field has been set.

### SetResolvedFromNil

`func (o *PackageArtifactMetadata) SetResolvedFromNil(b bool)`

 SetResolvedFromNil sets the value for ResolvedFrom to be an explicit nil

### UnsetResolvedFrom
`func (o *PackageArtifactMetadata) UnsetResolvedFrom()`

UnsetResolvedFrom ensures that no value is present for ResolvedFrom, not even an explicit nil
### GetResolvedAt

`func (o *PackageArtifactMetadata) GetResolvedAt() int64`

GetResolvedAt returns the ResolvedAt field if non-nil, zero value otherwise.

### GetResolvedAtOk

`func (o *PackageArtifactMetadata) GetResolvedAtOk() (*int64, bool)`

GetResolvedAtOk returns a tuple with the ResolvedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolvedAt

`func (o *PackageArtifactMetadata) SetResolvedAt(v int64)`

SetResolvedAt sets ResolvedAt field to given value.

### HasResolvedAt

`func (o *PackageArtifactMetadata) HasResolvedAt() bool`

HasResolvedAt returns a boolean if a field has been set.

### SetResolvedAtNil

`func (o *PackageArtifactMetadata) SetResolvedAtNil(b bool)`

 SetResolvedAtNil sets the value for ResolvedAt to be an explicit nil

### UnsetResolvedAt
`func (o *PackageArtifactMetadata) UnsetResolvedAt()`

UnsetResolvedAt ensures that no value is present for ResolvedAt, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


