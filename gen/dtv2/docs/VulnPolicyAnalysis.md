# VulnPolicyAnalysis

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**State** | **string** |  | 
**Justification** | Pointer to **string** |  | [optional] 
**VendorResponse** | Pointer to **string** |  | [optional] 
**Details** | Pointer to **string** |  | [optional] 
**Suppress** | Pointer to **bool** |  | [optional] [default to false]

## Methods

### NewVulnPolicyAnalysis

`func NewVulnPolicyAnalysis(state string, ) *VulnPolicyAnalysis`

NewVulnPolicyAnalysis instantiates a new VulnPolicyAnalysis object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVulnPolicyAnalysisWithDefaults

`func NewVulnPolicyAnalysisWithDefaults() *VulnPolicyAnalysis`

NewVulnPolicyAnalysisWithDefaults instantiates a new VulnPolicyAnalysis object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetState

`func (o *VulnPolicyAnalysis) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *VulnPolicyAnalysis) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *VulnPolicyAnalysis) SetState(v string)`

SetState sets State field to given value.


### GetJustification

`func (o *VulnPolicyAnalysis) GetJustification() string`

GetJustification returns the Justification field if non-nil, zero value otherwise.

### GetJustificationOk

`func (o *VulnPolicyAnalysis) GetJustificationOk() (*string, bool)`

GetJustificationOk returns a tuple with the Justification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJustification

`func (o *VulnPolicyAnalysis) SetJustification(v string)`

SetJustification sets Justification field to given value.

### HasJustification

`func (o *VulnPolicyAnalysis) HasJustification() bool`

HasJustification returns a boolean if a field has been set.

### GetVendorResponse

`func (o *VulnPolicyAnalysis) GetVendorResponse() string`

GetVendorResponse returns the VendorResponse field if non-nil, zero value otherwise.

### GetVendorResponseOk

`func (o *VulnPolicyAnalysis) GetVendorResponseOk() (*string, bool)`

GetVendorResponseOk returns a tuple with the VendorResponse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendorResponse

`func (o *VulnPolicyAnalysis) SetVendorResponse(v string)`

SetVendorResponse sets VendorResponse field to given value.

### HasVendorResponse

`func (o *VulnPolicyAnalysis) HasVendorResponse() bool`

HasVendorResponse returns a boolean if a field has been set.

### GetDetails

`func (o *VulnPolicyAnalysis) GetDetails() string`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *VulnPolicyAnalysis) GetDetailsOk() (*string, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *VulnPolicyAnalysis) SetDetails(v string)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *VulnPolicyAnalysis) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### GetSuppress

`func (o *VulnPolicyAnalysis) GetSuppress() bool`

GetSuppress returns the Suppress field if non-nil, zero value otherwise.

### GetSuppressOk

`func (o *VulnPolicyAnalysis) GetSuppressOk() (*bool, bool)`

GetSuppressOk returns a tuple with the Suppress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuppress

`func (o *VulnPolicyAnalysis) SetSuppress(v bool)`

SetSuppress sets Suppress field to given value.

### HasSuppress

`func (o *VulnPolicyAnalysis) HasSuppress() bool`

HasSuppress returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


