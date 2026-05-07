# \WorkflowAPI

All URIs are relative to *https://hyades-api.iris-flair-alpha.vlair-staging.defra01.iris-sensing.net/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetWorkflowStates**](WorkflowAPI.md#GetWorkflowStates) | **Get** /v1/workflow/token/{uuid}/status | Retrieves workflow states associated with the token received from bom upload .



## GetWorkflowStates

> WorkflowState GetWorkflowStates(ctx, uuid).Execute()

Retrieves workflow states associated with the token received from bom upload .



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
	uuid := "uuid_example" // string | The UUID of the token to query

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.GetWorkflowStates(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.GetWorkflowStates``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWorkflowStates`: WorkflowState
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.GetWorkflowStates`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the token to query | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWorkflowStatesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**WorkflowState**](WorkflowState.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

