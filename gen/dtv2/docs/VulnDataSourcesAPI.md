# \VulnDataSourcesAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetLatestVulnDataSourceMirrorRun**](VulnDataSourcesAPI.md#GetLatestVulnDataSourceMirrorRun) | **Get** /vuln-data-sources/{name}/mirror-runs/latest | Get the latest vulnerability data source mirror run
[**TriggerVulnDataSourceMirrorRun**](VulnDataSourcesAPI.md#TriggerVulnDataSourceMirrorRun) | **Post** /vuln-data-sources/{name}/mirror-runs | Trigger a vulnerability data source mirror run



## GetLatestVulnDataSourceMirrorRun

> VulnDataSourceMirrorStatus GetLatestVulnDataSourceMirrorRun(ctx, name).Execute()

Get the latest vulnerability data source mirror run



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
	name := "name_example" // string | Name of the vulnerability data source (e.g. `nvd`, `osv`, `github`).

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VulnDataSourcesAPI.GetLatestVulnDataSourceMirrorRun(context.Background(), name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnDataSourcesAPI.GetLatestVulnDataSourceMirrorRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLatestVulnDataSourceMirrorRun`: VulnDataSourceMirrorStatus
	fmt.Fprintf(os.Stdout, "Response from `VulnDataSourcesAPI.GetLatestVulnDataSourceMirrorRun`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | Name of the vulnerability data source (e.g. &#x60;nvd&#x60;, &#x60;osv&#x60;, &#x60;github&#x60;). | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLatestVulnDataSourceMirrorRunRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**VulnDataSourceMirrorStatus**](VulnDataSourceMirrorStatus.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TriggerVulnDataSourceMirrorRun

> TriggerVulnDataSourceMirrorRun(ctx, name).Execute()

Trigger a vulnerability data source mirror run



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
	name := "name_example" // string | Name of the vulnerability data source (e.g. `nvd`, `osv`, `github`).

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VulnDataSourcesAPI.TriggerVulnDataSourceMirrorRun(context.Background(), name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VulnDataSourcesAPI.TriggerVulnDataSourceMirrorRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | Name of the vulnerability data source (e.g. &#x60;nvd&#x60;, &#x60;osv&#x60;, &#x60;github&#x60;). | 

### Other Parameters

Other parameters are passed through a pointer to a apiTriggerVulnDataSourceMirrorRunRequest struct via the builder pattern


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

