# \ExtensionsAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetExtensionConfig**](ExtensionsAPI.md#GetExtensionConfig) | **Get** /extension-points/{extension_point_name}/extensions/{extension_name}/config | Get extension configuration
[**GetExtensionConfigSchema**](ExtensionsAPI.md#GetExtensionConfigSchema) | **Get** /extension-points/{extension_point_name}/extensions/{extension_name}/config-schema | Get extension configuration schema
[**ListExtensionPoints**](ExtensionsAPI.md#ListExtensionPoints) | **Get** /extension-points | List all extension points
[**ListExtensions**](ExtensionsAPI.md#ListExtensions) | **Get** /extension-points/{extension_point_name}/extensions | List all extensions
[**TestExtension**](ExtensionsAPI.md#TestExtension) | **Post** /extension-points/{extension_point_name}/extensions/{extension_name}/test | Test extension
[**UpdateExtensionConfig**](ExtensionsAPI.md#UpdateExtensionConfig) | **Put** /extension-points/{extension_point_name}/extensions/{extension_name}/config | Update extension configuration



## GetExtensionConfig

> GetExtensionConfigResponse GetExtensionConfig(ctx, extensionPointName, extensionName).Execute()

Get extension configuration



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
	extensionPointName := "extensionPointName_example" // string | Name of the extension point
	extensionName := "extensionName_example" // string | Name of the extension

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExtensionsAPI.GetExtensionConfig(context.Background(), extensionPointName, extensionName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.GetExtensionConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExtensionConfig`: GetExtensionConfigResponse
	fmt.Fprintf(os.Stdout, "Response from `ExtensionsAPI.GetExtensionConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**extensionPointName** | **string** | Name of the extension point | 
**extensionName** | **string** | Name of the extension | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExtensionConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetExtensionConfigResponse**](GetExtensionConfigResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExtensionConfigSchema

> ExtensionConfigSchema GetExtensionConfigSchema(ctx, extensionPointName, extensionName).Execute()

Get extension configuration schema



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
	extensionPointName := "extensionPointName_example" // string | Name of the extension point
	extensionName := "extensionName_example" // string | Name of the extension

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExtensionsAPI.GetExtensionConfigSchema(context.Background(), extensionPointName, extensionName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.GetExtensionConfigSchema``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExtensionConfigSchema`: ExtensionConfigSchema
	fmt.Fprintf(os.Stdout, "Response from `ExtensionsAPI.GetExtensionConfigSchema`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**extensionPointName** | **string** | Name of the extension point | 
**extensionName** | **string** | Name of the extension | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExtensionConfigSchemaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ExtensionConfigSchema**](ExtensionConfigSchema.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListExtensionPoints

> ListExtensionPointsResponse ListExtensionPoints(ctx).Execute()

List all extension points



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
	resp, r, err := apiClient.ExtensionsAPI.ListExtensionPoints(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.ListExtensionPoints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListExtensionPoints`: ListExtensionPointsResponse
	fmt.Fprintf(os.Stdout, "Response from `ExtensionsAPI.ListExtensionPoints`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListExtensionPointsRequest struct via the builder pattern


### Return type

[**ListExtensionPointsResponse**](ListExtensionPointsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListExtensions

> ListExtensionsResponse ListExtensions(ctx, extensionPointName).Execute()

List all extensions



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
	extensionPointName := "extensionPointName_example" // string | Name of the extension point

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExtensionsAPI.ListExtensions(context.Background(), extensionPointName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.ListExtensions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListExtensions`: ListExtensionsResponse
	fmt.Fprintf(os.Stdout, "Response from `ExtensionsAPI.ListExtensions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**extensionPointName** | **string** | Name of the extension point | 

### Other Parameters

Other parameters are passed through a pointer to a apiListExtensionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListExtensionsResponse**](ListExtensionsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TestExtension

> TestExtensionResponse TestExtension(ctx, extensionPointName, extensionName).TestExtensionRequest(testExtensionRequest).Execute()

Test extension



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
	extensionPointName := "extensionPointName_example" // string | Name of the extension point
	extensionName := "extensionName_example" // string | Name of the extension
	testExtensionRequest := *openapiclient.NewTestExtensionRequest() // TestExtensionRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExtensionsAPI.TestExtension(context.Background(), extensionPointName, extensionName).TestExtensionRequest(testExtensionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.TestExtension``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TestExtension`: TestExtensionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExtensionsAPI.TestExtension`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**extensionPointName** | **string** | Name of the extension point | 
**extensionName** | **string** | Name of the extension | 

### Other Parameters

Other parameters are passed through a pointer to a apiTestExtensionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **testExtensionRequest** | [**TestExtensionRequest**](TestExtensionRequest.md) |  | 

### Return type

[**TestExtensionResponse**](TestExtensionResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateExtensionConfig

> UpdateExtensionConfig(ctx, extensionPointName, extensionName).UpdateExtensionConfigRequest(updateExtensionConfigRequest).Execute()

Update extension configuration



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
	extensionPointName := "extensionPointName_example" // string | Name of the extension point
	extensionName := "extensionName_example" // string | Name of the extension
	updateExtensionConfigRequest := *openapiclient.NewUpdateExtensionConfigRequest(map[string]interface{}{"key": interface{}(123)}) // UpdateExtensionConfigRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExtensionsAPI.UpdateExtensionConfig(context.Background(), extensionPointName, extensionName).UpdateExtensionConfigRequest(updateExtensionConfigRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtensionsAPI.UpdateExtensionConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**extensionPointName** | **string** | Name of the extension point | 
**extensionName** | **string** | Name of the extension | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateExtensionConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateExtensionConfigRequest** | [**UpdateExtensionConfigRequest**](UpdateExtensionConfigRequest.md) |  | 

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

