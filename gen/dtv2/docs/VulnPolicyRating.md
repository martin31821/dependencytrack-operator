# VulnPolicyRating

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Method** | **string** |  | 
**Severity** | **string** |  | 
**Vector** | Pointer to **string** |  | [optional] 
**Score** | Pointer to **float64** |  | [optional] 

## Methods

### NewVulnPolicyRating

`func NewVulnPolicyRating(method string, severity string, ) *VulnPolicyRating`

NewVulnPolicyRating instantiates a new VulnPolicyRating object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVulnPolicyRatingWithDefaults

`func NewVulnPolicyRatingWithDefaults() *VulnPolicyRating`

NewVulnPolicyRatingWithDefaults instantiates a new VulnPolicyRating object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMethod

`func (o *VulnPolicyRating) GetMethod() string`

GetMethod returns the Method field if non-nil, zero value otherwise.

### GetMethodOk

`func (o *VulnPolicyRating) GetMethodOk() (*string, bool)`

GetMethodOk returns a tuple with the Method field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethod

`func (o *VulnPolicyRating) SetMethod(v string)`

SetMethod sets Method field to given value.


### GetSeverity

`func (o *VulnPolicyRating) GetSeverity() string`

GetSeverity returns the Severity field if non-nil, zero value otherwise.

### GetSeverityOk

`func (o *VulnPolicyRating) GetSeverityOk() (*string, bool)`

GetSeverityOk returns a tuple with the Severity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeverity

`func (o *VulnPolicyRating) SetSeverity(v string)`

SetSeverity sets Severity field to given value.


### GetVector

`func (o *VulnPolicyRating) GetVector() string`

GetVector returns the Vector field if non-nil, zero value otherwise.

### GetVectorOk

`func (o *VulnPolicyRating) GetVectorOk() (*string, bool)`

GetVectorOk returns a tuple with the Vector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVector

`func (o *VulnPolicyRating) SetVector(v string)`

SetVector sets Vector field to given value.

### HasVector

`func (o *VulnPolicyRating) HasVector() bool`

HasVector returns a boolean if a field has been set.

### GetScore

`func (o *VulnPolicyRating) GetScore() float64`

GetScore returns the Score field if non-nil, zero value otherwise.

### GetScoreOk

`func (o *VulnPolicyRating) GetScoreOk() (*float64, bool)`

GetScoreOk returns a tuple with the Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScore

`func (o *VulnPolicyRating) SetScore(v float64)`

SetScore sets Score field to given value.

### HasScore

`func (o *VulnPolicyRating) HasScore() bool`

HasScore returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


