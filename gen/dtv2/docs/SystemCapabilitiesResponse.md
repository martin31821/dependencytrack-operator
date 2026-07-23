# SystemCapabilitiesResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Capabilities** | **map[string]map[string]interface{}** | Map of namespaces to capability flags. | 

## Methods

### NewSystemCapabilitiesResponse

`func NewSystemCapabilitiesResponse(capabilities map[string]map[string]interface{}, ) *SystemCapabilitiesResponse`

NewSystemCapabilitiesResponse instantiates a new SystemCapabilitiesResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSystemCapabilitiesResponseWithDefaults

`func NewSystemCapabilitiesResponseWithDefaults() *SystemCapabilitiesResponse`

NewSystemCapabilitiesResponseWithDefaults instantiates a new SystemCapabilitiesResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapabilities

`func (o *SystemCapabilitiesResponse) GetCapabilities() map[string]map[string]interface{}`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *SystemCapabilitiesResponse) GetCapabilitiesOk() (*map[string]map[string]interface{}, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *SystemCapabilitiesResponse) SetCapabilities(v map[string]map[string]interface{})`

SetCapabilities sets Capabilities field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


