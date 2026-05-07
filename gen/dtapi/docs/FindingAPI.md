# \FindingAPI

All URIs are relative to *https://hyades-api.iris-flair-alpha.vlair-staging.defra01.iris-sensing.net/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AnalyzePortfolio**](FindingAPI.md#AnalyzePortfolio) | **Post** /v1/finding/portfolio/analyze | Triggers Vulnerability Analysis for the entire portfolio
[**AnalyzeProject**](FindingAPI.md#AnalyzeProject) | **Post** /v1/finding/project/{uuid}/analyze | Triggers Vulnerability Analysis on a specific project
[**ExportFindingsByProject**](FindingAPI.md#ExportFindingsByProject) | **Get** /v1/finding/project/{uuid}/export | Returns the findings for the specified project as FPF
[**GetAllFindings**](FindingAPI.md#GetAllFindings) | **Get** /v1/finding/grouped | Returns a list of all findings grouped by vulnerability
[**GetAllFindings1**](FindingAPI.md#GetAllFindings1) | **Get** /v1/finding | Returns a list of all findings
[**GetFindingsByProject**](FindingAPI.md#GetFindingsByProject) | **Get** /v1/finding/project/{uuid} | Returns a list of all findings for a specific project or generates SARIF file if Accept: application/sarif+json header is provided



## AnalyzePortfolio

> AnalyzePortfolio(ctx).Execute()

Triggers Vulnerability Analysis for the entire portfolio



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FindingAPI.AnalyzePortfolio(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.AnalyzePortfolio``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAnalyzePortfolioRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AnalyzeProject

> BomUploadResponse AnalyzeProject(ctx, uuid).Execute()

Triggers Vulnerability Analysis on a specific project



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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the project to analyze

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FindingAPI.AnalyzeProject(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.AnalyzeProject``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AnalyzeProject`: BomUploadResponse
	fmt.Fprintf(os.Stdout, "Response from `FindingAPI.AnalyzeProject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the project to analyze | 

### Other Parameters

Other parameters are passed through a pointer to a apiAnalyzeProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BomUploadResponse**](BomUploadResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportFindingsByProject

> string ExportFindingsByProject(ctx, uuid).Execute()

Returns the findings for the specified project as FPF



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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the project

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FindingAPI.ExportFindingsByProject(context.Background(), uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.ExportFindingsByProject``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportFindingsByProject`: string
	fmt.Fprintf(os.Stdout, "Response from `FindingAPI.ExportFindingsByProject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the project | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportFindingsByProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**string**

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllFindings

> []Finding GetAllFindings(ctx).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).ShowInactive(showInactive).Severity(severity).PublishDateFrom(publishDateFrom).PublishDateTo(publishDateTo).TextSearchField(textSearchField).TextSearchInput(textSearchInput).Cvssv2From(cvssv2From).Cvssv2To(cvssv2To).Cvssv3From(cvssv3From).Cvssv3To(cvssv3To).EpssFrom(epssFrom).EpssTo(epssTo).EpssPercentileFrom(epssPercentileFrom).EpssPercentileTo(epssPercentileTo).OccurrencesFrom(occurrencesFrom).OccurrencesTo(occurrencesTo).Execute()

Returns a list of all findings grouped by vulnerability



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
	pageNumber := "pageNumber_example" // string | The page to return. To be used in conjunction with <code>pageSize</code>. (optional) (default to "1")
	pageSize := "pageSize_example" // string | Number of elements to return per page. To be used in conjunction with <code>pageNumber</code>. (optional) (default to "100")
	offset := "offset_example" // string | Offset to start returning elements from. To be used in conjunction with <code>limit</code>. (optional)
	limit := "limit_example" // string | Number of elements to return per page. To be used in conjunction with <code>offset</code>. (optional)
	sortName := "sortName_example" // string | Name of the resource field to sort on. (optional)
	sortOrder := "sortOrder_example" // string | Ordering of items when sorting with <code>sortName</code>. (optional)
	showInactive := true // bool | Show inactive projects (optional)
	severity := "severity_example" // string | Filter by severity (optional)
	publishDateFrom := "publishDateFrom_example" // string | Filter published from this date (optional)
	publishDateTo := "publishDateTo_example" // string | Filter published to this date (optional)
	textSearchField := "textSearchField_example" // string | Filter the text input in these fields (optional)
	textSearchInput := "textSearchInput_example" // string | Filter by this text input (optional)
	cvssv2From := "cvssv2From_example" // string | Filter CVSSv2 from this value (optional)
	cvssv2To := "cvssv2To_example" // string | Filter CVSSv2 to this value (optional)
	cvssv3From := "cvssv3From_example" // string | Filter CVSSv3 from this value (optional)
	cvssv3To := "cvssv3To_example" // string | Filter CVSSv3 to this value (optional)
	epssFrom := "epssFrom_example" // string | Filter EPSS from this value (optional)
	epssTo := "epssTo_example" // string | Filter EPSS to this value (optional)
	epssPercentileFrom := "epssPercentileFrom_example" // string | Filter EPSS Percentile from this value (optional)
	epssPercentileTo := "epssPercentileTo_example" // string | Filter EPSS Percentile to this value (optional)
	occurrencesFrom := "occurrencesFrom_example" // string | Filter occurrences in projects from this value (optional)
	occurrencesTo := "occurrencesTo_example" // string | Filter occurrences in projects to this value (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FindingAPI.GetAllFindings(context.Background()).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).ShowInactive(showInactive).Severity(severity).PublishDateFrom(publishDateFrom).PublishDateTo(publishDateTo).TextSearchField(textSearchField).TextSearchInput(textSearchInput).Cvssv2From(cvssv2From).Cvssv2To(cvssv2To).Cvssv3From(cvssv3From).Cvssv3To(cvssv3To).EpssFrom(epssFrom).EpssTo(epssTo).EpssPercentileFrom(epssPercentileFrom).EpssPercentileTo(epssPercentileTo).OccurrencesFrom(occurrencesFrom).OccurrencesTo(occurrencesTo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.GetAllFindings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllFindings`: []Finding
	fmt.Fprintf(os.Stdout, "Response from `FindingAPI.GetAllFindings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAllFindingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageNumber** | **string** | The page to return. To be used in conjunction with &lt;code&gt;pageSize&lt;/code&gt;. | [default to &quot;1&quot;]
 **pageSize** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;pageNumber&lt;/code&gt;. | [default to &quot;100&quot;]
 **offset** | **string** | Offset to start returning elements from. To be used in conjunction with &lt;code&gt;limit&lt;/code&gt;. | 
 **limit** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;offset&lt;/code&gt;. | 
 **sortName** | **string** | Name of the resource field to sort on. | 
 **sortOrder** | **string** | Ordering of items when sorting with &lt;code&gt;sortName&lt;/code&gt;. | 
 **showInactive** | **bool** | Show inactive projects | 
 **severity** | **string** | Filter by severity | 
 **publishDateFrom** | **string** | Filter published from this date | 
 **publishDateTo** | **string** | Filter published to this date | 
 **textSearchField** | **string** | Filter the text input in these fields | 
 **textSearchInput** | **string** | Filter by this text input | 
 **cvssv2From** | **string** | Filter CVSSv2 from this value | 
 **cvssv2To** | **string** | Filter CVSSv2 to this value | 
 **cvssv3From** | **string** | Filter CVSSv3 from this value | 
 **cvssv3To** | **string** | Filter CVSSv3 to this value | 
 **epssFrom** | **string** | Filter EPSS from this value | 
 **epssTo** | **string** | Filter EPSS to this value | 
 **epssPercentileFrom** | **string** | Filter EPSS Percentile from this value | 
 **epssPercentileTo** | **string** | Filter EPSS Percentile to this value | 
 **occurrencesFrom** | **string** | Filter occurrences in projects from this value | 
 **occurrencesTo** | **string** | Filter occurrences in projects to this value | 

### Return type

[**[]Finding**](Finding.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllFindings1

> []Finding GetAllFindings1(ctx).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).ShowInactive(showInactive).ShowSuppressed(showSuppressed).Severity(severity).AnalysisStatus(analysisStatus).VendorResponse(vendorResponse).PublishDateFrom(publishDateFrom).PublishDateTo(publishDateTo).AttributedOnDateFrom(attributedOnDateFrom).AttributedOnDateTo(attributedOnDateTo).TextSearchField(textSearchField).TextSearchInput(textSearchInput).Cvssv2From(cvssv2From).Cvssv2To(cvssv2To).Cvssv3From(cvssv3From).Cvssv3To(cvssv3To).EpssFrom(epssFrom).EpssTo(epssTo).EpssPercentileFrom(epssPercentileFrom).EpssPercentileTo(epssPercentileTo).Execute()

Returns a list of all findings



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
	pageNumber := "pageNumber_example" // string | The page to return. To be used in conjunction with <code>pageSize</code>. (optional) (default to "1")
	pageSize := "pageSize_example" // string | Number of elements to return per page. To be used in conjunction with <code>pageNumber</code>. (optional) (default to "100")
	offset := "offset_example" // string | Offset to start returning elements from. To be used in conjunction with <code>limit</code>. (optional)
	limit := "limit_example" // string | Number of elements to return per page. To be used in conjunction with <code>offset</code>. (optional)
	sortName := "sortName_example" // string | Name of the resource field to sort on. (optional)
	sortOrder := "sortOrder_example" // string | Ordering of items when sorting with <code>sortName</code>. (optional)
	showInactive := true // bool | Show inactive projects (optional)
	showSuppressed := true // bool | Show suppressed findings (optional)
	severity := "severity_example" // string | Filter by severity (optional)
	analysisStatus := "analysisStatus_example" // string | Filter by analysis status (optional)
	vendorResponse := "vendorResponse_example" // string | Filter by vendor response (optional)
	publishDateFrom := "publishDateFrom_example" // string | Filter published from this date (optional)
	publishDateTo := "publishDateTo_example" // string | Filter published to this date (optional)
	attributedOnDateFrom := "attributedOnDateFrom_example" // string | Filter attributed on from this date (optional)
	attributedOnDateTo := "attributedOnDateTo_example" // string | Filter attributed on to this date (optional)
	textSearchField := "textSearchField_example" // string | Filter the text input in these fields (optional)
	textSearchInput := "textSearchInput_example" // string | Filter by this text input (optional)
	cvssv2From := "cvssv2From_example" // string | Filter CVSSv2 from this value (optional)
	cvssv2To := "cvssv2To_example" // string | Filter CVSSv2 from this Value (optional)
	cvssv3From := "cvssv3From_example" // string | Filter CVSSv3 from this value (optional)
	cvssv3To := "cvssv3To_example" // string | Filter CVSSv3 from this Value (optional)
	epssFrom := "epssFrom_example" // string | Filter EPSS from this value (optional)
	epssTo := "epssTo_example" // string | Filter EPSS to this value (optional)
	epssPercentileFrom := "epssPercentileFrom_example" // string | Filter EPSS Percentile from this value (optional)
	epssPercentileTo := "epssPercentileTo_example" // string | Filter EPSS Percentile to this value (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FindingAPI.GetAllFindings1(context.Background()).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).ShowInactive(showInactive).ShowSuppressed(showSuppressed).Severity(severity).AnalysisStatus(analysisStatus).VendorResponse(vendorResponse).PublishDateFrom(publishDateFrom).PublishDateTo(publishDateTo).AttributedOnDateFrom(attributedOnDateFrom).AttributedOnDateTo(attributedOnDateTo).TextSearchField(textSearchField).TextSearchInput(textSearchInput).Cvssv2From(cvssv2From).Cvssv2To(cvssv2To).Cvssv3From(cvssv3From).Cvssv3To(cvssv3To).EpssFrom(epssFrom).EpssTo(epssTo).EpssPercentileFrom(epssPercentileFrom).EpssPercentileTo(epssPercentileTo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.GetAllFindings1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllFindings1`: []Finding
	fmt.Fprintf(os.Stdout, "Response from `FindingAPI.GetAllFindings1`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAllFindings1Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageNumber** | **string** | The page to return. To be used in conjunction with &lt;code&gt;pageSize&lt;/code&gt;. | [default to &quot;1&quot;]
 **pageSize** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;pageNumber&lt;/code&gt;. | [default to &quot;100&quot;]
 **offset** | **string** | Offset to start returning elements from. To be used in conjunction with &lt;code&gt;limit&lt;/code&gt;. | 
 **limit** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;offset&lt;/code&gt;. | 
 **sortName** | **string** | Name of the resource field to sort on. | 
 **sortOrder** | **string** | Ordering of items when sorting with &lt;code&gt;sortName&lt;/code&gt;. | 
 **showInactive** | **bool** | Show inactive projects | 
 **showSuppressed** | **bool** | Show suppressed findings | 
 **severity** | **string** | Filter by severity | 
 **analysisStatus** | **string** | Filter by analysis status | 
 **vendorResponse** | **string** | Filter by vendor response | 
 **publishDateFrom** | **string** | Filter published from this date | 
 **publishDateTo** | **string** | Filter published to this date | 
 **attributedOnDateFrom** | **string** | Filter attributed on from this date | 
 **attributedOnDateTo** | **string** | Filter attributed on to this date | 
 **textSearchField** | **string** | Filter the text input in these fields | 
 **textSearchInput** | **string** | Filter by this text input | 
 **cvssv2From** | **string** | Filter CVSSv2 from this value | 
 **cvssv2To** | **string** | Filter CVSSv2 from this Value | 
 **cvssv3From** | **string** | Filter CVSSv3 from this value | 
 **cvssv3To** | **string** | Filter CVSSv3 from this Value | 
 **epssFrom** | **string** | Filter EPSS from this value | 
 **epssTo** | **string** | Filter EPSS to this value | 
 **epssPercentileFrom** | **string** | Filter EPSS Percentile from this value | 
 **epssPercentileTo** | **string** | Filter EPSS Percentile to this value | 

### Return type

[**[]Finding**](Finding.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetFindingsByProject

> []Finding GetFindingsByProject(ctx, uuid).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).Suppressed(suppressed).Source(source).Accept(accept).HasAnalysis(hasAnalysis).Execute()

Returns a list of all findings for a specific project or generates SARIF file if Accept: application/sarif+json header is provided



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
	uuid := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The UUID of the project
	pageNumber := "pageNumber_example" // string | The page to return. To be used in conjunction with <code>pageSize</code>. (optional) (default to "1")
	pageSize := "pageSize_example" // string | Number of elements to return per page. To be used in conjunction with <code>pageNumber</code>. (optional) (default to "100")
	offset := "offset_example" // string | Offset to start returning elements from. To be used in conjunction with <code>limit</code>. (optional)
	limit := "limit_example" // string | Number of elements to return per page. To be used in conjunction with <code>offset</code>. (optional)
	sortName := "sortName_example" // string | Name of the resource field to sort on. (optional)
	sortOrder := "sortOrder_example" // string | Ordering of items when sorting with <code>sortName</code>. (optional)
	suppressed := true // bool | Optionally includes suppressed findings (optional)
	source := "source_example" // string | Optionally limit findings to specific sources of vulnerability intelligence (optional)
	accept := "accept_example" // string |  (optional)
	hasAnalysis := true // bool | Whether to include only projects with existing analysis. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FindingAPI.GetFindingsByProject(context.Background(), uuid).PageNumber(pageNumber).PageSize(pageSize).Offset(offset).Limit(limit).SortName(sortName).SortOrder(sortOrder).Suppressed(suppressed).Source(source).Accept(accept).HasAnalysis(hasAnalysis).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FindingAPI.GetFindingsByProject``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetFindingsByProject`: []Finding
	fmt.Fprintf(os.Stdout, "Response from `FindingAPI.GetFindingsByProject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**uuid** | **string** | The UUID of the project | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetFindingsByProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pageNumber** | **string** | The page to return. To be used in conjunction with &lt;code&gt;pageSize&lt;/code&gt;. | [default to &quot;1&quot;]
 **pageSize** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;pageNumber&lt;/code&gt;. | [default to &quot;100&quot;]
 **offset** | **string** | Offset to start returning elements from. To be used in conjunction with &lt;code&gt;limit&lt;/code&gt;. | 
 **limit** | **string** | Number of elements to return per page. To be used in conjunction with &lt;code&gt;offset&lt;/code&gt;. | 
 **sortName** | **string** | Name of the resource field to sort on. | 
 **sortOrder** | **string** | Ordering of items when sorting with &lt;code&gt;sortName&lt;/code&gt;. | 
 **suppressed** | **bool** | Optionally includes suppressed findings | 
 **source** | **string** | Optionally limit findings to specific sources of vulnerability intelligence | 
 **accept** | **string** |  | 
 **hasAnalysis** | **bool** | Whether to include only projects with existing analysis. | 

### Return type

[**[]Finding**](Finding.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/sarif+json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

