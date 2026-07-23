# \VulnPoliciesAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateVulnPolicy**](VulnPoliciesAPI.md#CreateVulnPolicy) | **Post** /vuln-policies | Create a vulnerability policy
[**DeleteVulnPolicy**](VulnPoliciesAPI.md#DeleteVulnPolicy) | **Delete** /vuln-policies/{uuid} | Delete a vulnerability policy
[**DeleteVulnPolicyBundle**](VulnPoliciesAPI.md#DeleteVulnPolicyBundle) | **Delete** /vuln-policy-bundles/{uuid} | Delete a vulnerability policy bundle
[**GetLatestVulnPolicyBundleSyncRun**](VulnPoliciesAPI.md#GetLatestVulnPolicyBundleSyncRun) | **Get** /vuln-policy-bundles/{uuid}/sync-runs/latest | Get the latest vulnerability policy bundle sync run
[**GetVulnPolicy**](VulnPoliciesAPI.md#GetVulnPolicy) | **Get** /vuln-policies/{uuid} | Get a vulnerability policy
[**ListVulnPolicies**](VulnPoliciesAPI.md#ListVulnPolicies) | **Get** /vuln-policies | List vulnerability policies
[**ListVulnPolicyBundles**](VulnPoliciesAPI.md#ListVulnPolicyBundles) | **Get** /vuln-policy-bundles | List vulnerability policy bundles
[**TriggerVulnPolicyBundleSyncRun**](VulnPoliciesAPI.md#TriggerVulnPolicyBundleSyncRun) | **Post** /vuln-policy-bundles/{uuid}/sync-runs | Trigger a vulnerability policy bundle sync run
[**UpdateVulnPolicy**](VulnPoliciesAPI.md#UpdateVulnPolicy) | **Put** /vuln-policies/{uuid} | Replace a vulnerability policy



## CreateVulnPolicy

> CreateVulnPolicy201Response CreateVulnPolicy(ctx).CreateVulnPolicyRequest(createVulnPolicyRequest).Execute()

Create a vulnerability policy



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createVulnPolicyRequest := *openapiclient.NewCreateVulnPolicyRequest("Name_example", "Condition_example", *openapiclient.NewVulnPolicyAnalysis("State_example")) // CreateVulnPolicyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnPoliciesAPI.CreateVulnPolicy(context.Background()).CreateVulnPolicyRequest(createVulnPolicyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.CreateVulnPolicy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateVulnPolicy`: CreateVulnPolicy201Response
	fmt.Fprintf(os.Stdout, "Response from `VulnPoliciesAPI.CreateVulnPolicy`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateVulnPolicyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createVulnPolicyRequest** | [**CreateVulnPolicyRequest**](CreateVulnPolicyRequest.md) |  | 

### Return type

[**CreateVulnPolicy201Response**](CreateVulnPolicy201Response.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVulnPolicy

> DeleteVulnPolicy(ctx, uuid).Execute()

Delete a vulnerability policy



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VulnPoliciesAPI.DeleteVulnPolicy(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.DeleteVulnPolicy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteVulnPolicyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVulnPolicyBundle

> DeleteVulnPolicyBundle(ctx, uuid).Execute()

Delete a vulnerability policy bundle



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy bundle

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VulnPoliciesAPI.DeleteVulnPolicyBundle(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.DeleteVulnPolicyBundle``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy bundle | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteVulnPolicyBundleRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLatestVulnPolicyBundleSyncRun

> VulnPolicyBundleSyncStatus GetLatestVulnPolicyBundleSyncRun(ctx, uuid).Execute()

Get the latest vulnerability policy bundle sync run



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy bundle

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnPoliciesAPI.GetLatestVulnPolicyBundleSyncRun(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.GetLatestVulnPolicyBundleSyncRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLatestVulnPolicyBundleSyncRun`: VulnPolicyBundleSyncStatus
	fmt.Fprintf(os.Stdout, "Response from `VulnPoliciesAPI.GetLatestVulnPolicyBundleSyncRun`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy bundle | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLatestVulnPolicyBundleSyncRunRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**VulnPolicyBundleSyncStatus**](VulnPolicyBundleSyncStatus.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVulnPolicy

> GetVulnPolicyResponse GetVulnPolicy(ctx, uuid).Execute()

Get a vulnerability policy



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnPoliciesAPI.GetVulnPolicy(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.GetVulnPolicy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVulnPolicy`: GetVulnPolicyResponse
	fmt.Fprintf(os.Stdout, "Response from `VulnPoliciesAPI.GetVulnPolicy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVulnPolicyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetVulnPolicyResponse**](GetVulnPolicyResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVulnPolicies

> ListVulnPoliciesResponse ListVulnPolicies(ctx).Limit(limit).PageToken(pageToken).Name(name).Execute()

List vulnerability policies



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)
	name := "name_example" // string | Filter by name (partial match, case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnPoliciesAPI.ListVulnPolicies(context.Background()).Limit(limit).PageToken(pageToken).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.ListVulnPolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVulnPolicies`: ListVulnPoliciesResponse
	fmt.Fprintf(os.Stdout, "Response from `VulnPoliciesAPI.ListVulnPolicies`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVulnPoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 
 **name** | **string** | Filter by name (partial match, case-insensitive) | 

### Return type

[**ListVulnPoliciesResponse**](ListVulnPoliciesResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVulnPolicyBundles

> ListVulnPolicyBundlesResponse ListVulnPolicyBundles(ctx).Execute()

List vulnerability policy bundles



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnPoliciesAPI.ListVulnPolicyBundles(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.ListVulnPolicyBundles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListVulnPolicyBundles`: ListVulnPolicyBundlesResponse
	fmt.Fprintf(os.Stdout, "Response from `VulnPoliciesAPI.ListVulnPolicyBundles`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListVulnPolicyBundlesRequest struct via the builder pattern


### Return type

[**ListVulnPolicyBundlesResponse**](ListVulnPolicyBundlesResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TriggerVulnPolicyBundleSyncRun

> TriggerVulnPolicyBundleSyncRun(ctx, uuid).Execute()

Trigger a vulnerability policy bundle sync run



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy bundle

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VulnPoliciesAPI.TriggerVulnPolicyBundleSyncRun(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.TriggerVulnPolicyBundleSyncRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy bundle | 

### Other Parameters

Other parameters are passed through a pointer to a apiTriggerVulnPolicyBundleSyncRunRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVulnPolicy

> UpdateVulnPolicy(ctx, uuid).UpdateVulnPolicyRequest(updateVulnPolicyRequest).Execute()

Replace a vulnerability policy



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | UUID of the vulnerability policy
	updateVulnPolicyRequest := *openapiclient.NewUpdateVulnPolicyRequest("Name_example", "Condition_example", *openapiclient.NewVulnPolicyAnalysis("State_example")) // UpdateVulnPolicyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VulnPoliciesAPI.UpdateVulnPolicy(context.Background(), uuid).UpdateVulnPolicyRequest(updateVulnPolicyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnPoliciesAPI.UpdateVulnPolicy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | UUID of the vulnerability policy | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVulnPolicyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateVulnPolicyRequest** | [**UpdateVulnPolicyRequest**](UpdateVulnPolicyRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

