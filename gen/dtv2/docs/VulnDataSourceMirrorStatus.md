# VulnDataSourceMirrorStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | Status of the mirror run. | 
**StartedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**CompletedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**FailureReason** | Pointer to **string** | Reason for why the mirror run failed. | [optional] 

## Methods

### NewVulnDataSourceMirrorStatus

`func NewVulnDataSourceMirrorStatus(status string, ) *VulnDataSourceMirrorStatus`

NewVulnDataSourceMirrorStatus instantiates a new VulnDataSourceMirrorStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVulnDataSourceMirrorStatusWithDefaults

`func NewVulnDataSourceMirrorStatusWithDefaults() *VulnDataSourceMirrorStatus`

NewVulnDataSourceMirrorStatusWithDefaults instantiates a new VulnDataSourceMirrorStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *VulnDataSourceMirrorStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VulnDataSourceMirrorStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VulnDataSourceMirrorStatus) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetStartedAt

`func (o *VulnDataSourceMirrorStatus) GetStartedAt() int64`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *VulnDataSourceMirrorStatus) GetStartedAtOk() (*int64, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *VulnDataSourceMirrorStatus) SetStartedAt(v int64)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *VulnDataSourceMirrorStatus) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetCompletedAt

`func (o *VulnDataSourceMirrorStatus) GetCompletedAt() int64`

GetCompletedAt returns the CompletedAt field if non-nil, zero value otherwise.

### GetCompletedAtOk

`func (o *VulnDataSourceMirrorStatus) GetCompletedAtOk() (*int64, bool)`

GetCompletedAtOk returns a tuple with the CompletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedAt

`func (o *VulnDataSourceMirrorStatus) SetCompletedAt(v int64)`

SetCompletedAt sets CompletedAt field to given value.

### HasCompletedAt

`func (o *VulnDataSourceMirrorStatus) HasCompletedAt() bool`

HasCompletedAt returns a boolean if a field has been set.

### GetFailureReason

`func (o *VulnDataSourceMirrorStatus) GetFailureReason() string`

GetFailureReason returns the FailureReason field if non-nil, zero value otherwise.

### GetFailureReasonOk

`func (o *VulnDataSourceMirrorStatus) GetFailureReasonOk() (*string, bool)`

GetFailureReasonOk returns a tuple with the FailureReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailureReason

`func (o *VulnDataSourceMirrorStatus) SetFailureReason(v string)`

SetFailureReason sets FailureReason field to given value.

### HasFailureReason

`func (o *VulnDataSourceMirrorStatus) HasFailureReason() bool`

HasFailureReason returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


