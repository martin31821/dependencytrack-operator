/*
OWASP Dependency-Track

# REST API of OWASP Dependency-Track  ## Pagination  This API implements token-based pagination. Collection responses have the following structure:  ```json {   \"items\": [...],   \"next_page_token\": \"abcdefg\",   \"total\": {     \"count\": 100,     \"type\": \"EXACT\"   } } ```  `next_page_token` is present when more items exist, and absent otherwise. To fetch the next page, pass it as the `page_token` query parameter.  To navigate backwards, clients should keep track of previous page tokens as they paginate through collections. The API does *not* provide backward navigation!  Collections that support sorting will only consider the `sort_by` and `sort_direction` query parameters for the request of the first page. For subsequent pages, sorting preferences are bound to the page token.  Page tokens are opaque strings. Clients should not try to interpret or generate them. Their format may change without notice.  The `total` object discloses how many items exist in the collection *across all pages*. Because counting is expensive, some collections that hold *a lot* of items may return partial counts (type `AT_LEAST`) instead of exact counts (type `EXACT`). Which type to expect is usually documented in the operation's description.  ## Sorting  Items in a collection can be sorted using the `sort_by` and `sort_direction` query parameters. Which fields are sortable is documented in the respective operation's description.  Note that if no sortable fields are documented for an operation, sorting is not supported *at all*.  ## Field expansion  Some collection endpoints support an `expand` query parameter. Passing an expand value includes optional fields in each response item that are omitted by default, typically because they are expensive to compute and only needed in specific contexts.  Valid `expand` values for an endpoint are listed in its operation description. Unknown values are silently ignored.  ## Errors  All error responses use the `application/problem+json` media type as defined in [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html).  Example:  ```json {   \"type\": \"about:blank\",   \"status\": 404,   \"title\": \"Resource Not Found\",   \"detail\": \"No project with UUID 0976ef46-e4a0-4be4-8b0e-00e9d3625c1f exists\" } ```  ## Authentication  Two credential types are accepted:  - **API key** via the `X-Api-Key` request header. API keys are   long-lived and intended for machine-to-machine integrations. - **Bearer token** via the `Authorization: Bearer <token>` request   header. Bearer tokens are short-lived, user-bound, and opaque   server-issued session tokens.  Bearer tokens are obtained from one of the following endpoints:  - `POST /api/v1/user/login` - `POST /api/v1/user/oidc/login`  API v2 does not yet provide its own login endpoints; clients use the v1 endpoints above to acquire tokens and then call v2 with them.  Tokens are valid for 8 hours by default and **cannot be refreshed**. Clients must re-authenticate once a token expires.  Requests with missing or invalid credentials are rejected with `401 Unauthorized`.  ## Authorization  Access is gated by named permissions. Operations document the permission(s) they require; operations without a documented permission requirement only require authentication.  When the *Portfolio Access Control* feature is enabled (disabled by default), project-scoped operations additionally enforce per-project access via team membership. The `PORTFOLIO_ACCESS_CONTROL_BYPASS` permission grants access to all projects regardless of team mappings. When the feature is disabled, all authenticated callers holding the required permission can access all projects.  Authenticated requests that lack the required permission, or that target a project the caller cannot access, are rejected with `403 Forbidden`.  ## HTTP Methods  | Method   | Semantics                  | |----------|----------------------------| | `GET`    | Retrieve a resource        | | `POST`   | Create a new resource      | | `PUT`    | Update a resource          | | `PATCH`  | Partially update a resource| | `DELETE` | Delete a resource          |  ## Response Conventions  Create and update operations (`POST`, `PUT`, `PATCH`) do not return the full resource in the response. They return either no body, or only server-generated identifiers (e.g. a UUID). `POST` responses may include a `Location` header linking to the created resource.  Delete operations return `204 No Content` with no body.  ## Deprecations  Operations may be removed or replaced over time. When a response carries the `X-API-Deprecated: true` header, the operation that produced it is deprecated and may be removed in a future release. Clients should check for this header on every response and surface it (e.g. via a log warning) so that operators are aware of upcoming breakages. The respective operation's description points out which alternative operation(s) to use.  ## Internal operations  Operations under the `/internal` path prefix expose system internals and are reserved for first-party use. They are **not** part of the stable v2 API contract and may change or be removed without notice. Third-party clients should not depend on them.

API version: 2.0.0
Contact: dependencytrack@owasp.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dtv2

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// ComponentsAPIService ComponentsAPI service
type ComponentsAPIService service

type ApiCreateComponentRequest struct {
	ctx                    context.Context
	ApiService             *ComponentsAPIService
	createComponentRequest *CreateComponentRequest
}

func (r ApiCreateComponentRequest) CreateComponentRequest(createComponentRequest CreateComponentRequest) ApiCreateComponentRequest {
	r.createComponentRequest = &createComponentRequest
	return r
}

func (r ApiCreateComponentRequest) Execute() (*http.Response, error) {
	return r.ApiService.CreateComponentExecute(r)
}

/*
CreateComponent Creates a new component for the project

Requires permission `PORTFOLIO_MANAGEMENT` or `PORTFOLIO_MANAGEMENT_UPDATE`

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateComponentRequest
*/
func (a *ComponentsAPIService) CreateComponent(ctx context.Context) ApiCreateComponentRequest {
	return ApiCreateComponentRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *ComponentsAPIService) CreateComponentExecute(r ApiCreateComponentRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ComponentsAPIService.CreateComponent")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/components"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createComponentRequest == nil {
		return nil, reportError("createComponentRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/problem+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.createComponentRequest
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["apiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Api-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v CreateComponent400Response
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 409 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		var v ProblemDetails
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiListComponentsRequest struct {
	ctx                            context.Context
	ApiService                     *ComponentsAPIService
	groupContains                  *string
	nameContains                   *string
	versionContains                *string
	purlPrefix                     *string
	cpe                            *string
	swidTagIdContains              *string
	hashType                       *string
	hash                           *string
	packageArtifactPublishedSince  *int64
	packageArtifactPublishedBefore *int64
	projectState                   *ProjectState
	projectLatestVersion           *bool
	expand                         *[]string
	limit                          *int32
	pageToken                      *string
	sortDirection                  *SortDirection
	sortBy                         *string
}

// Filter by group (substring match)
func (r ApiListComponentsRequest) GroupContains(groupContains string) ApiListComponentsRequest {
	r.groupContains = &groupContains
	return r
}

// Filter by name (substring match)
func (r ApiListComponentsRequest) NameContains(nameContains string) ApiListComponentsRequest {
	r.nameContains = &nameContains
	return r
}

// Filter by version (substring match)
func (r ApiListComponentsRequest) VersionContains(versionContains string) ApiListComponentsRequest {
	r.versionContains = &versionContains
	return r
}

// Filter by PURL (prefix match).  Must be a valid PURL, with at least &#x60;pkg:&lt;ecosystem&gt;/&lt;name&gt;&#x60; populated.
func (r ApiListComponentsRequest) PurlPrefix(purlPrefix string) ApiListComponentsRequest {
	r.purlPrefix = &purlPrefix
	return r
}

// Filter by CPE (exact match).  Must be a valid CPE.
func (r ApiListComponentsRequest) Cpe(cpe string) ApiListComponentsRequest {
	r.cpe = &cpe
	return r
}

// Filter by SWID Tag ID (substring match)
func (r ApiListComponentsRequest) SwidTagIdContains(swidTagIdContains string) ApiListComponentsRequest {
	r.swidTagIdContains = &swidTagIdContains
	return r
}

// The hash type to filter by
func (r ApiListComponentsRequest) HashType(hashType string) ApiListComponentsRequest {
	r.hashType = &hashType
	return r
}

// Filter by hash value (exact match).  Requires &#x60;hash_type&#x60; to be set.
func (r ApiListComponentsRequest) Hash(hash string) ApiListComponentsRequest {
	r.hash = &hash
	return r
}

// Filter by package artifact publish date (inclusive lower bound).  Note that components without resolved package artifact metadata, or whose upstream repository did not report a publication date, are excluded whenever &#x60;package_artifact_published_since&#x60; or &#x60;package_artifact_published_before&#x60; is set.
func (r ApiListComponentsRequest) PackageArtifactPublishedSince(packageArtifactPublishedSince int64) ApiListComponentsRequest {
	r.packageArtifactPublishedSince = &packageArtifactPublishedSince
	return r
}

// Filter by package artifact publish date (exclusive upper bound).
func (r ApiListComponentsRequest) PackageArtifactPublishedBefore(packageArtifactPublishedBefore int64) ApiListComponentsRequest {
	r.packageArtifactPublishedBefore = &packageArtifactPublishedBefore
	return r
}

// Filter by the state of the project that the component belongs to.  Omit to include components from projects in any state.
func (r ApiListComponentsRequest) ProjectState(projectState ProjectState) ApiListComponentsRequest {
	r.projectState = &projectState
	return r
}

// Filter by whether the project the component belongs to is flagged as the latest version.  When &#x60;true&#x60;, only components from latest-version projects are returned. When &#x60;false&#x60;, only components from non-latest projects are returned. Omit to include components regardless of the flag.
func (r ApiListComponentsRequest) ProjectLatestVersion(projectLatestVersion bool) ApiListComponentsRequest {
	r.projectLatestVersion = &projectLatestVersion
	return r
}

// Optional fields to include in each component response item. Unknown values are silently ignored.
func (r ApiListComponentsRequest) Expand(expand []string) ApiListComponentsRequest {
	r.expand = &expand
	return r
}

// Maximum number of items to retrieve from the collection
func (r ApiListComponentsRequest) Limit(limit int32) ApiListComponentsRequest {
	r.limit = &limit
	return r
}

// Opaque token pointing to a specific position in a collection
func (r ApiListComponentsRequest) PageToken(pageToken string) ApiListComponentsRequest {
	r.pageToken = &pageToken
	return r
}

func (r ApiListComponentsRequest) SortDirection(sortDirection SortDirection) ApiListComponentsRequest {
	r.sortDirection = &sortDirection
	return r
}

// Field to sort by. Refer to the operation description for information about which fields are sortable.
func (r ApiListComponentsRequest) SortBy(sortBy string) ApiListComponentsRequest {
	r.sortBy = &sortBy
	return r
}

func (r ApiListComponentsRequest) Execute() (*ListComponentsResponse, *http.Response, error) {
	return r.ApiService.ListComponentsExecute(r)
}

/*
ListComponents List all components

Retrieves a list of all components matching the provided filter criteria.

Text filters are case-insensitive.

### Sortable fields

Sorting is supported for the following fields:

* `name`
* `group`
* `last_inherited_risk_score`

### Expandable fields

The following fields can be included via `expand`:

* `metrics`
* `package_metadata`
* `package_artifact_metadata`

Requires permission `VIEW_PORTFOLIO`

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiListComponentsRequest
*/
func (a *ComponentsAPIService) ListComponents(ctx context.Context) ApiListComponentsRequest {
	return ApiListComponentsRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return ListComponentsResponse
func (a *ComponentsAPIService) ListComponentsExecute(r ApiListComponentsRequest) (*ListComponentsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ListComponentsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ComponentsAPIService.ListComponents")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/components"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.groupContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "group_contains", r.groupContains, "form", "")
	}
	if r.nameContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name_contains", r.nameContains, "form", "")
	}
	if r.versionContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "version_contains", r.versionContains, "form", "")
	}
	if r.purlPrefix != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "purl_prefix", r.purlPrefix, "form", "")
	}
	if r.cpe != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "cpe", r.cpe, "form", "")
	}
	if r.swidTagIdContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "swid_tag_id_contains", r.swidTagIdContains, "form", "")
	}
	if r.hashType != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "hash_type", r.hashType, "form", "")
	}
	if r.hash != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "hash", r.hash, "form", "")
	}
	if r.packageArtifactPublishedSince != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "package_artifact_published_since", r.packageArtifactPublishedSince, "form", "")
	}
	if r.packageArtifactPublishedBefore != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "package_artifact_published_before", r.packageArtifactPublishedBefore, "form", "")
	}
	if r.projectState != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "project_state", r.projectState, "form", "")
	}
	if r.projectLatestVersion != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "project_latest_version", r.projectLatestVersion, "form", "")
	}
	if r.expand != nil {
		t := *r.expand
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "expand", s.Index(i).Interface(), "form", "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "expand", t, "form", "multi")
		}
	}
	if r.limit != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", r.limit, "form", "")
	} else {
		var defaultValue int32 = 100
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", defaultValue, "form", "")
		r.limit = &defaultValue
	}
	if r.pageToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "page_token", r.pageToken, "form", "")
	}
	if r.sortDirection != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "sort_direction", r.sortDirection, "form", "")
	}
	if r.sortBy != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "sort_by", r.sortBy, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json", "application/problem+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["apiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Api-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v ListComponents400Response
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v ProblemDetails
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
