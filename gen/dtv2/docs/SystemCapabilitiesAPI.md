# \SystemCapabilitiesAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetSystemCapabilities**](SystemCapabilitiesAPI.md#GetSystemCapabilities) | **Get** /internal/system-capabilities | Get system capabilities



## GetSystemCapabilities

> SystemCapabilitiesResponse GetSystemCapabilities(ctx).Execute()

Get system capabilities



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
	resp, r, err := apiClient.SystemCapabilitiesAPI.GetSystemCapabilities(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SystemCapabilitiesAPI.GetSystemCapabilities``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSystemCapabilities`: SystemCapabilitiesResponse
	fmt.Fprintf(os.Stdout, "Response from `SystemCapabilitiesAPI.GetSystemCapabilities`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetSystemCapabilitiesRequest struct via the builder pattern


### Return type

[**SystemCapabilitiesResponse**](SystemCapabilitiesResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

