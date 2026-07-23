# \ComponentsAPI

All URIs are relative to */api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateComponent**](ComponentsAPI.md#CreateComponent) | **Post** /components | Creates a new component for the project
[**ListComponents**](ComponentsAPI.md#ListComponents) | **Get** /components | List all components



## CreateComponent

> CreateComponent(ctx).CreateComponentRequest(createComponentRequest).Execute()

Creates a new component for the project



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
	createComponentRequest := *openapiclient.NewCreateComponentRequest("ProjectUuid_example", "Name_example") // CreateComponentRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ComponentsAPI.CreateComponent(context.Background()).CreateComponentRequest(createComponentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ComponentsAPI.CreateComponent``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateComponentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createComponentRequest** | [**CreateComponentRequest**](CreateComponentRequest.md) |  | 

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


## ListComponents

> ListComponentsResponse ListComponents(ctx).GroupContains(groupContains).NameContains(nameContains).VersionContains(versionContains).PurlPrefix(purlPrefix).Cpe(cpe).SwidTagIdContains(swidTagIdContains).HashType(hashType).Hash(hash).PackageArtifactPublishedSince(packageArtifactPublishedSince).PackageArtifactPublishedBefore(packageArtifactPublishedBefore).ProjectState(projectState).ProjectLatestVersion(projectLatestVersion).Expand(expand).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()

List all components



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
	groupContains := "groupContains_example" // string | Filter by group (substring match) (optional)
	nameContains := "nameContains_example" // string | Filter by name (substring match) (optional)
	versionContains := "versionContains_example" // string | Filter by version (substring match) (optional)
	purlPrefix := "purlPrefix_example" // string | Filter by PURL (prefix match).  Must be a valid PURL, with at least `pkg:<ecosystem>/<name>` populated. (optional)
	cpe := "cpe_example" // string | Filter by CPE (exact match).  Must be a valid CPE. (optional)
	swidTagIdContains := "swidTagIdContains_example" // string | Filter by SWID Tag ID (substring match) (optional)
	hashType := "hashType_example" // string | The hash type to filter by (optional)
	hash := "hash_example" // string | Filter by hash value (exact match).  Requires `hash_type` to be set. (optional)
	packageArtifactPublishedSince := int64(789) // int64 | Filter by package artifact publish date (inclusive lower bound).  Note that components without resolved package artifact metadata, or whose upstream repository did not report a publication date, are excluded whenever `package_artifact_published_since` or `package_artifact_published_before` is set. (optional)
	packageArtifactPublishedBefore := int64(789) // int64 | Filter by package artifact publish date (exclusive upper bound). (optional)
	projectState := openapiclient.project-state("ACTIVE") // ProjectState | Filter by the state of the project that the component belongs to.  Omit to include components from projects in any state. (optional)
	projectLatestVersion := true // bool | Filter by whether the project the component belongs to is flagged as the latest version.  When `true`, only components from latest-version projects are returned. When `false`, only components from non-latest projects are returned. Omit to include components regardless of the flag. (optional)
	expand := []string{"Inner_example"} // []string | Optional fields to include in each component response item. Unknown values are silently ignored. (optional)
	limit := int32(56) // int32 | Maximum number of items to retrieve from the collection (optional) (default to 100)
	pageToken := "pageToken_example" // string | Opaque token pointing to a specific position in a collection (optional)
	sortDirection := openapiclient.sort-direction("ASC") // SortDirection |  (optional)
	sortBy := "sortBy_example" // string | Field to sort by. Refer to the operation description for information about which fields are sortable. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ComponentsAPI.ListComponents(context.Background()).GroupContains(groupContains).NameContains(nameContains).VersionContains(versionContains).PurlPrefix(purlPrefix).Cpe(cpe).SwidTagIdContains(swidTagIdContains).HashType(hashType).Hash(hash).PackageArtifactPublishedSince(packageArtifactPublishedSince).PackageArtifactPublishedBefore(packageArtifactPublishedBefore).ProjectState(projectState).ProjectLatestVersion(projectLatestVersion).Expand(expand).Limit(limit).PageToken(pageToken).SortDirection(sortDirection).SortBy(sortBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ComponentsAPI.ListComponents``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListComponents`: ListComponentsResponse
	fmt.Fprintf(os.Stdout, "Response from `ComponentsAPI.ListComponents`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListComponentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupContains** | **string** | Filter by group (substring match) | 
 **nameContains** | **string** | Filter by name (substring match) | 
 **versionContains** | **string** | Filter by version (substring match) | 
 **purlPrefix** | **string** | Filter by PURL (prefix match).  Must be a valid PURL, with at least &#x60;pkg:&lt;ecosystem&gt;/&lt;name&gt;&#x60; populated. | 
 **cpe** | **string** | Filter by CPE (exact match).  Must be a valid CPE. | 
 **swidTagIdContains** | **string** | Filter by SWID Tag ID (substring match) | 
 **hashType** | **string** | The hash type to filter by | 
 **hash** | **string** | Filter by hash value (exact match).  Requires &#x60;hash_type&#x60; to be set. | 
 **packageArtifactPublishedSince** | **int64** | Filter by package artifact publish date (inclusive lower bound).  Note that components without resolved package artifact metadata, or whose upstream repository did not report a publication date, are excluded whenever &#x60;package_artifact_published_since&#x60; or &#x60;package_artifact_published_before&#x60; is set. | 
 **packageArtifactPublishedBefore** | **int64** | Filter by package artifact publish date (exclusive upper bound). | 
 **projectState** | [**ProjectState**](ProjectState.md) | Filter by the state of the project that the component belongs to.  Omit to include components from projects in any state. | 
 **projectLatestVersion** | **bool** | Filter by whether the project the component belongs to is flagged as the latest version.  When &#x60;true&#x60;, only components from latest-version projects are returned. When &#x60;false&#x60;, only components from non-latest projects are returned. Omit to include components regardless of the flag. | 
 **expand** | **[]string** | Optional fields to include in each component response item. Unknown values are silently ignored. | 
 **limit** | **int32** | Maximum number of items to retrieve from the collection | [default to 100]
 **pageToken** | **string** | Opaque token pointing to a specific position in a collection | 
 **sortDirection** | [**SortDirection**](SortDirection.md) |  | 
 **sortBy** | **string** | Field to sort by. Refer to the operation description for information about which fields are sortable. | 

### Return type

[**ListComponentsResponse**](ListComponentsResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

