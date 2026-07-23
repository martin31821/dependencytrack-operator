# ListVulnPolicyBundlesResponseItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**Url** | **string** |  | 
**Hash** | Pointer to **string** |  | [optional] 
**LastSuccessfulSync** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**Created** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**Updated** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 

## Methods

### NewListVulnPolicyBundlesResponseItem

`func NewListVulnPolicyBundlesResponseItem(uuid string, url string, ) *ListVulnPolicyBundlesResponseItem`

NewListVulnPolicyBundlesResponseItem instantiates a new ListVulnPolicyBundlesResponseItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListVulnPolicyBundlesResponseItemWithDefaults

`func NewListVulnPolicyBundlesResponseItemWithDefaults() *ListVulnPolicyBundlesResponseItem`

NewListVulnPolicyBundlesResponseItemWithDefaults instantiates a new ListVulnPolicyBundlesResponseItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *ListVulnPolicyBundlesResponseItem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ListVulnPolicyBundlesResponseItem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ListVulnPolicyBundlesResponseItem) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetUrl

`func (o *ListVulnPolicyBundlesResponseItem) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ListVulnPolicyBundlesResponseItem) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ListVulnPolicyBundlesResponseItem) SetUrl(v string)`

SetUrl sets Url field to given value.


### GetHash

`func (o *ListVulnPolicyBundlesResponseItem) GetHash() string`

GetHash returns the Hash field if non-nil, zero value otherwise.

### GetHashOk

`func (o *ListVulnPolicyBundlesResponseItem) GetHashOk() (*string, bool)`

GetHashOk returns a tuple with the Hash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHash

`func (o *ListVulnPolicyBundlesResponseItem) SetHash(v string)`

SetHash sets Hash field to given value.

### HasHash

`func (o *ListVulnPolicyBundlesResponseItem) HasHash() bool`

HasHash returns a boolean if a field has been set.

### GetLastSuccessfulSync

`func (o *ListVulnPolicyBundlesResponseItem) GetLastSuccessfulSync() int64`

GetLastSuccessfulSync returns the LastSuccessfulSync field if non-nil, zero value otherwise.

### GetLastSuccessfulSyncOk

`func (o *ListVulnPolicyBundlesResponseItem) GetLastSuccessfulSyncOk() (*int64, bool)`

GetLastSuccessfulSyncOk returns a tuple with the LastSuccessfulSync field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSuccessfulSync

`func (o *ListVulnPolicyBundlesResponseItem) SetLastSuccessfulSync(v int64)`

SetLastSuccessfulSync sets LastSuccessfulSync field to given value.

### HasLastSuccessfulSync

`func (o *ListVulnPolicyBundlesResponseItem) HasLastSuccessfulSync() bool`

HasLastSuccessfulSync returns a boolean if a field has been set.

### GetCreated

`func (o *ListVulnPolicyBundlesResponseItem) GetCreated() int64`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *ListVulnPolicyBundlesResponseItem) GetCreatedOk() (*int64, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *ListVulnPolicyBundlesResponseItem) SetCreated(v int64)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *ListVulnPolicyBundlesResponseItem) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *ListVulnPolicyBundlesResponseItem) GetUpdated() int64`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *ListVulnPolicyBundlesResponseItem) GetUpdatedOk() (*int64, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *ListVulnPolicyBundlesResponseItem) SetUpdated(v int64)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *ListVulnPolicyBundlesResponseItem) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


