# \PermissionAPI

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddPermissionToTeam**](PermissionAPI.md#AddPermissionToTeam) | **Post** /v1/permission/{permission}/team/{uuid} | 
[**AddPermissionToUser**](PermissionAPI.md#AddPermissionToUser) | **Post** /v1/permission/{permission}/user/{username} | Adds the permission to the specified username.
[**GetAllPermissions**](PermissionAPI.md#GetAllPermissions) | **Get** /v1/permission | Returns a list of all permissions
[**RemovePermissionFromTeam**](PermissionAPI.md#RemovePermissionFromTeam) | **Delete** /v1/permission/{permission}/team/{uuid} | 
[**RemovePermissionFromUser**](PermissionAPI.md#RemovePermissionFromUser) | **Delete** /v1/permission/{permission}/user/{username} | Removes the permission from the user.
[**SetTeamPermissions**](PermissionAPI.md#SetTeamPermissions) | **Put** /v1/permission/team | Replaces a team&#39;s permissions with the specified list
[**SetUserPermissions**](PermissionAPI.md#SetUserPermissions) | **Put** /v1/permission/user | Replaces a users&#39;s permissions with the specified list



## AddPermissionToTeam

> Team AddPermissionToTeam(ctx, uuid, permission).Execute()





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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | A valid team uuid
	permission := "permission_example" // string | A valid permission

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.AddPermissionToTeam(context.Background(), uuid, permission).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.AddPermissionToTeam``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPermissionToTeam`: Team
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.AddPermissionToTeam`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | A valid team uuid | 
**permission** | **string** | A valid permission | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPermissionToTeamRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Team**](Team.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddPermissionToUser

> User AddPermissionToUser(ctx, username, permission).Execute()

Adds the permission to the specified username.



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
	username := "username_example" // string | A valid username
	permission := "permission_example" // string | A valid permission

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.AddPermissionToUser(context.Background(), username, permission).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.AddPermissionToUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddPermissionToUser`: User
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.AddPermissionToUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**username** | **string** | A valid username | 
**permission** | **string** | A valid permission | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddPermissionToUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**User**](User.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllPermissions

> string GetAllPermissions(ctx).Execute()

Returns a list of all permissions



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
	resp, r, err := apiClient.PermissionAPI.GetAllPermissions(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.GetAllPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllPermissions`: string
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.GetAllPermissions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllPermissionsRequest struct via the builder pattern


### Return type

**string**

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemovePermissionFromTeam

> Team RemovePermissionFromTeam(ctx, uuid, permission).Execute()





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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | A valid team uuid
	permission := "permission_example" // string | A valid permission

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.RemovePermissionFromTeam(context.Background(), uuid, permission).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.RemovePermissionFromTeam``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePermissionFromTeam`: Team
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.RemovePermissionFromTeam`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | A valid team uuid | 
**permission** | **string** | A valid permission | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePermissionFromTeamRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Team**](Team.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemovePermissionFromUser

> User RemovePermissionFromUser(ctx, username, permission).Execute()

Removes the permission from the user.



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
	username := "username_example" // string | A valid username
	permission := "permission_example" // string | A valid permission

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.RemovePermissionFromUser(context.Background(), username, permission).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.RemovePermissionFromUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemovePermissionFromUser`: User
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.RemovePermissionFromUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**username** | **string** | A valid username | 
**permission** | **string** | A valid permission | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemovePermissionFromUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**User**](User.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetTeamPermissions

> Team SetTeamPermissions(ctx).TeamPermissionsSetRequest(teamPermissionsSetRequest).Execute()

Replaces a team's permissions with the specified list



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
	teamPermissionsSetRequest := *openapiclient.NewTeamPermissionsSetRequest([]string{"Permissions_example"}, "Team_example") // TeamPermissionsSetRequest | Team UUID and requested permissions (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.SetTeamPermissions(context.Background()).TeamPermissionsSetRequest(teamPermissionsSetRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.SetTeamPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetTeamPermissions`: Team
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.SetTeamPermissions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSetTeamPermissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **teamPermissionsSetRequest** | [**TeamPermissionsSetRequest**](TeamPermissionsSetRequest.md) | Team UUID and requested permissions | 

### Return type

[**Team**](Team.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetUserPermissions

> User SetUserPermissions(ctx).UserPermissionsSetRequest(userPermissionsSetRequest).Execute()

Replaces a users's permissions with the specified list



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
	userPermissionsSetRequest := *openapiclient.NewUserPermissionsSetRequest([]string{"Permissions_example"}, "Username_example") // UserPermissionsSetRequest | A username and valid list permission (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PermissionAPI.SetUserPermissions(context.Background()).UserPermissionsSetRequest(userPermissionsSetRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PermissionAPI.SetUserPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetUserPermissions`: User
	fmt.Fprintf(os.Stdout, "Response from `PermissionAPI.SetUserPermissions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSetUserPermissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userPermissionsSetRequest** | [**UserPermissionsSetRequest**](UserPermissionsSetRequest.md) | A username and valid list permission | 

### Return type

[**User**](User.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

