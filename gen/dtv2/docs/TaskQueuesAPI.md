# \TaskQueuesAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListTaskQueues**](TaskQueuesAPI.md#ListTaskQueues) | **Get** /internal/task-queues/{type} | List task queues
[**UpdateTaskQueue**](TaskQueuesAPI.md#UpdateTaskQueue) | **Patch** /internal/task-queues/{type}/{name} | Update a task queue



## ListTaskQueues

> ListTaskQueuesResponse ListTaskQueues(ctx, type_).Limit(limit).PageToken(pageToken).Execute()

List task queues



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
	type_ := openapiclient.task-queue-type("ACTIVITY") // TaskQueueType | Type of task queues to list
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaskQueuesAPI.ListTaskQueues(context.Background(), type_).Limit(limit).PageToken(pageToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaskQueuesAPI.ListTaskQueues``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTaskQueues`: ListTaskQueuesResponse
	fmt.Fprintf(os.Stdout, "Response from `TaskQueuesAPI.ListTaskQueues`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**type_** | [**TaskQueueType**](.md) | Type of task queues to list | 

### Other Parameters

Other parameters are passed through a pointer to a apiListTaskQueuesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 

### Return type

[**ListTaskQueuesResponse**](ListTaskQueuesResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTaskQueue

> UpdateTaskQueue(ctx, type_, name).UpdateTaskQueueRequest(updateTaskQueueRequest).Execute()

Update a task queue



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
	type_ := openapiclient.task-queue-type("ACTIVITY") // TaskQueueType | Type of the task queue
	name := "name_example" // string | Name of the task queue
	updateTaskQueueRequest := *openapiclient.NewUpdateTaskQueueRequest() // UpdateTaskQueueRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TaskQueuesAPI.UpdateTaskQueue(context.Background(), type_, name).UpdateTaskQueueRequest(updateTaskQueueRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaskQueuesAPI.UpdateTaskQueue``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**type_** | [**TaskQueueType**](.md) | Type of the task queue | 
**name** | **string** | Name of the task queue | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTaskQueueRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateTaskQueueRequest** | [**UpdateTaskQueueRequest**](UpdateTaskQueueRequest.md) |  | 

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

