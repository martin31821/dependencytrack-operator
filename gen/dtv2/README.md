# Go API client for dtv2

# REST API of OWASP Dependency-Track

## Pagination

This API implements token-based pagination.
Collection responses have the following structure:

```json
{
  \"items\": [...],
  \"next_page_token\": \"abcdefg\",
  \"total\": {
    \"count\": 100,
    \"type\": \"EXACT\"
  }
}
```

`next_page_token` is present when more items exist, and absent otherwise.
To fetch the next page, pass it as the `page_token` query parameter.

To navigate backwards, clients should keep track of previous page tokens
as they paginate through collections. The API does *not* provide
backward navigation!

Collections that support sorting will only consider the `sort_by`
and `sort_direction` query parameters for the request of the first
page. For subsequent pages, sorting preferences are bound to the
page token.

Page tokens are opaque strings. Clients should not try to interpret
or generate them. Their format may change without notice.

The `total` object discloses how many items exist in the collection
*across all pages*. Because counting is expensive, some collections
that hold *a lot* of items may return partial counts (type `AT_LEAST`)
instead of exact counts (type `EXACT`). Which type to expect is usually
documented in the operation's description.

## Sorting

Items in a collection can be sorted using the `sort_by` and `sort_direction`
query parameters. Which fields are sortable is documented in the respective
operation's description.

Note that if no sortable fields are documented for an operation,
sorting is not supported *at all*.

## Field expansion

Some collection endpoints support an `expand` query parameter.
Passing an expand value includes optional fields in each response item
that are omitted by default, typically because they are expensive to compute
and only needed in specific contexts.

Valid `expand` values for an endpoint are listed in its operation description.
Unknown values are silently ignored.

## Errors

All error responses use the `application/problem+json` media type
as defined in [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html).

Example:

```json
{
  \"type\": \"about:blank\",
  \"status\": 404,
  \"title\": \"Resource Not Found\",
  \"detail\": \"No project with UUID 0976ef46-e4a0-4be4-8b0e-00e9d3625c1f exists\"
}
```

## Authentication

Two credential types are accepted:

- **API key** via the `X-Api-Key` request header. API keys are
  long-lived and intended for machine-to-machine integrations.
- **Bearer token** via the `Authorization: Bearer <token>` request
  header. Bearer tokens are short-lived, user-bound, and opaque
  server-issued session tokens.

Bearer tokens are obtained from one of the following endpoints:

- `POST /api/v1/user/login`
- `POST /api/v1/user/oidc/login`

API v2 does not yet provide its own login endpoints; clients use
the v1 endpoints above to acquire tokens and then call v2 with them.

Tokens are valid for 8 hours by default and **cannot be refreshed**.
Clients must re-authenticate once a token expires.

Requests with missing or invalid credentials are rejected with
`401 Unauthorized`.

## Authorization

Access is gated by named permissions. Operations document
the permission(s) they require; operations without a documented
permission requirement only require authentication.

When the *Portfolio Access Control* feature is enabled
(disabled by default), project-scoped operations additionally enforce
per-project access via team membership. The `PORTFOLIO_ACCESS_CONTROL_BYPASS`
permission grants access to all projects regardless of team mappings.
When the feature is disabled, all authenticated callers holding the required
permission can access all projects.

Authenticated requests that lack the required permission, or that
target a project the caller cannot access, are rejected with
`403 Forbidden`.

## HTTP Methods

| Method   | Semantics                  |
|----------|----------------------------|
| `GET`    | Retrieve a resource        |
| `POST`   | Create a new resource      |
| `PUT`    | Update a resource          |
| `PATCH`  | Partially update a resource|
| `DELETE` | Delete a resource          |

## Response Conventions

Create and update operations (`POST`, `PUT`, `PATCH`) do not return
the full resource in the response. They return either no body,
or only server-generated identifiers (e.g. a UUID).
`POST` responses may include a `Location` header linking to the
created resource.

Delete operations return `204 No Content` with no body.

## Deprecations

Operations may be removed or replaced over time. When a response
carries the `X-API-Deprecated: true` header, the operation that
produced it is deprecated and may be removed in a future release.
Clients should check for this header on every response and surface
it (e.g. via a log warning) so that operators are aware of upcoming
breakages. The respective operation's description points out which
alternative operation(s) to use.

## Internal operations

Operations under the `/internal` path prefix expose system internals
and are reserved for first-party use. They are **not** part of the
stable v2 API contract and may change or be removed without notice.
Third-party clients should not depend on them.

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 2.0.0
- Package version: 1.0.0
- Generator version: 7.22.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [https://github.com/DependencyTrack/dependency-track](https://github.com/DependencyTrack/dependency-track)

## Installation

Import the package in a go file in your project and run `go mod tidy`:

```go
import dtv2 "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `dtv2.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), dtv2.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `dtv2.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), dtv2.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `dtv2.ContextOperationServerIndices` and `dtv2.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), dtv2.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), dtv2.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to */api/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ComponentsAPI* | [**CreateComponent**](docs/ComponentsAPI.md#createcomponent) | **Post** /components | Creates a new component for the project
*ComponentsAPI* | [**ListComponents**](docs/ComponentsAPI.md#listcomponents) | **Get** /components | List all components
*ExtensionsAPI* | [**GetExtensionConfig**](docs/ExtensionsAPI.md#getextensionconfig) | **Get** /extension-points/{extension_point_name}/extensions/{extension_name}/config | Get extension configuration
*ExtensionsAPI* | [**GetExtensionConfigSchema**](docs/ExtensionsAPI.md#getextensionconfigschema) | **Get** /extension-points/{extension_point_name}/extensions/{extension_name}/config-schema | Get extension configuration schema
*ExtensionsAPI* | [**ListExtensionPoints**](docs/ExtensionsAPI.md#listextensionpoints) | **Get** /extension-points | List all extension points
*ExtensionsAPI* | [**ListExtensions**](docs/ExtensionsAPI.md#listextensions) | **Get** /extension-points/{extension_point_name}/extensions | List all extensions
*ExtensionsAPI* | [**TestExtension**](docs/ExtensionsAPI.md#testextension) | **Post** /extension-points/{extension_point_name}/extensions/{extension_name}/test | Test extension
*ExtensionsAPI* | [**UpdateExtensionConfig**](docs/ExtensionsAPI.md#updateextensionconfig) | **Put** /extension-points/{extension_point_name}/extensions/{extension_name}/config | Update extension configuration
*ProjectsAPI* | [**CloneProject**](docs/ProjectsAPI.md#cloneproject) | **Post** /projects/{uuid}/clone | Clones a given project.
*ProjectsAPI* | [**ListProjectComponents**](docs/ProjectsAPI.md#listprojectcomponents) | **Get** /projects/{uuid}/components | Retrieves a list of all components for a given project.
*SecretsAPI* | [**CreateSecret**](docs/SecretsAPI.md#createsecret) | **Post** /secrets | Create a secret
*SecretsAPI* | [**DeleteSecret**](docs/SecretsAPI.md#deletesecret) | **Delete** /secrets/{name} | Delete a secret
*SecretsAPI* | [**GetSecretMetadata**](docs/SecretsAPI.md#getsecretmetadata) | **Get** /secrets/{name} | Get secret metadata
*SecretsAPI* | [**ListSecretMetadata**](docs/SecretsAPI.md#listsecretmetadata) | **Get** /secrets | List secret metadata
*SecretsAPI* | [**UpdateSecret**](docs/SecretsAPI.md#updatesecret) | **Patch** /secrets/{name} | Update a secret
*SystemCapabilitiesAPI* | [**GetSystemCapabilities**](docs/SystemCapabilitiesAPI.md#getsystemcapabilities) | **Get** /internal/system-capabilities | Get system capabilities
*TaskQueuesAPI* | [**ListTaskQueues**](docs/TaskQueuesAPI.md#listtaskqueues) | **Get** /internal/task-queues/{type} | List task queues
*TaskQueuesAPI* | [**UpdateTaskQueue**](docs/TaskQueuesAPI.md#updatetaskqueue) | **Patch** /internal/task-queues/{type}/{name} | Update a task queue
*VulnDataSourcesAPI* | [**GetLatestVulnDataSourceMirrorRun**](docs/VulnDataSourcesAPI.md#getlatestvulndatasourcemirrorrun) | **Get** /vuln-data-sources/{name}/mirror-runs/latest | Get the latest vulnerability data source mirror run
*VulnDataSourcesAPI* | [**TriggerVulnDataSourceMirrorRun**](docs/VulnDataSourcesAPI.md#triggervulndatasourcemirrorrun) | **Post** /vuln-data-sources/{name}/mirror-runs | Trigger a vulnerability data source mirror run
*VulnPoliciesAPI* | [**CreateVulnPolicy**](docs/VulnPoliciesAPI.md#createvulnpolicy) | **Post** /vuln-policies | Create a vulnerability policy
*VulnPoliciesAPI* | [**DeleteVulnPolicy**](docs/VulnPoliciesAPI.md#deletevulnpolicy) | **Delete** /vuln-policies/{uuid} | Delete a vulnerability policy
*VulnPoliciesAPI* | [**DeleteVulnPolicyBundle**](docs/VulnPoliciesAPI.md#deletevulnpolicybundle) | **Delete** /vuln-policy-bundles/{uuid} | Delete a vulnerability policy bundle
*VulnPoliciesAPI* | [**GetLatestVulnPolicyBundleSyncRun**](docs/VulnPoliciesAPI.md#getlatestvulnpolicybundlesyncrun) | **Get** /vuln-policy-bundles/{uuid}/sync-runs/latest | Get the latest vulnerability policy bundle sync run
*VulnPoliciesAPI* | [**GetVulnPolicy**](docs/VulnPoliciesAPI.md#getvulnpolicy) | **Get** /vuln-policies/{uuid} | Get a vulnerability policy
*VulnPoliciesAPI* | [**ListVulnPolicies**](docs/VulnPoliciesAPI.md#listvulnpolicies) | **Get** /vuln-policies | List vulnerability policies
*VulnPoliciesAPI* | [**ListVulnPolicyBundles**](docs/VulnPoliciesAPI.md#listvulnpolicybundles) | **Get** /vuln-policy-bundles | List vulnerability policy bundles
*VulnPoliciesAPI* | [**TriggerVulnPolicyBundleSyncRun**](docs/VulnPoliciesAPI.md#triggervulnpolicybundlesyncrun) | **Post** /vuln-policy-bundles/{uuid}/sync-runs | Trigger a vulnerability policy bundle sync run
*VulnPoliciesAPI* | [**UpdateVulnPolicy**](docs/VulnPoliciesAPI.md#updatevulnpolicy) | **Put** /vuln-policies/{uuid} | Replace a vulnerability policy
*WorkflowsAPI* | [**GetWorkflowInstance**](docs/WorkflowsAPI.md#getworkflowinstance) | **Get** /internal/workflow-instances/{id} | Get a workflow instance
*WorkflowsAPI* | [**GetWorkflowRun**](docs/WorkflowsAPI.md#getworkflowrun) | **Get** /internal/workflow-runs/{id} | Get a workflow run
*WorkflowsAPI* | [**ListWorkflowRunEvents**](docs/WorkflowsAPI.md#listworkflowrunevents) | **Get** /internal/workflow-runs/{id}/events | List all events of a workflow run
*WorkflowsAPI* | [**ListWorkflowRuns**](docs/WorkflowsAPI.md#listworkflowruns) | **Get** /internal/workflow-runs | List all workflow runs


## Documentation For Models

 - [Classifier](docs/Classifier.md)
 - [CloneProjectInclude](docs/CloneProjectInclude.md)
 - [CloneProjectRequest](docs/CloneProjectRequest.md)
 - [CloneProjectResponse](docs/CloneProjectResponse.md)
 - [ComponentProject](docs/ComponentProject.md)
 - [ConstraintViolationError](docs/ConstraintViolationError.md)
 - [CreateComponent400Response](docs/CreateComponent400Response.md)
 - [CreateComponentRequest](docs/CreateComponentRequest.md)
 - [CreateSecretRequest](docs/CreateSecretRequest.md)
 - [CreateVulnPolicy201Response](docs/CreateVulnPolicy201Response.md)
 - [CreateVulnPolicy400Response](docs/CreateVulnPolicy400Response.md)
 - [CreateVulnPolicyRequest](docs/CreateVulnPolicyRequest.md)
 - [DependencyMetrics](docs/DependencyMetrics.md)
 - [ExtensionConfigSchema](docs/ExtensionConfigSchema.md)
 - [ExtensionTestCheck](docs/ExtensionTestCheck.md)
 - [ExtensionTestCheckStatus](docs/ExtensionTestCheckStatus.md)
 - [GetExtensionConfigResponse](docs/GetExtensionConfigResponse.md)
 - [GetVulnPolicyResponse](docs/GetVulnPolicyResponse.md)
 - [Hashes](docs/Hashes.md)
 - [InvalidRequestProblemDetails](docs/InvalidRequestProblemDetails.md)
 - [InvalidSortFieldProblemDetails](docs/InvalidSortFieldProblemDetails.md)
 - [InvalidVulnPolicyConditionProblemDetails](docs/InvalidVulnPolicyConditionProblemDetails.md)
 - [JsonSchemaValidationError](docs/JsonSchemaValidationError.md)
 - [JsonSchemaValidationProblemDetails](docs/JsonSchemaValidationProblemDetails.md)
 - [License](docs/License.md)
 - [ListComponents400Response](docs/ListComponents400Response.md)
 - [ListComponentsResponse](docs/ListComponentsResponse.md)
 - [ListComponentsResponseItem](docs/ListComponentsResponseItem.md)
 - [ListExtensionPointsResponse](docs/ListExtensionPointsResponse.md)
 - [ListExtensionPointsResponseItem](docs/ListExtensionPointsResponseItem.md)
 - [ListExtensionsResponse](docs/ListExtensionsResponse.md)
 - [ListExtensionsResponseItem](docs/ListExtensionsResponseItem.md)
 - [ListProjectComponentsResponse](docs/ListProjectComponentsResponse.md)
 - [ListProjectComponentsResponseItem](docs/ListProjectComponentsResponseItem.md)
 - [ListSecretsResponse](docs/ListSecretsResponse.md)
 - [ListTaskQueuesResponse](docs/ListTaskQueuesResponse.md)
 - [ListVulnPoliciesResponse](docs/ListVulnPoliciesResponse.md)
 - [ListVulnPoliciesResponseItem](docs/ListVulnPoliciesResponseItem.md)
 - [ListVulnPolicyBundlesResponse](docs/ListVulnPolicyBundlesResponse.md)
 - [ListVulnPolicyBundlesResponseItem](docs/ListVulnPolicyBundlesResponseItem.md)
 - [ListWorkflowRunEventsResponse](docs/ListWorkflowRunEventsResponse.md)
 - [ListWorkflowRunEventsResponseItem](docs/ListWorkflowRunEventsResponseItem.md)
 - [ListWorkflowRuns400Response](docs/ListWorkflowRuns400Response.md)
 - [ListWorkflowRunsResponse](docs/ListWorkflowRunsResponse.md)
 - [OrganizationalContact](docs/OrganizationalContact.md)
 - [OrganizationalEntity](docs/OrganizationalEntity.md)
 - [PackageArtifactMetadata](docs/PackageArtifactMetadata.md)
 - [PackageMetadata](docs/PackageMetadata.md)
 - [PaginatedResponse](docs/PaginatedResponse.md)
 - [ProblemDetails](docs/ProblemDetails.md)
 - [ProjectState](docs/ProjectState.md)
 - [Scope](docs/Scope.md)
 - [SecretMetadata](docs/SecretMetadata.md)
 - [SortDirection](docs/SortDirection.md)
 - [SystemCapabilitiesResponse](docs/SystemCapabilitiesResponse.md)
 - [TaskQueue](docs/TaskQueue.md)
 - [TaskQueueStatus](docs/TaskQueueStatus.md)
 - [TaskQueueType](docs/TaskQueueType.md)
 - [TestExtensionRequest](docs/TestExtensionRequest.md)
 - [TestExtensionResponse](docs/TestExtensionResponse.md)
 - [TotalCount](docs/TotalCount.md)
 - [TotalCountType](docs/TotalCountType.md)
 - [UpdateExtensionConfig400Response](docs/UpdateExtensionConfig400Response.md)
 - [UpdateExtensionConfigRequest](docs/UpdateExtensionConfigRequest.md)
 - [UpdateSecretRequest](docs/UpdateSecretRequest.md)
 - [UpdateTaskQueueRequest](docs/UpdateTaskQueueRequest.md)
 - [UpdateVulnPolicyRequest](docs/UpdateVulnPolicyRequest.md)
 - [VulnDataSourceMirrorStatus](docs/VulnDataSourceMirrorStatus.md)
 - [VulnPolicyAnalysis](docs/VulnPolicyAnalysis.md)
 - [VulnPolicyBundleSyncStatus](docs/VulnPolicyBundleSyncStatus.md)
 - [VulnPolicyConditionError](docs/VulnPolicyConditionError.md)
 - [VulnPolicyOperationMode](docs/VulnPolicyOperationMode.md)
 - [VulnPolicyRating](docs/VulnPolicyRating.md)
 - [VulnPolicySource](docs/VulnPolicySource.md)
 - [WorkflowRunMetadata](docs/WorkflowRunMetadata.md)
 - [WorkflowRunStatus](docs/WorkflowRunStatus.md)


## Documentation For Authorization


Authentication schemes defined for the API:
### apiKeyAuth

- **Type**: API key
- **API key parameter name**: X-Api-Key
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: apiKeyAuth and passed in as the auth context for each request.

Example

```go
auth := context.WithValue(
		context.Background(),
		dtv2.ContextAPIKeys,
		map[string]dtv2.APIKey{
			"apiKeyAuth": {Key: "API_KEY_STRING"},
		},
	)
r, err := client.Service.Operation(auth, args)
```

### bearerAuth

- **Type**: HTTP Bearer token authentication

Example

```go
auth := context.WithValue(context.Background(), dtv2.ContextAccessToken, "BEARER_TOKEN_STRING")
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

dependencytrack@owasp.org

