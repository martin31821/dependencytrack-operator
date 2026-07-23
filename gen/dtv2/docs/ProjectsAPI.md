# \ProjectsAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CloneProject**](ProjectsAPI.md#CloneProject) | **Post** /projects/{uuid}/clone | Clones a given project.
[**ListProjectComponents**](ProjectsAPI.md#ListProjectComponents) | **Get** /projects/{uuid}/components | Retrieves a list of all components for a given project.



## CloneProject

> CloneProjectResponse CloneProject(ctx, uuid).CloneProjectRequest(cloneProjectRequest).Execute()

Clones a given project.



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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the project to clone
	cloneProjectRequest := *openapiclient.NewCloneProjectRequest("Version_example") // CloneProjectRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProjectsAPI.CloneProject(context.Background(), uuid).CloneProjectRequest(cloneProjectRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.CloneProject``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CloneProject`: CloneProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.CloneProject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the project to clone | 

### Other Parameters

Other parameters are passed through a pointer to a apiCloneProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **cloneProjectRequest** | [**CloneProjectRequest**](CloneProjectRequest.md) |  | 

### Return type

[**CloneProjectResponse**](CloneProjectResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListProjectComponents

> ListProjectComponentsResponse ListProjectComponents(ctx, uuid).OnlyOutdated(onlyOutdated).OnlyDirect(onlyDirect).Q(q).Expand(expand).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()

Retrieves a list of all components for a given project.



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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the project to retrieve components for
	onlyOutdated := true // bool | Optionally exclude recent components so only outdated components are returned (optional)
	onlyDirect := true // bool | Optionally exclude transitive dependencies so only direct dependencies are returned (optional)
	q := "q_example" // string | Optional free-text search term. Matches components whose `group` or `name` contains the given value (case-insensitive). (optional)
	expand := []string{"Inner_example"} // []string | Optional fields to include in each component response item. Unknown values are silently ignored. (optional)
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)
	sortDirection := openapiclient.sort-direction("ASC") // SortDirection |  (optional)
	sortBy := "sortBy_example" // string | Field to sort by. Refer to the operation description for information about which fields are sortable. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProjectsAPI.ListProjectComponents(context.Background(), uuid).OnlyOutdated(onlyOutdated).OnlyDirect(onlyDirect).Q(q).Expand(expand).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ListProjectComponents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListProjectComponents`: ListProjectComponentsResponse
	fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.ListProjectComponents`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the project to retrieve components for | 

### Other Parameters

Other parameters are passed through a pointer to a apiListProjectComponentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **onlyOutdated** | **bool** | Optionally exclude recent components so only outdated components are returned | 
 **onlyDirect** | **bool** | Optionally exclude transitive dependencies so only direct dependencies are returned | 
 **q** | **string** | Optional free-text search term. Matches components whose &#x60;group&#x60; or &#x60;name&#x60; contains the given value (case-insensitive). | 
 **expand** | **[]string** | Optional fields to include in each component response item. Unknown values are silently ignored. | 
 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 
 **sortDirection** | [**SortDirection**](SortDirection.md) |  | 
 **sortBy** | **string** | Field to sort by. Refer to the operation description for information about which fields are sortable. | 

### Return type

[**ListProjectComponentsResponse**](ListProjectComponentsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

