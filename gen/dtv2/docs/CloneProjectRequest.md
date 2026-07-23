# CloneProjectRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | **string** | Version of the cloned project. | 
**VersionIsLatest** | Pointer to **bool** | Whether to mark the cloned project version as latest. If another version is already marked as latest, it will be atomically un-unmarked as part of the cloning operation. | [optional] [default to false]
**Includes** | Pointer to [**[]CloneProjectInclude**](CloneProjectInclude.md) | List of items to include in the clone:    * &#x60;ACL&#x60;: Include portfolio ACL definitions.   * &#x60;COMPONENTS&#x60;: Include components.   * &#x60;FINDINGS&#x60;: Include findings.       * Has no effect unless &#x60;COMPONENTS&#x60; is also included.   * &#x60;FINDINGS_AUDIT_HISTORY&#x60;: Include audit history of findings.       * Has no effect unless &#x60;FINDINGS&#x60; is also included.   * &#x60;POLICY_VIOLATIONS&#x60;: Include policy violations.       * Has no effect unless &#x60;COMPONENTS&#x60; is also included.   * &#x60;POLICY_VIOLATIONS_AUDIT_HISTORY&#x60;: Include audit history of policy violations.       * Has no effect unless &#x60;POLICY_VIOLATIONS&#x60; is also included.   * &#x60;PROPERTIES&#x60;: Include project properties.   * &#x60;SERVICES&#x60;: Include services.   * &#x60;TAGS&#x60;: Include project tags. | [optional] [default to {}]

## Methods

### NewCloneProjectRequest

`func NewCloneProjectRequest(version string, ) *CloneProjectRequest`

NewCloneProjectRequest instantiates a new CloneProjectRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCloneProjectRequestWithDefaults

`func NewCloneProjectRequestWithDefaults() *CloneProjectRequest`

NewCloneProjectRequestWithDefaults instantiates a new CloneProjectRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *CloneProjectRequest) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *CloneProjectRequest) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *CloneProjectRequest) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetVersionIsLatest

`func (o *CloneProjectRequest) GetVersionIsLatest() bool`

GetVersionIsLatest returns the VersionIsLatest field if non-nil, zero value otherwise.

### GetVersionIsLatestOk

`func (o *CloneProjectRequest) GetVersionIsLatestOk() (*bool, bool)`

GetVersionIsLatestOk returns a tuple with the VersionIsLatest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionIsLatest

`func (o *CloneProjectRequest) SetVersionIsLatest(v bool)`

SetVersionIsLatest sets VersionIsLatest field to given value.

### HasVersionIsLatest

`func (o *CloneProjectRequest) HasVersionIsLatest() bool`

HasVersionIsLatest returns a boolean if a field has been set.

### GetIncludes

`func (o *CloneProjectRequest) GetIncludes() []CloneProjectInclude`

GetIncludes returns the Includes field if non-nil, zero value otherwise.

### GetIncludesOk

`func (o *CloneProjectRequest) GetIncludesOk() (*[]CloneProjectInclude, bool)`

GetIncludesOk returns a tuple with the Includes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncludes

`func (o *CloneProjectRequest) SetIncludes(v []CloneProjectInclude)`

SetIncludes sets Includes field to given value.

### HasIncludes

`func (o *CloneProjectRequest) HasIncludes() bool`

HasIncludes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


