# VulnPolicyBundleSyncStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | Status of the synchronization. | 
**StartedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**CompletedAt** | Pointer to **int64** | Epoch timestamp in milliseconds since January 1, 1970 UTC. | [optional] 
**FailureReason** | Pointer to **string** | Reason for why the synchronization failed. | [optional] 

## Methods

### NewVulnPolicyBundleSyncStatus

`func NewVulnPolicyBundleSyncStatus(status string, ) *VulnPolicyBundleSyncStatus`

NewVulnPolicyBundleSyncStatus instantiates a new VulnPolicyBundleSyncStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVulnPolicyBundleSyncStatusWithDefaults

`func NewVulnPolicyBundleSyncStatusWithDefaults() *VulnPolicyBundleSyncStatus`

NewVulnPolicyBundleSyncStatusWithDefaults instantiates a new VulnPolicyBundleSyncStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *VulnPolicyBundleSyncStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VulnPolicyBundleSyncStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VulnPolicyBundleSyncStatus) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetStartedAt

`func (o *VulnPolicyBundleSyncStatus) GetStartedAt() int64`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *VulnPolicyBundleSyncStatus) GetStartedAtOk() (*int64, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *VulnPolicyBundleSyncStatus) SetStartedAt(v int64)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *VulnPolicyBundleSyncStatus) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetCompletedAt

`func (o *VulnPolicyBundleSyncStatus) GetCompletedAt() int64`

GetCompletedAt returns the CompletedAt field if non-nil, zero value otherwise.

### GetCompletedAtOk

`func (o *VulnPolicyBundleSyncStatus) GetCompletedAtOk() (*int64, bool)`

GetCompletedAtOk returns a tuple with the CompletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedAt

`func (o *VulnPolicyBundleSyncStatus) SetCompletedAt(v int64)`

SetCompletedAt sets CompletedAt field to given value.

### HasCompletedAt

`func (o *VulnPolicyBundleSyncStatus) HasCompletedAt() bool`

HasCompletedAt returns a boolean if a field has been set.

### GetFailureReason

`func (o *VulnPolicyBundleSyncStatus) GetFailureReason() string`

GetFailureReason returns the FailureReason field if non-nil, zero value otherwise.

### GetFailureReasonOk

`func (o *VulnPolicyBundleSyncStatus) GetFailureReasonOk() (*string, bool)`

GetFailureReasonOk returns a tuple with the FailureReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailureReason

`func (o *VulnPolicyBundleSyncStatus) SetFailureReason(v string)`

SetFailureReason sets FailureReason field to given value.

### HasFailureReason

`func (o *VulnPolicyBundleSyncStatus) HasFailureReason() bool`

HasFailureReason returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


