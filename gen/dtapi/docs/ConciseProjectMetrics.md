# ConciseProjectMetrics

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Components** | **int32** | Total number of components | 
**Critical** | **int32** | Number of vulnerabilities with critical severity | 
**High** | **int32** | Number of vulnerabilities with high severity | 
**InheritedRiskScore** | **float64** | The inherited risk score | 
**Low** | **int32** | Number of vulnerabilities with low severity | 
**Medium** | **int32** | Number of vulnerabilities with medium severity | 
**PolicyViolationsFail** | **int32** | Number of policy violations with status FAIL | 
**PolicyViolationsInfo** | **int32** | Number of policy violations with status WARN | 
**PolicyViolationsLicenseTotal** | **int32** | Number of license policy violations | 
**PolicyViolationsOperationalTotal** | **int32** | Number of operational policy violations | 
**PolicyViolationsSecurityTotal** | **int32** | Number of security policy violations | 
**PolicyViolationsTotal** | **int32** | Total number of policy violations | 
**PolicyViolationsWarn** | **int32** | Number of policy violations with status WARN | 
**Unassigned** | **int32** | Number of vulnerabilities with unassigned severity | 
**Vulnerabilities** | **int32** | Total number of vulnerabilities | 

## Methods

### NewConciseProjectMetrics

`func NewConciseProjectMetrics(components int32, critical int32, high int32, inheritedRiskScore float64, low int32, medium int32, policyViolationsFail int32, policyViolationsInfo int32, policyViolationsLicenseTotal int32, policyViolationsOperationalTotal int32, policyViolationsSecurityTotal int32, policyViolationsTotal int32, policyViolationsWarn int32, unassigned int32, vulnerabilities int32, ) *ConciseProjectMetrics`

NewConciseProjectMetrics instantiates a new ConciseProjectMetrics object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConciseProjectMetricsWithDefaults

`func NewConciseProjectMetricsWithDefaults() *ConciseProjectMetrics`

NewConciseProjectMetricsWithDefaults instantiates a new ConciseProjectMetrics object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetComponents

`func (o *ConciseProjectMetrics) GetComponents() int32`

GetComponents returns the Components field if non-nil, zero value otherwise.

### GetComponentsOk

`func (o *ConciseProjectMetrics) GetComponentsOk() (*int32, bool)`

GetComponentsOk returns a tuple with the Components field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponents

`func (o *ConciseProjectMetrics) SetComponents(v int32)`

SetComponents sets Components field to given value.


### GetCritical

`func (o *ConciseProjectMetrics) GetCritical() int32`

GetCritical returns the Critical field if non-nil, zero value otherwise.

### GetCriticalOk

`func (o *ConciseProjectMetrics) GetCriticalOk() (*int32, bool)`

GetCriticalOk returns a tuple with the Critical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCritical

`func (o *ConciseProjectMetrics) SetCritical(v int32)`

SetCritical sets Critical field to given value.


### GetHigh

`func (o *ConciseProjectMetrics) GetHigh() int32`

GetHigh returns the High field if non-nil, zero value otherwise.

### GetHighOk

`func (o *ConciseProjectMetrics) GetHighOk() (*int32, bool)`

GetHighOk returns a tuple with the High field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHigh

`func (o *ConciseProjectMetrics) SetHigh(v int32)`

SetHigh sets High field to given value.


### GetInheritedRiskScore

`func (o *ConciseProjectMetrics) GetInheritedRiskScore() float64`

GetInheritedRiskScore returns the InheritedRiskScore field if non-nil, zero value otherwise.

### GetInheritedRiskScoreOk

`func (o *ConciseProjectMetrics) GetInheritedRiskScoreOk() (*float64, bool)`

GetInheritedRiskScoreOk returns a tuple with the InheritedRiskScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInheritedRiskScore

`func (o *ConciseProjectMetrics) SetInheritedRiskScore(v float64)`

SetInheritedRiskScore sets InheritedRiskScore field to given value.


### GetLow

`func (o *ConciseProjectMetrics) GetLow() int32`

GetLow returns the Low field if non-nil, zero value otherwise.

### GetLowOk

`func (o *ConciseProjectMetrics) GetLowOk() (*int32, bool)`

GetLowOk returns a tuple with the Low field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLow

`func (o *ConciseProjectMetrics) SetLow(v int32)`

SetLow sets Low field to given value.


### GetMedium

`func (o *ConciseProjectMetrics) GetMedium() int32`

GetMedium returns the Medium field if non-nil, zero value otherwise.

### GetMediumOk

`func (o *ConciseProjectMetrics) GetMediumOk() (*int32, bool)`

GetMediumOk returns a tuple with the Medium field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMedium

`func (o *ConciseProjectMetrics) SetMedium(v int32)`

SetMedium sets Medium field to given value.


### GetPolicyViolationsFail

`func (o *ConciseProjectMetrics) GetPolicyViolationsFail() int32`

GetPolicyViolationsFail returns the PolicyViolationsFail field if non-nil, zero value otherwise.

### GetPolicyViolationsFailOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsFailOk() (*int32, bool)`

GetPolicyViolationsFailOk returns a tuple with the PolicyViolationsFail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsFail

`func (o *ConciseProjectMetrics) SetPolicyViolationsFail(v int32)`

SetPolicyViolationsFail sets PolicyViolationsFail field to given value.


### GetPolicyViolationsInfo

`func (o *ConciseProjectMetrics) GetPolicyViolationsInfo() int32`

GetPolicyViolationsInfo returns the PolicyViolationsInfo field if non-nil, zero value otherwise.

### GetPolicyViolationsInfoOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsInfoOk() (*int32, bool)`

GetPolicyViolationsInfoOk returns a tuple with the PolicyViolationsInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsInfo

`func (o *ConciseProjectMetrics) SetPolicyViolationsInfo(v int32)`

SetPolicyViolationsInfo sets PolicyViolationsInfo field to given value.


### GetPolicyViolationsLicenseTotal

`func (o *ConciseProjectMetrics) GetPolicyViolationsLicenseTotal() int32`

GetPolicyViolationsLicenseTotal returns the PolicyViolationsLicenseTotal field if non-nil, zero value otherwise.

### GetPolicyViolationsLicenseTotalOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsLicenseTotalOk() (*int32, bool)`

GetPolicyViolationsLicenseTotalOk returns a tuple with the PolicyViolationsLicenseTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsLicenseTotal

`func (o *ConciseProjectMetrics) SetPolicyViolationsLicenseTotal(v int32)`

SetPolicyViolationsLicenseTotal sets PolicyViolationsLicenseTotal field to given value.


### GetPolicyViolationsOperationalTotal

`func (o *ConciseProjectMetrics) GetPolicyViolationsOperationalTotal() int32`

GetPolicyViolationsOperationalTotal returns the PolicyViolationsOperationalTotal field if non-nil, zero value otherwise.

### GetPolicyViolationsOperationalTotalOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsOperationalTotalOk() (*int32, bool)`

GetPolicyViolationsOperationalTotalOk returns a tuple with the PolicyViolationsOperationalTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsOperationalTotal

`func (o *ConciseProjectMetrics) SetPolicyViolationsOperationalTotal(v int32)`

SetPolicyViolationsOperationalTotal sets PolicyViolationsOperationalTotal field to given value.


### GetPolicyViolationsSecurityTotal

`func (o *ConciseProjectMetrics) GetPolicyViolationsSecurityTotal() int32`

GetPolicyViolationsSecurityTotal returns the PolicyViolationsSecurityTotal field if non-nil, zero value otherwise.

### GetPolicyViolationsSecurityTotalOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsSecurityTotalOk() (*int32, bool)`

GetPolicyViolationsSecurityTotalOk returns a tuple with the PolicyViolationsSecurityTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsSecurityTotal

`func (o *ConciseProjectMetrics) SetPolicyViolationsSecurityTotal(v int32)`

SetPolicyViolationsSecurityTotal sets PolicyViolationsSecurityTotal field to given value.


### GetPolicyViolationsTotal

`func (o *ConciseProjectMetrics) GetPolicyViolationsTotal() int32`

GetPolicyViolationsTotal returns the PolicyViolationsTotal field if non-nil, zero value otherwise.

### GetPolicyViolationsTotalOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsTotalOk() (*int32, bool)`

GetPolicyViolationsTotalOk returns a tuple with the PolicyViolationsTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsTotal

`func (o *ConciseProjectMetrics) SetPolicyViolationsTotal(v int32)`

SetPolicyViolationsTotal sets PolicyViolationsTotal field to given value.


### GetPolicyViolationsWarn

`func (o *ConciseProjectMetrics) GetPolicyViolationsWarn() int32`

GetPolicyViolationsWarn returns the PolicyViolationsWarn field if non-nil, zero value otherwise.

### GetPolicyViolationsWarnOk

`func (o *ConciseProjectMetrics) GetPolicyViolationsWarnOk() (*int32, bool)`

GetPolicyViolationsWarnOk returns a tuple with the PolicyViolationsWarn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyViolationsWarn

`func (o *ConciseProjectMetrics) SetPolicyViolationsWarn(v int32)`

SetPolicyViolationsWarn sets PolicyViolationsWarn field to given value.


### GetUnassigned

`func (o *ConciseProjectMetrics) GetUnassigned() int32`

GetUnassigned returns the Unassigned field if non-nil, zero value otherwise.

### GetUnassignedOk

`func (o *ConciseProjectMetrics) GetUnassignedOk() (*int32, bool)`

GetUnassignedOk returns a tuple with the Unassigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnassigned

`func (o *ConciseProjectMetrics) SetUnassigned(v int32)`

SetUnassigned sets Unassigned field to given value.


### GetVulnerabilities

`func (o *ConciseProjectMetrics) GetVulnerabilities() int32`

GetVulnerabilities returns the Vulnerabilities field if non-nil, zero value otherwise.

### GetVulnerabilitiesOk

`func (o *ConciseProjectMetrics) GetVulnerabilitiesOk() (*int32, bool)`

GetVulnerabilitiesOk returns a tuple with the Vulnerabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVulnerabilities

`func (o *ConciseProjectMetrics) SetVulnerabilities(v int32)`

SetVulnerabilities sets Vulnerabilities field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


