# ComponentMetaInformation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IntegrityMatchStatus** | Pointer to **string** |  | [optional] 
**IntegrityRepoUrl** | Pointer to **string** |  | [optional] 
**LastFetched** | Pointer to **util.DTTime** |  | [optional] 
**PublishedDate** | Pointer to **util.DTTime** |  | [optional] 

## Methods

### NewComponentMetaInformation

`func NewComponentMetaInformation() *ComponentMetaInformation`

NewComponentMetaInformation instantiates a new ComponentMetaInformation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewComponentMetaInformationWithDefaults

`func NewComponentMetaInformationWithDefaults() *ComponentMetaInformation`

NewComponentMetaInformationWithDefaults instantiates a new ComponentMetaInformation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIntegrityMatchStatus

`func (o *ComponentMetaInformation) GetIntegrityMatchStatus() string`

GetIntegrityMatchStatus returns the IntegrityMatchStatus field if non-nil, zero value otherwise.

### GetIntegrityMatchStatusOk

`func (o *ComponentMetaInformation) GetIntegrityMatchStatusOk() (*string, bool)`

GetIntegrityMatchStatusOk returns a tuple with the IntegrityMatchStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrityMatchStatus

`func (o *ComponentMetaInformation) SetIntegrityMatchStatus(v string)`

SetIntegrityMatchStatus sets IntegrityMatchStatus field to given value.

### HasIntegrityMatchStatus

`func (o *ComponentMetaInformation) HasIntegrityMatchStatus() bool`

HasIntegrityMatchStatus returns a boolean if a field has been set.

### GetIntegrityRepoUrl

`func (o *ComponentMetaInformation) GetIntegrityRepoUrl() string`

GetIntegrityRepoUrl returns the IntegrityRepoUrl field if non-nil, zero value otherwise.

### GetIntegrityRepoUrlOk

`func (o *ComponentMetaInformation) GetIntegrityRepoUrlOk() (*string, bool)`

GetIntegrityRepoUrlOk returns a tuple with the IntegrityRepoUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrityRepoUrl

`func (o *ComponentMetaInformation) SetIntegrityRepoUrl(v string)`

SetIntegrityRepoUrl sets IntegrityRepoUrl field to given value.

### HasIntegrityRepoUrl

`func (o *ComponentMetaInformation) HasIntegrityRepoUrl() bool`

HasIntegrityRepoUrl returns a boolean if a field has been set.

### GetLastFetched

`func (o *ComponentMetaInformation) GetLastFetched() util.DTTime`

GetLastFetched returns the LastFetched field if non-nil, zero value otherwise.

### GetLastFetchedOk

`func (o *ComponentMetaInformation) GetLastFetchedOk() (*util.DTTime, bool)`

GetLastFetchedOk returns a tuple with the LastFetched field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastFetched

`func (o *ComponentMetaInformation) SetLastFetched(v util.DTTime)`

SetLastFetched sets LastFetched field to given value.

### HasLastFetched

`func (o *ComponentMetaInformation) HasLastFetched() bool`

HasLastFetched returns a boolean if a field has been set.

### GetPublishedDate

`func (o *ComponentMetaInformation) GetPublishedDate() util.DTTime`

GetPublishedDate returns the PublishedDate field if non-nil, zero value otherwise.

### GetPublishedDateOk

`func (o *ComponentMetaInformation) GetPublishedDateOk() (*util.DTTime, bool)`

GetPublishedDateOk returns a tuple with the PublishedDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishedDate

`func (o *ComponentMetaInformation) SetPublishedDate(v util.DTTime)`

SetPublishedDate sets PublishedDate field to given value.

### HasPublishedDate

`func (o *ComponentMetaInformation) HasPublishedDate() bool`

HasPublishedDate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


