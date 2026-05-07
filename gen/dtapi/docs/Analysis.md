# Analysis

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AnalysisComments** | Pointer to [**[]AnalysisComment**](AnalysisComment.md) |  | [optional] 
**AnalysisDetails** | **string** |  | 
**AnalysisJustification** | **string** |  | 
**AnalysisResponse** | **string** |  | 
**AnalysisState** | **string** |  | 
**CvssV2Score** | Pointer to **float32** |  | [optional] 
**CvssV2Vector** | Pointer to **string** |  | [optional] 
**CvssV3Score** | Pointer to **float32** |  | [optional] 
**CvssV3Vector** | Pointer to **string** |  | [optional] 
**CvssV4Score** | Pointer to **float32** |  | [optional] 
**CvssV4Vector** | Pointer to **string** |  | [optional] 
**IsSuppressed** | Pointer to **bool** |  | [optional] 
**OwaspScore** | Pointer to **float32** |  | [optional] 
**OwaspVector** | Pointer to **string** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 

## Methods

### NewAnalysis

`func NewAnalysis(analysisDetails string, analysisJustification string, analysisResponse string, analysisState string, ) *Analysis`

NewAnalysis instantiates a new Analysis object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnalysisWithDefaults

`func NewAnalysisWithDefaults() *Analysis`

NewAnalysisWithDefaults instantiates a new Analysis object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnalysisComments

`func (o *Analysis) GetAnalysisComments() []AnalysisComment`

GetAnalysisComments returns the AnalysisComments field if non-nil, zero value otherwise.

### GetAnalysisCommentsOk

`func (o *Analysis) GetAnalysisCommentsOk() (*[]AnalysisComment, bool)`

GetAnalysisCommentsOk returns a tuple with the AnalysisComments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysisComments

`func (o *Analysis) SetAnalysisComments(v []AnalysisComment)`

SetAnalysisComments sets AnalysisComments field to given value.

### HasAnalysisComments

`func (o *Analysis) HasAnalysisComments() bool`

HasAnalysisComments returns a boolean if a field has been set.

### GetAnalysisDetails

`func (o *Analysis) GetAnalysisDetails() string`

GetAnalysisDetails returns the AnalysisDetails field if non-nil, zero value otherwise.

### GetAnalysisDetailsOk

`func (o *Analysis) GetAnalysisDetailsOk() (*string, bool)`

GetAnalysisDetailsOk returns a tuple with the AnalysisDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysisDetails

`func (o *Analysis) SetAnalysisDetails(v string)`

SetAnalysisDetails sets AnalysisDetails field to given value.


### GetAnalysisJustification

`func (o *Analysis) GetAnalysisJustification() string`

GetAnalysisJustification returns the AnalysisJustification field if non-nil, zero value otherwise.

### GetAnalysisJustificationOk

`func (o *Analysis) GetAnalysisJustificationOk() (*string, bool)`

GetAnalysisJustificationOk returns a tuple with the AnalysisJustification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysisJustification

`func (o *Analysis) SetAnalysisJustification(v string)`

SetAnalysisJustification sets AnalysisJustification field to given value.


### GetAnalysisResponse

`func (o *Analysis) GetAnalysisResponse() string`

GetAnalysisResponse returns the AnalysisResponse field if non-nil, zero value otherwise.

### GetAnalysisResponseOk

`func (o *Analysis) GetAnalysisResponseOk() (*string, bool)`

GetAnalysisResponseOk returns a tuple with the AnalysisResponse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysisResponse

`func (o *Analysis) SetAnalysisResponse(v string)`

SetAnalysisResponse sets AnalysisResponse field to given value.


### GetAnalysisState

`func (o *Analysis) GetAnalysisState() string`

GetAnalysisState returns the AnalysisState field if non-nil, zero value otherwise.

### GetAnalysisStateOk

`func (o *Analysis) GetAnalysisStateOk() (*string, bool)`

GetAnalysisStateOk returns a tuple with the AnalysisState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnalysisState

`func (o *Analysis) SetAnalysisState(v string)`

SetAnalysisState sets AnalysisState field to given value.


### GetCvssV2Score

`func (o *Analysis) GetCvssV2Score() float32`

GetCvssV2Score returns the CvssV2Score field if non-nil, zero value otherwise.

### GetCvssV2ScoreOk

`func (o *Analysis) GetCvssV2ScoreOk() (*float32, bool)`

GetCvssV2ScoreOk returns a tuple with the CvssV2Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV2Score

`func (o *Analysis) SetCvssV2Score(v float32)`

SetCvssV2Score sets CvssV2Score field to given value.

### HasCvssV2Score

`func (o *Analysis) HasCvssV2Score() bool`

HasCvssV2Score returns a boolean if a field has been set.

### GetCvssV2Vector

`func (o *Analysis) GetCvssV2Vector() string`

GetCvssV2Vector returns the CvssV2Vector field if non-nil, zero value otherwise.

### GetCvssV2VectorOk

`func (o *Analysis) GetCvssV2VectorOk() (*string, bool)`

GetCvssV2VectorOk returns a tuple with the CvssV2Vector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV2Vector

`func (o *Analysis) SetCvssV2Vector(v string)`

SetCvssV2Vector sets CvssV2Vector field to given value.

### HasCvssV2Vector

`func (o *Analysis) HasCvssV2Vector() bool`

HasCvssV2Vector returns a boolean if a field has been set.

### GetCvssV3Score

`func (o *Analysis) GetCvssV3Score() float32`

GetCvssV3Score returns the CvssV3Score field if non-nil, zero value otherwise.

### GetCvssV3ScoreOk

`func (o *Analysis) GetCvssV3ScoreOk() (*float32, bool)`

GetCvssV3ScoreOk returns a tuple with the CvssV3Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV3Score

`func (o *Analysis) SetCvssV3Score(v float32)`

SetCvssV3Score sets CvssV3Score field to given value.

### HasCvssV3Score

`func (o *Analysis) HasCvssV3Score() bool`

HasCvssV3Score returns a boolean if a field has been set.

### GetCvssV3Vector

`func (o *Analysis) GetCvssV3Vector() string`

GetCvssV3Vector returns the CvssV3Vector field if non-nil, zero value otherwise.

### GetCvssV3VectorOk

`func (o *Analysis) GetCvssV3VectorOk() (*string, bool)`

GetCvssV3VectorOk returns a tuple with the CvssV3Vector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV3Vector

`func (o *Analysis) SetCvssV3Vector(v string)`

SetCvssV3Vector sets CvssV3Vector field to given value.

### HasCvssV3Vector

`func (o *Analysis) HasCvssV3Vector() bool`

HasCvssV3Vector returns a boolean if a field has been set.

### GetCvssV4Score

`func (o *Analysis) GetCvssV4Score() float32`

GetCvssV4Score returns the CvssV4Score field if non-nil, zero value otherwise.

### GetCvssV4ScoreOk

`func (o *Analysis) GetCvssV4ScoreOk() (*float32, bool)`

GetCvssV4ScoreOk returns a tuple with the CvssV4Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV4Score

`func (o *Analysis) SetCvssV4Score(v float32)`

SetCvssV4Score sets CvssV4Score field to given value.

### HasCvssV4Score

`func (o *Analysis) HasCvssV4Score() bool`

HasCvssV4Score returns a boolean if a field has been set.

### GetCvssV4Vector

`func (o *Analysis) GetCvssV4Vector() string`

GetCvssV4Vector returns the CvssV4Vector field if non-nil, zero value otherwise.

### GetCvssV4VectorOk

`func (o *Analysis) GetCvssV4VectorOk() (*string, bool)`

GetCvssV4VectorOk returns a tuple with the CvssV4Vector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvssV4Vector

`func (o *Analysis) SetCvssV4Vector(v string)`

SetCvssV4Vector sets CvssV4Vector field to given value.

### HasCvssV4Vector

`func (o *Analysis) HasCvssV4Vector() bool`

HasCvssV4Vector returns a boolean if a field has been set.

### GetIsSuppressed

`func (o *Analysis) GetIsSuppressed() bool`

GetIsSuppressed returns the IsSuppressed field if non-nil, zero value otherwise.

### GetIsSuppressedOk

`func (o *Analysis) GetIsSuppressedOk() (*bool, bool)`

GetIsSuppressedOk returns a tuple with the IsSuppressed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSuppressed

`func (o *Analysis) SetIsSuppressed(v bool)`

SetIsSuppressed sets IsSuppressed field to given value.

### HasIsSuppressed

`func (o *Analysis) HasIsSuppressed() bool`

HasIsSuppressed returns a boolean if a field has been set.

### GetOwaspScore

`func (o *Analysis) GetOwaspScore() float32`

GetOwaspScore returns the OwaspScore field if non-nil, zero value otherwise.

### GetOwaspScoreOk

`func (o *Analysis) GetOwaspScoreOk() (*float32, bool)`

GetOwaspScoreOk returns a tuple with the OwaspScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwaspScore

`func (o *Analysis) SetOwaspScore(v float32)`

SetOwaspScore sets OwaspScore field to given value.

### HasOwaspScore

`func (o *Analysis) HasOwaspScore() bool`

HasOwaspScore returns a boolean if a field has been set.

### GetOwaspVector

`func (o *Analysis) GetOwaspVector() string`

GetOwaspVector returns the OwaspVector field if non-nil, zero value otherwise.

### GetOwaspVectorOk

`func (o *Analysis) GetOwaspVectorOk() (*string, bool)`

GetOwaspVectorOk returns a tuple with the OwaspVector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwaspVector

`func (o *Analysis) SetOwaspVector(v string)`

SetOwaspVector sets OwaspVector field to given value.

### HasOwaspVector

`func (o *Analysis) HasOwaspVector() bool`

HasOwaspVector returns a boolean if a field has been set.

### GetSeverity

`func (o *Analysis) GetSeverity() string`

GetSeverity returns the Severity field if non-nil, zero value otherwise.

### GetSeverityOk

`func (o *Analysis) GetSeverityOk() (*string, bool)`

GetSeverityOk returns a tuple with the Severity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeverity

`func (o *Analysis) SetSeverity(v string)`

SetSeverity sets Severity field to given value.

### HasSeverity

`func (o *Analysis) HasSeverity() bool`

HasSeverity returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


