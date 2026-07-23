# \WorkflowsAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetWorkflowInstance**](WorkflowsAPI.md#GetWorkflowInstance) | **Get** /internal/workflow-instances/{id} | Get a workflow instance
[**GetWorkflowRun**](WorkflowsAPI.md#GetWorkflowRun) | **Get** /internal/workflow-runs/{id} | Get a workflow run
[**ListWorkflowRunEvents**](WorkflowsAPI.md#ListWorkflowRunEvents) | **Get** /internal/workflow-runs/{id}/events | List all events of a workflow run
[**ListWorkflowRuns**](WorkflowsAPI.md#ListWorkflowRuns) | **Get** /internal/workflow-runs | List all workflow runs



## GetWorkflowInstance

> WorkflowRunMetadata GetWorkflowInstance(ctx, id).Execute()

Get a workflow instance



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
	id := "id_example" // string | ID of the workflow instance

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowsAPI.GetWorkflowInstance(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowsAPI.GetWorkflowInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWorkflowInstance`: WorkflowRunMetadata
	fmt.Fprintf(os.Stdout, "Response from `WorkflowsAPI.GetWorkflowInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the workflow instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWorkflowInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**WorkflowRunMetadata**](WorkflowRunMetadata.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetWorkflowRun

> WorkflowRunMetadata GetWorkflowRun(ctx, id).Execute()

Get a workflow run



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the workflow run

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowsAPI.GetWorkflowRun(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowsAPI.GetWorkflowRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWorkflowRun`: WorkflowRunMetadata
	fmt.Fprintf(os.Stdout, "Response from `WorkflowsAPI.GetWorkflowRun`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the workflow run | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWorkflowRunRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**WorkflowRunMetadata**](WorkflowRunMetadata.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWorkflowRunEvents

> ListWorkflowRunEventsResponse ListWorkflowRunEvents(ctx, id).FromSequenceNumber(fromSequenceNumber).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).Execute()

List all events of a workflow run



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the workflow run
	fromSequenceNumber := int32(56) // int32 | Sequence number of the last seen event. May be used to continuously poll for new events. Can not be used together with `page_token`. (optional)
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)
	sortDirection := openapiclient.sort-direction("ASC") // SortDirection |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowsAPI.ListWorkflowRunEvents(context.Background(), id).FromSequenceNumber(fromSequenceNumber).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowsAPI.ListWorkflowRunEvents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListWorkflowRunEvents`: ListWorkflowRunEventsResponse
	fmt.Fprintf(os.Stdout, "Response from `WorkflowsAPI.ListWorkflowRunEvents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the workflow run | 

### Other Parameters

Other parameters are passed through a pointer to a apiListWorkflowRunEventsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **fromSequenceNumber** | **int32** | Sequence number of the last seen event. May be used to continuously poll for new events. Can not be used together with &#x60;page_token&#x60;. | 
 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 
 **sortDirection** | [**SortDirection**](SortDirection.md) |  | 

### Return type

[**ListWorkflowRunEventsResponse**](ListWorkflowRunEventsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWorkflowRuns

> ListWorkflowRunsResponse ListWorkflowRuns(ctx).WorkflowName(workflowName).WorkflowVersion(workflowVersion).WorkflowInstanceId(workflowInstanceId).Status(status).Label(label).CreatedSince(createdSince).CreatedBefore(createdBefore).CompletedSince(completedSince).CompletedBefore(completedBefore).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()

List all workflow runs



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
	workflowName := "workflowName_example" // string | Name of the workflow to filter by. Must be an exact match. (optional)
	workflowVersion := int32(56) // int32 | Version of the workflow to filter by. Must be an exact match. (optional)
	workflowInstanceId := "workflowInstanceId_example" // string | Workflow instance ID to filter by. Must be an exact match. (optional)
	status := openapiclient.workflow-run-status("CREATED") // WorkflowRunStatus | Status to filter by (optional)
	label := []string{"Inner_example"} // []string | Filter by label in `key=value` form. Repeat to require multiple labels. A run matches only if it carries every supplied label. On duplicate keys the last occurrence wins. (optional)
	createdSince := int64(789) // int64 | Filter runs created on or after this timestamp. (optional)
	createdBefore := int64(789) // int64 | Filter runs created before this timestamp. (optional)
	completedSince := int64(789) // int64 | Filter runs completed on or after this timestamp. (optional)
	completedBefore := int64(789) // int64 | Filter runs completed before this timestamp. (optional)
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)
	sortDirection := openapiclient.sort-direction("ASC") // SortDirection |  (optional)
	sortBy := "sortBy_example" // string | Field to sort by. Refer to the operation description for information about which fields are sortable. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowsAPI.ListWorkflowRuns(context.Background()).WorkflowName(workflowName).WorkflowVersion(workflowVersion).WorkflowInstanceId(workflowInstanceId).Status(status).Label(label).CreatedSince(createdSince).CreatedBefore(createdBefore).CompletedSince(completedSince).CompletedBefore(completedBefore).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowsAPI.ListWorkflowRuns``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListWorkflowRuns`: ListWorkflowRunsResponse
	fmt.Fprintf(os.Stdout, "Response from `WorkflowsAPI.ListWorkflowRuns`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListWorkflowRunsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowName** | **string** | Name of the workflow to filter by. Must be an exact match. | 
 **workflowVersion** | **int32** | Version of the workflow to filter by. Must be an exact match. | 
 **workflowInstanceId** | **string** | Workflow instance ID to filter by. Must be an exact match. | 
 **status** | [**WorkflowRunStatus**](WorkflowRunStatus.md) | Status to filter by | 
 **label** | **[]string** | Filter by label in &#x60;key&#x3D;value&#x60; form. Repeat to require multiple labels. A run matches only if it carries every supplied label. On duplicate keys the last occurrence wins. | 
 **createdSince** | **int64** | Filter runs created on or after this timestamp. | 
 **createdBefore** | **int64** | Filter runs created before this timestamp. | 
 **completedSince** | **int64** | Filter runs completed on or after this timestamp. | 
 **completedBefore** | **int64** | Filter runs completed before this timestamp. | 
 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 
 **sortDirection** | [**SortDirection**](SortDirection.md) |  | 
 **sortBy** | **string** | Field to sort by. Refer to the operation description for information about which fields are sortable. | 

### Return type

[**ListWorkflowRunsResponse**](ListWorkflowRunsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

