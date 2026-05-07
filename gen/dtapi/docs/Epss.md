# Epss

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cve** | **string** |  | 
**Percentile** | Pointer to **float32** |  | [optional] 
**Score** | Pointer to **float32** |  | [optional] 

## Methods

### NewEpss

`func NewEpss(cve string, ) *Epss`

NewEpss instantiates a new Epss object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEpssWithDefaults

`func NewEpssWithDefaults() *Epss`

NewEpssWithDefaults instantiates a new Epss object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCve

`func (o *Epss) GetCve() string`

GetCve returns the Cve field if non-nil, zero value otherwise.

### GetCveOk

`func (o *Epss) GetCveOk() (*string, bool)`

GetCveOk returns a tuple with the Cve field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCve

`func (o *Epss) SetCve(v string)`

SetCve sets Cve field to given value.


### GetPercentile

`func (o *Epss) GetPercentile() float32`

GetPercentile returns the Percentile field if non-nil, zero value otherwise.

### GetPercentileOk

`func (o *Epss) GetPercentileOk() (*float32, bool)`

GetPercentileOk returns a tuple with the Percentile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPercentile

`func (o *Epss) SetPercentile(v float32)`

SetPercentile sets Percentile field to given value.

### HasPercentile

`func (o *Epss) HasPercentile() bool`

HasPercentile returns a boolean if a field has been set.

### GetScore

`func (o *Epss) GetScore() float32`

GetScore returns the Score field if non-nil, zero value otherwise.

### GetScoreOk

`func (o *Epss) GetScoreOk() (*float32, bool)`

GetScoreOk returns a tuple with the Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScore

`func (o *Epss) SetScore(v float32)`

SetScore sets Score field to given value.

### HasScore

`func (o *Epss) HasScore() bool`

HasScore returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


