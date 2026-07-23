# ConciseProject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | **bool** | Whether the project is active | 
**Classifier** | Pointer to **string** | Classifier of the project | [optional] 
**CollectionLogic** | Pointer to **string** | Collection logic for aggregating child metrics | [optional] 
**Group** | Pointer to **string** | Group or namespace of the project | [optional] 
**HasChildren** | **bool** | Whether the project has children | 
**IsLatest** | Pointer to **bool** | Whether the project version is latest | [optional] 
**LastBomImport** | Pointer to **int64** | Timestamp of the last BOM import | [optional] 
**LastBomImportFormat** | Pointer to **string** | Format of the last imported BOM | [optional] 
**LastRiskScore** | Pointer to **float64** | Last observed risk score | [optional] 
**Metrics** | Pointer to [**ConciseProjectMetrics**](ConciseProjectMetrics.md) |  | [optional] 
**Name** | **string** | Name of the project | 
**Tags** | Pointer to [**[]Tag**](Tag.md) | Tags associated with the project | [optional] 
**Teams** | Pointer to [**[]Team**](Team.md) | Teams associated with the project | [optional] 
**Uuid** | **string** | UUID of the project | 
**Version** | Pointer to **string** | Version of the project | [optional] 

## Methods

### NewConciseProject

`func NewConciseProject(active bool, hasChildren bool, name string, uuid string, ) *ConciseProject`

NewConciseProject instantiates a new ConciseProject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConciseProjectWithDefaults

`func NewConciseProjectWithDefaults() *ConciseProject`

NewConciseProjectWithDefaults instantiates a new ConciseProject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *ConciseProject) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *ConciseProject) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *ConciseProject) SetActive(v bool)`

SetActive sets Active field to given value.


### GetClassifier

`func (o *ConciseProject) GetClassifier() string`

GetClassifier returns the Classifier field if non-nil, zero value otherwise.

### GetClassifierOk

`func (o *ConciseProject) GetClassifierOk() (*string, bool)`

GetClassifierOk returns a tuple with the Classifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifier

`func (o *ConciseProject) SetClassifier(v string)`

SetClassifier sets Classifier field to given value.

### HasClassifier

`func (o *ConciseProject) HasClassifier() bool`

HasClassifier returns a boolean if a field has been set.

### GetCollectionLogic

`func (o *ConciseProject) GetCollectionLogic() string`

GetCollectionLogic returns the CollectionLogic field if non-nil, zero value otherwise.

### GetCollectionLogicOk

`func (o *ConciseProject) GetCollectionLogicOk() (*string, bool)`

GetCollectionLogicOk returns a tuple with the CollectionLogic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollectionLogic

`func (o *ConciseProject) SetCollectionLogic(v string)`

SetCollectionLogic sets CollectionLogic field to given value.

### HasCollectionLogic

`func (o *ConciseProject) HasCollectionLogic() bool`

HasCollectionLogic returns a boolean if a field has been set.

### GetGroup

`func (o *ConciseProject) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ConciseProject) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ConciseProject) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ConciseProject) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetHasChildren

`func (o *ConciseProject) GetHasChildren() bool`

GetHasChildren returns the HasChildren field if non-nil, zero value otherwise.

### GetHasChildrenOk

`func (o *ConciseProject) GetHasChildrenOk() (*bool, bool)`

GetHasChildrenOk returns a tuple with the HasChildren field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasChildren

`func (o *ConciseProject) SetHasChildren(v bool)`

SetHasChildren sets HasChildren field to given value.


### GetIsLatest

`func (o *ConciseProject) GetIsLatest() bool`

GetIsLatest returns the IsLatest field if non-nil, zero value otherwise.

### GetIsLatestOk

`func (o *ConciseProject) GetIsLatestOk() (*bool, bool)`

GetIsLatestOk returns a tuple with the IsLatest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLatest

`func (o *ConciseProject) SetIsLatest(v bool)`

SetIsLatest sets IsLatest field to given value.

### HasIsLatest

`func (o *ConciseProject) HasIsLatest() bool`

HasIsLatest returns a boolean if a field has been set.

### GetLastBomImport

`func (o *ConciseProject) GetLastBomImport() int64`

GetLastBomImport returns the LastBomImport field if non-nil, zero value otherwise.

### GetLastBomImportOk

`func (o *ConciseProject) GetLastBomImportOk() (*int64, bool)`

GetLastBomImportOk returns a tuple with the LastBomImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastBomImport

`func (o *ConciseProject) SetLastBomImport(v int64)`

SetLastBomImport sets LastBomImport field to given value.

### HasLastBomImport

`func (o *ConciseProject) HasLastBomImport() bool`

HasLastBomImport returns a boolean if a field has been set.

### GetLastBomImportFormat

`func (o *ConciseProject) GetLastBomImportFormat() string`

GetLastBomImportFormat returns the LastBomImportFormat field if non-nil, zero value otherwise.

### GetLastBomImportFormatOk

`func (o *ConciseProject) GetLastBomImportFormatOk() (*string, bool)`

GetLastBomImportFormatOk returns a tuple with the LastBomImportFormat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastBomImportFormat

`func (o *ConciseProject) SetLastBomImportFormat(v string)`

SetLastBomImportFormat sets LastBomImportFormat field to given value.

### HasLastBomImportFormat

`func (o *ConciseProject) HasLastBomImportFormat() bool`

HasLastBomImportFormat returns a boolean if a field has been set.

### GetLastRiskScore

`func (o *ConciseProject) GetLastRiskScore() float64`

GetLastRiskScore returns the LastRiskScore field if non-nil, zero value otherwise.

### GetLastRiskScoreOk

`func (o *ConciseProject) GetLastRiskScoreOk() (*float64, bool)`

GetLastRiskScoreOk returns a tuple with the LastRiskScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastRiskScore

`func (o *ConciseProject) SetLastRiskScore(v float64)`

SetLastRiskScore sets LastRiskScore field to given value.

### HasLastRiskScore

`func (o *ConciseProject) HasLastRiskScore() bool`

HasLastRiskScore returns a boolean if a field has been set.

### GetMetrics

`func (o *ConciseProject) GetMetrics() ConciseProjectMetrics`

GetMetrics returns the Metrics field if non-nil, zero value otherwise.

### GetMetricsOk

`func (o *ConciseProject) GetMetricsOk() (*ConciseProjectMetrics, bool)`

GetMetricsOk returns a tuple with the Metrics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetrics

`func (o *ConciseProject) SetMetrics(v ConciseProjectMetrics)`

SetMetrics sets Metrics field to given value.

### HasMetrics

`func (o *ConciseProject) HasMetrics() bool`

HasMetrics returns a boolean if a field has been set.

### GetName

`func (o *ConciseProject) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConciseProject) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConciseProject) SetName(v string)`

SetName sets Name field to given value.


### GetTags

`func (o *ConciseProject) GetTags() []Tag`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ConciseProject) GetTagsOk() (*[]Tag, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ConciseProject) SetTags(v []Tag)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ConciseProject) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetTeams

`func (o *ConciseProject) GetTeams() []Team`

GetTeams returns the Teams field if non-nil, zero value otherwise.

### GetTeamsOk

`func (o *ConciseProject) GetTeamsOk() (*[]Team, bool)`

GetTeamsOk returns a tuple with the Teams field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTeams

`func (o *ConciseProject) SetTeams(v []Team)`

SetTeams sets Teams field to given value.

### HasTeams

`func (o *ConciseProject) HasTeams() bool`

HasTeams returns a boolean if a field has been set.

### GetUuid

`func (o *ConciseProject) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ConciseProject) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ConciseProject) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetVersion

`func (o *ConciseProject) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ConciseProject) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ConciseProject) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ConciseProject) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


