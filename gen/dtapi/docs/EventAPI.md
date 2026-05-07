# \EventAPI

All URIs are relative to *https://hyades-api.iris-flair-alpha.vlair-staging.defra01.iris-sensing.net/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**IsTokenBeingProcessed**](EventAPI.md#IsTokenBeingProcessed) | **Get** /v1/event/token/{uuid} | Determines if there are any tasks associated with the token that are being processed, or in the queue to be processed.



## IsTokenBeingProcessed

> IsTokenBeingProcessedResponse IsTokenBeingProcessed(ctx, uuid).Execute()

Determines if there are any tasks associated with the token that are being processed, or in the queue to be processed.



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/martin31821/dependencytrack-operator/gen/dtapi"
)

func main() {
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the token to query

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EventAPI.IsTokenBeingProcessed(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EventAPI.IsTokenBeingProcessed``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `IsTokenBeingProcessed`: IsTokenBeingProcessedResponse
	fmt.Fprintf(os.Stdout, "Response from `EventAPI.IsTokenBeingProcessed`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the token to query | 

### Other Parameters

Other parameters are passed through a pointer to a apiIsTokenBeingProcessedRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**IsTokenBeingProcessedResponse**](IsTokenBeingProcessedResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

