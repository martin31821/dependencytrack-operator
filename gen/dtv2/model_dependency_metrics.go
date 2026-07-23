/*
OWASP Dependency-Track

# REST API of OWASP Dependency-Track  ## Pagination  This API implements token-based pagination. Collection responses have the following structure:  ```json {   \"items\": [...],   \"next_page_token\": \"abcdefg\",   \"total\": {     \"count\": 100,     \"type\": \"EXACT\"   } } ```  `next_page_token` is present when more items exist, and absent otherwise. To fetch the next page, pass it as the `page_token` query parameter.  To navigate backwards, clients should keep track of previous page tokens as they paginate through collections. The API does *not* provide backward navigation!  Collections that support sorting will only consider the `sort_by` and `sort_direction` query parameters for the request of the first page. For subsequent pages, sorting preferences are bound to the page token.  Page tokens are opaque strings. Clients should not try to interpret or generate them. Their format may change without notice.  The `total` object discloses how many items exist in the collection *across all pages*. Because counting is expensive, some collections that hold *a lot* of items may return partial counts (type `AT_LEAST`) instead of exact counts (type `EXACT`). Which type to expect is usually documented in the operation's description.  ## Sorting  Items in a collection can be sorted using the `sort_by` and `sort_direction` query parameters. Which fields are sortable is documented in the respective operation's description.  Note that if no sortable fields are documented for an operation, sorting is not supported *at all*.  ## Field expansion  Some collection endpoints support an `expand` query parameter. Passing an expand value includes optional fields in each response item that are omitted by default, typically because they are expensive to compute and only needed in specific contexts.  Valid `expand` values for an endpoint are listed in its operation description. Unknown values are silently ignored.  ## Errors  All error responses use the `application/problem+json` media type as defined in [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html).  Example:  ```json {   \"type\": \"about:blank\",   \"status\": 404,   \"title\": \"Resource Not Found\",   \"detail\": \"No project with UUID 0976ef46-e4a0-4be4-8b0e-00e9d3625c1f exists\" } ```  ## Authentication  Two credential types are accepted:  - **API key** via the `X-Api-Key` request header. API keys are   long-lived and intended for machine-to-machine integrations. - **Bearer token** via the `Authorization: Bearer <token>` request   header. Bearer tokens are short-lived, user-bound, and opaque   server-issued session tokens.  Bearer tokens are obtained from one of the following endpoints:  - `POST /api/v1/user/login` - `POST /api/v1/user/oidc/login`  API v2 does not yet provide its own login endpoints; clients use the v1 endpoints above to acquire tokens and then call v2 with them.  Tokens are valid for 8 hours by default and **cannot be refreshed**. Clients must re-authenticate once a token expires.  Requests with missing or invalid credentials are rejected with `401 Unauthorized`.  ## Authorization  Access is gated by named permissions. Operations document the permission(s) they require; operations without a documented permission requirement only require authentication.  When the *Portfolio Access Control* feature is enabled (disabled by default), project-scoped operations additionally enforce per-project access via team membership. The `PORTFOLIO_ACCESS_CONTROL_BYPASS` permission grants access to all projects regardless of team mappings. When the feature is disabled, all authenticated callers holding the required permission can access all projects.  Authenticated requests that lack the required permission, or that target a project the caller cannot access, are rejected with `403 Forbidden`.  ## HTTP Methods  | Method   | Semantics                  | |----------|----------------------------| | `GET`    | Retrieve a resource        | | `POST`   | Create a new resource      | | `PUT`    | Update a resource          | | `PATCH`  | Partially update a resource| | `DELETE` | Delete a resource          |  ## Response Conventions  Create and update operations (`POST`, `PUT`, `PATCH`) do not return the full resource in the response. They return either no body, or only server-generated identifiers (e.g. a UUID). `POST` responses may include a `Location` header linking to the created resource.  Delete operations return `204 No Content` with no body.  ## Deprecations  Operations may be removed or replaced over time. When a response carries the `X-API-Deprecated: true` header, the operation that produced it is deprecated and may be removed in a future release. Clients should check for this header on every response and surface it (e.g. via a log warning) so that operators are aware of upcoming breakages. The respective operation's description points out which alternative operation(s) to use.  ## Internal operations  Operations under the `/internal` path prefix expose system internals and are reserved for first-party use. They are **not** part of the stable v2 API contract and may change or be removed without notice. Third-party clients should not depend on them.

API version: 2.0.0
Contact: dependencytrack@owasp.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dtv2

import (
	"encoding/json"
)

// checks if the DependencyMetrics type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DependencyMetrics{}

// DependencyMetrics struct for DependencyMetrics
type DependencyMetrics struct {
	Critical                             *int32   `json:"critical,omitempty"`
	High                                 *int32   `json:"high,omitempty"`
	Medium                               *int32   `json:"medium,omitempty"`
	Low                                  *int32   `json:"low,omitempty"`
	Unassigned                           *int32   `json:"unassigned,omitempty"`
	Vulnerabilities                      *int32   `json:"vulnerabilities,omitempty"`
	Suppressed                           *int32   `json:"suppressed,omitempty"`
	InheritedRiskScore                   *float64 `json:"inherited_risk_score,omitempty"`
	FindingsTotal                        *int32   `json:"findings_total,omitempty"`
	FindingsAudited                      *int32   `json:"findings_audited,omitempty"`
	FindingsUnaudited                    *int32   `json:"findings_unaudited,omitempty"`
	PolicyViolationsFail                 *int32   `json:"policy_violations_fail,omitempty"`
	PolicyViolationsWarn                 *int32   `json:"policy_violations_warn,omitempty"`
	PolicyViolationsInfo                 *int32   `json:"policy_violations_info,omitempty"`
	PolicyViolationsTotal                *int32   `json:"policy_violations_total,omitempty"`
	PolicyViolationsAudited              *int32   `json:"policy_violations_audited,omitempty"`
	PolicyViolationsUnaudited            *int32   `json:"policy_violations_unaudited,omitempty"`
	PolicyViolationsSecurityTotal        *int32   `json:"policy_violations_security_total,omitempty"`
	PolicyViolationsSecurityAudited      *int32   `json:"policy_violations_security_audited,omitempty"`
	PolicyViolationsSecurityUnaudited    *int32   `json:"policy_violations_security_unaudited,omitempty"`
	PolicyViolationsLicenseTotal         *int32   `json:"policy_violations_license_total,omitempty"`
	PolicyViolationsLicenseAudited       *int32   `json:"policy_violations_license_audited,omitempty"`
	PolicyViolationsLicenseUnaudited     *int32   `json:"policy_violations_license_unaudited,omitempty"`
	PolicyViolationsOperationalTotal     *int32   `json:"policy_violations_operational_total,omitempty"`
	PolicyViolationsOperationalAudited   *int32   `json:"policy_violations_operational_audited,omitempty"`
	PolicyViolationsOperationalUnaudited *int32   `json:"policy_violations_operational_unaudited,omitempty"`
}

// NewDependencyMetrics instantiates a new DependencyMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDependencyMetrics() *DependencyMetrics {
	this := DependencyMetrics{}
	return &this
}

// NewDependencyMetricsWithDefaults instantiates a new DependencyMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDependencyMetricsWithDefaults() *DependencyMetrics {
	this := DependencyMetrics{}
	return &this
}

// GetCritical returns the Critical field value if set, zero value otherwise.
func (o *DependencyMetrics) GetCritical() int32 {
	if o == nil || IsNil(o.Critical) {
		var ret int32
		return ret
	}
	return *o.Critical
}

// GetCriticalOk returns a tuple with the Critical field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetCriticalOk() (*int32, bool) {
	if o == nil || IsNil(o.Critical) {
		return nil, false
	}
	return o.Critical, true
}

// HasCritical returns a boolean if a field has been set.
func (o *DependencyMetrics) HasCritical() bool {
	if o != nil && !IsNil(o.Critical) {
		return true
	}

	return false
}

// SetCritical gets a reference to the given int32 and assigns it to the Critical field.
func (o *DependencyMetrics) SetCritical(v int32) {
	o.Critical = &v
}

// GetHigh returns the High field value if set, zero value otherwise.
func (o *DependencyMetrics) GetHigh() int32 {
	if o == nil || IsNil(o.High) {
		var ret int32
		return ret
	}
	return *o.High
}

// GetHighOk returns a tuple with the High field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetHighOk() (*int32, bool) {
	if o == nil || IsNil(o.High) {
		return nil, false
	}
	return o.High, true
}

// HasHigh returns a boolean if a field has been set.
func (o *DependencyMetrics) HasHigh() bool {
	if o != nil && !IsNil(o.High) {
		return true
	}

	return false
}

// SetHigh gets a reference to the given int32 and assigns it to the High field.
func (o *DependencyMetrics) SetHigh(v int32) {
	o.High = &v
}

// GetMedium returns the Medium field value if set, zero value otherwise.
func (o *DependencyMetrics) GetMedium() int32 {
	if o == nil || IsNil(o.Medium) {
		var ret int32
		return ret
	}
	return *o.Medium
}

// GetMediumOk returns a tuple with the Medium field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetMediumOk() (*int32, bool) {
	if o == nil || IsNil(o.Medium) {
		return nil, false
	}
	return o.Medium, true
}

// HasMedium returns a boolean if a field has been set.
func (o *DependencyMetrics) HasMedium() bool {
	if o != nil && !IsNil(o.Medium) {
		return true
	}

	return false
}

// SetMedium gets a reference to the given int32 and assigns it to the Medium field.
func (o *DependencyMetrics) SetMedium(v int32) {
	o.Medium = &v
}

// GetLow returns the Low field value if set, zero value otherwise.
func (o *DependencyMetrics) GetLow() int32 {
	if o == nil || IsNil(o.Low) {
		var ret int32
		return ret
	}
	return *o.Low
}

// GetLowOk returns a tuple with the Low field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetLowOk() (*int32, bool) {
	if o == nil || IsNil(o.Low) {
		return nil, false
	}
	return o.Low, true
}

// HasLow returns a boolean if a field has been set.
func (o *DependencyMetrics) HasLow() bool {
	if o != nil && !IsNil(o.Low) {
		return true
	}

	return false
}

// SetLow gets a reference to the given int32 and assigns it to the Low field.
func (o *DependencyMetrics) SetLow(v int32) {
	o.Low = &v
}

// GetUnassigned returns the Unassigned field value if set, zero value otherwise.
func (o *DependencyMetrics) GetUnassigned() int32 {
	if o == nil || IsNil(o.Unassigned) {
		var ret int32
		return ret
	}
	return *o.Unassigned
}

// GetUnassignedOk returns a tuple with the Unassigned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetUnassignedOk() (*int32, bool) {
	if o == nil || IsNil(o.Unassigned) {
		return nil, false
	}
	return o.Unassigned, true
}

// HasUnassigned returns a boolean if a field has been set.
func (o *DependencyMetrics) HasUnassigned() bool {
	if o != nil && !IsNil(o.Unassigned) {
		return true
	}

	return false
}

// SetUnassigned gets a reference to the given int32 and assigns it to the Unassigned field.
func (o *DependencyMetrics) SetUnassigned(v int32) {
	o.Unassigned = &v
}

// GetVulnerabilities returns the Vulnerabilities field value if set, zero value otherwise.
func (o *DependencyMetrics) GetVulnerabilities() int32 {
	if o == nil || IsNil(o.Vulnerabilities) {
		var ret int32
		return ret
	}
	return *o.Vulnerabilities
}

// GetVulnerabilitiesOk returns a tuple with the Vulnerabilities field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetVulnerabilitiesOk() (*int32, bool) {
	if o == nil || IsNil(o.Vulnerabilities) {
		return nil, false
	}
	return o.Vulnerabilities, true
}

// HasVulnerabilities returns a boolean if a field has been set.
func (o *DependencyMetrics) HasVulnerabilities() bool {
	if o != nil && !IsNil(o.Vulnerabilities) {
		return true
	}

	return false
}

// SetVulnerabilities gets a reference to the given int32 and assigns it to the Vulnerabilities field.
func (o *DependencyMetrics) SetVulnerabilities(v int32) {
	o.Vulnerabilities = &v
}

// GetSuppressed returns the Suppressed field value if set, zero value otherwise.
func (o *DependencyMetrics) GetSuppressed() int32 {
	if o == nil || IsNil(o.Suppressed) {
		var ret int32
		return ret
	}
	return *o.Suppressed
}

// GetSuppressedOk returns a tuple with the Suppressed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetSuppressedOk() (*int32, bool) {
	if o == nil || IsNil(o.Suppressed) {
		return nil, false
	}
	return o.Suppressed, true
}

// HasSuppressed returns a boolean if a field has been set.
func (o *DependencyMetrics) HasSuppressed() bool {
	if o != nil && !IsNil(o.Suppressed) {
		return true
	}

	return false
}

// SetSuppressed gets a reference to the given int32 and assigns it to the Suppressed field.
func (o *DependencyMetrics) SetSuppressed(v int32) {
	o.Suppressed = &v
}

// GetInheritedRiskScore returns the InheritedRiskScore field value if set, zero value otherwise.
func (o *DependencyMetrics) GetInheritedRiskScore() float64 {
	if o == nil || IsNil(o.InheritedRiskScore) {
		var ret float64
		return ret
	}
	return *o.InheritedRiskScore
}

// GetInheritedRiskScoreOk returns a tuple with the InheritedRiskScore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetInheritedRiskScoreOk() (*float64, bool) {
	if o == nil || IsNil(o.InheritedRiskScore) {
		return nil, false
	}
	return o.InheritedRiskScore, true
}

// HasInheritedRiskScore returns a boolean if a field has been set.
func (o *DependencyMetrics) HasInheritedRiskScore() bool {
	if o != nil && !IsNil(o.InheritedRiskScore) {
		return true
	}

	return false
}

// SetInheritedRiskScore gets a reference to the given float64 and assigns it to the InheritedRiskScore field.
func (o *DependencyMetrics) SetInheritedRiskScore(v float64) {
	o.InheritedRiskScore = &v
}

// GetFindingsTotal returns the FindingsTotal field value if set, zero value otherwise.
func (o *DependencyMetrics) GetFindingsTotal() int32 {
	if o == nil || IsNil(o.FindingsTotal) {
		var ret int32
		return ret
	}
	return *o.FindingsTotal
}

// GetFindingsTotalOk returns a tuple with the FindingsTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetFindingsTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.FindingsTotal) {
		return nil, false
	}
	return o.FindingsTotal, true
}

// HasFindingsTotal returns a boolean if a field has been set.
func (o *DependencyMetrics) HasFindingsTotal() bool {
	if o != nil && !IsNil(o.FindingsTotal) {
		return true
	}

	return false
}

// SetFindingsTotal gets a reference to the given int32 and assigns it to the FindingsTotal field.
func (o *DependencyMetrics) SetFindingsTotal(v int32) {
	o.FindingsTotal = &v
}

// GetFindingsAudited returns the FindingsAudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetFindingsAudited() int32 {
	if o == nil || IsNil(o.FindingsAudited) {
		var ret int32
		return ret
	}
	return *o.FindingsAudited
}

// GetFindingsAuditedOk returns a tuple with the FindingsAudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetFindingsAuditedOk() (*int32, bool) {
	if o == nil || IsNil(o.FindingsAudited) {
		return nil, false
	}
	return o.FindingsAudited, true
}

// HasFindingsAudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasFindingsAudited() bool {
	if o != nil && !IsNil(o.FindingsAudited) {
		return true
	}

	return false
}

// SetFindingsAudited gets a reference to the given int32 and assigns it to the FindingsAudited field.
func (o *DependencyMetrics) SetFindingsAudited(v int32) {
	o.FindingsAudited = &v
}

// GetFindingsUnaudited returns the FindingsUnaudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetFindingsUnaudited() int32 {
	if o == nil || IsNil(o.FindingsUnaudited) {
		var ret int32
		return ret
	}
	return *o.FindingsUnaudited
}

// GetFindingsUnauditedOk returns a tuple with the FindingsUnaudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetFindingsUnauditedOk() (*int32, bool) {
	if o == nil || IsNil(o.FindingsUnaudited) {
		return nil, false
	}
	return o.FindingsUnaudited, true
}

// HasFindingsUnaudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasFindingsUnaudited() bool {
	if o != nil && !IsNil(o.FindingsUnaudited) {
		return true
	}

	return false
}

// SetFindingsUnaudited gets a reference to the given int32 and assigns it to the FindingsUnaudited field.
func (o *DependencyMetrics) SetFindingsUnaudited(v int32) {
	o.FindingsUnaudited = &v
}

// GetPolicyViolationsFail returns the PolicyViolationsFail field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsFail() int32 {
	if o == nil || IsNil(o.PolicyViolationsFail) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsFail
}

// GetPolicyViolationsFailOk returns a tuple with the PolicyViolationsFail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsFailOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsFail) {
		return nil, false
	}
	return o.PolicyViolationsFail, true
}

// HasPolicyViolationsFail returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsFail() bool {
	if o != nil && !IsNil(o.PolicyViolationsFail) {
		return true
	}

	return false
}

// SetPolicyViolationsFail gets a reference to the given int32 and assigns it to the PolicyViolationsFail field.
func (o *DependencyMetrics) SetPolicyViolationsFail(v int32) {
	o.PolicyViolationsFail = &v
}

// GetPolicyViolationsWarn returns the PolicyViolationsWarn field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsWarn() int32 {
	if o == nil || IsNil(o.PolicyViolationsWarn) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsWarn
}

// GetPolicyViolationsWarnOk returns a tuple with the PolicyViolationsWarn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsWarnOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsWarn) {
		return nil, false
	}
	return o.PolicyViolationsWarn, true
}

// HasPolicyViolationsWarn returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsWarn() bool {
	if o != nil && !IsNil(o.PolicyViolationsWarn) {
		return true
	}

	return false
}

// SetPolicyViolationsWarn gets a reference to the given int32 and assigns it to the PolicyViolationsWarn field.
func (o *DependencyMetrics) SetPolicyViolationsWarn(v int32) {
	o.PolicyViolationsWarn = &v
}

// GetPolicyViolationsInfo returns the PolicyViolationsInfo field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsInfo() int32 {
	if o == nil || IsNil(o.PolicyViolationsInfo) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsInfo
}

// GetPolicyViolationsInfoOk returns a tuple with the PolicyViolationsInfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsInfoOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsInfo) {
		return nil, false
	}
	return o.PolicyViolationsInfo, true
}

// HasPolicyViolationsInfo returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsInfo() bool {
	if o != nil && !IsNil(o.PolicyViolationsInfo) {
		return true
	}

	return false
}

// SetPolicyViolationsInfo gets a reference to the given int32 and assigns it to the PolicyViolationsInfo field.
func (o *DependencyMetrics) SetPolicyViolationsInfo(v int32) {
	o.PolicyViolationsInfo = &v
}

// GetPolicyViolationsTotal returns the PolicyViolationsTotal field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsTotal() int32 {
	if o == nil || IsNil(o.PolicyViolationsTotal) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsTotal
}

// GetPolicyViolationsTotalOk returns a tuple with the PolicyViolationsTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsTotal) {
		return nil, false
	}
	return o.PolicyViolationsTotal, true
}

// HasPolicyViolationsTotal returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsTotal() bool {
	if o != nil && !IsNil(o.PolicyViolationsTotal) {
		return true
	}

	return false
}

// SetPolicyViolationsTotal gets a reference to the given int32 and assigns it to the PolicyViolationsTotal field.
func (o *DependencyMetrics) SetPolicyViolationsTotal(v int32) {
	o.PolicyViolationsTotal = &v
}

// GetPolicyViolationsAudited returns the PolicyViolationsAudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsAudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsAudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsAudited
}

// GetPolicyViolationsAuditedOk returns a tuple with the PolicyViolationsAudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsAuditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsAudited) {
		return nil, false
	}
	return o.PolicyViolationsAudited, true
}

// HasPolicyViolationsAudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsAudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsAudited) {
		return true
	}

	return false
}

// SetPolicyViolationsAudited gets a reference to the given int32 and assigns it to the PolicyViolationsAudited field.
func (o *DependencyMetrics) SetPolicyViolationsAudited(v int32) {
	o.PolicyViolationsAudited = &v
}

// GetPolicyViolationsUnaudited returns the PolicyViolationsUnaudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsUnaudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsUnaudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsUnaudited
}

// GetPolicyViolationsUnauditedOk returns a tuple with the PolicyViolationsUnaudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsUnauditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsUnaudited) {
		return nil, false
	}
	return o.PolicyViolationsUnaudited, true
}

// HasPolicyViolationsUnaudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsUnaudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsUnaudited) {
		return true
	}

	return false
}

// SetPolicyViolationsUnaudited gets a reference to the given int32 and assigns it to the PolicyViolationsUnaudited field.
func (o *DependencyMetrics) SetPolicyViolationsUnaudited(v int32) {
	o.PolicyViolationsUnaudited = &v
}

// GetPolicyViolationsSecurityTotal returns the PolicyViolationsSecurityTotal field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsSecurityTotal() int32 {
	if o == nil || IsNil(o.PolicyViolationsSecurityTotal) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsSecurityTotal
}

// GetPolicyViolationsSecurityTotalOk returns a tuple with the PolicyViolationsSecurityTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsSecurityTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsSecurityTotal) {
		return nil, false
	}
	return o.PolicyViolationsSecurityTotal, true
}

// HasPolicyViolationsSecurityTotal returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsSecurityTotal() bool {
	if o != nil && !IsNil(o.PolicyViolationsSecurityTotal) {
		return true
	}

	return false
}

// SetPolicyViolationsSecurityTotal gets a reference to the given int32 and assigns it to the PolicyViolationsSecurityTotal field.
func (o *DependencyMetrics) SetPolicyViolationsSecurityTotal(v int32) {
	o.PolicyViolationsSecurityTotal = &v
}

// GetPolicyViolationsSecurityAudited returns the PolicyViolationsSecurityAudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsSecurityAudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsSecurityAudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsSecurityAudited
}

// GetPolicyViolationsSecurityAuditedOk returns a tuple with the PolicyViolationsSecurityAudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsSecurityAuditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsSecurityAudited) {
		return nil, false
	}
	return o.PolicyViolationsSecurityAudited, true
}

// HasPolicyViolationsSecurityAudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsSecurityAudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsSecurityAudited) {
		return true
	}

	return false
}

// SetPolicyViolationsSecurityAudited gets a reference to the given int32 and assigns it to the PolicyViolationsSecurityAudited field.
func (o *DependencyMetrics) SetPolicyViolationsSecurityAudited(v int32) {
	o.PolicyViolationsSecurityAudited = &v
}

// GetPolicyViolationsSecurityUnaudited returns the PolicyViolationsSecurityUnaudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsSecurityUnaudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsSecurityUnaudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsSecurityUnaudited
}

// GetPolicyViolationsSecurityUnauditedOk returns a tuple with the PolicyViolationsSecurityUnaudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsSecurityUnauditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsSecurityUnaudited) {
		return nil, false
	}
	return o.PolicyViolationsSecurityUnaudited, true
}

// HasPolicyViolationsSecurityUnaudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsSecurityUnaudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsSecurityUnaudited) {
		return true
	}

	return false
}

// SetPolicyViolationsSecurityUnaudited gets a reference to the given int32 and assigns it to the PolicyViolationsSecurityUnaudited field.
func (o *DependencyMetrics) SetPolicyViolationsSecurityUnaudited(v int32) {
	o.PolicyViolationsSecurityUnaudited = &v
}

// GetPolicyViolationsLicenseTotal returns the PolicyViolationsLicenseTotal field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsLicenseTotal() int32 {
	if o == nil || IsNil(o.PolicyViolationsLicenseTotal) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsLicenseTotal
}

// GetPolicyViolationsLicenseTotalOk returns a tuple with the PolicyViolationsLicenseTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsLicenseTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsLicenseTotal) {
		return nil, false
	}
	return o.PolicyViolationsLicenseTotal, true
}

// HasPolicyViolationsLicenseTotal returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsLicenseTotal() bool {
	if o != nil && !IsNil(o.PolicyViolationsLicenseTotal) {
		return true
	}

	return false
}

// SetPolicyViolationsLicenseTotal gets a reference to the given int32 and assigns it to the PolicyViolationsLicenseTotal field.
func (o *DependencyMetrics) SetPolicyViolationsLicenseTotal(v int32) {
	o.PolicyViolationsLicenseTotal = &v
}

// GetPolicyViolationsLicenseAudited returns the PolicyViolationsLicenseAudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsLicenseAudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsLicenseAudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsLicenseAudited
}

// GetPolicyViolationsLicenseAuditedOk returns a tuple with the PolicyViolationsLicenseAudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsLicenseAuditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsLicenseAudited) {
		return nil, false
	}
	return o.PolicyViolationsLicenseAudited, true
}

// HasPolicyViolationsLicenseAudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsLicenseAudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsLicenseAudited) {
		return true
	}

	return false
}

// SetPolicyViolationsLicenseAudited gets a reference to the given int32 and assigns it to the PolicyViolationsLicenseAudited field.
func (o *DependencyMetrics) SetPolicyViolationsLicenseAudited(v int32) {
	o.PolicyViolationsLicenseAudited = &v
}

// GetPolicyViolationsLicenseUnaudited returns the PolicyViolationsLicenseUnaudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsLicenseUnaudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsLicenseUnaudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsLicenseUnaudited
}

// GetPolicyViolationsLicenseUnauditedOk returns a tuple with the PolicyViolationsLicenseUnaudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsLicenseUnauditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsLicenseUnaudited) {
		return nil, false
	}
	return o.PolicyViolationsLicenseUnaudited, true
}

// HasPolicyViolationsLicenseUnaudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsLicenseUnaudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsLicenseUnaudited) {
		return true
	}

	return false
}

// SetPolicyViolationsLicenseUnaudited gets a reference to the given int32 and assigns it to the PolicyViolationsLicenseUnaudited field.
func (o *DependencyMetrics) SetPolicyViolationsLicenseUnaudited(v int32) {
	o.PolicyViolationsLicenseUnaudited = &v
}

// GetPolicyViolationsOperationalTotal returns the PolicyViolationsOperationalTotal field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsOperationalTotal() int32 {
	if o == nil || IsNil(o.PolicyViolationsOperationalTotal) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsOperationalTotal
}

// GetPolicyViolationsOperationalTotalOk returns a tuple with the PolicyViolationsOperationalTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsOperationalTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsOperationalTotal) {
		return nil, false
	}
	return o.PolicyViolationsOperationalTotal, true
}

// HasPolicyViolationsOperationalTotal returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsOperationalTotal() bool {
	if o != nil && !IsNil(o.PolicyViolationsOperationalTotal) {
		return true
	}

	return false
}

// SetPolicyViolationsOperationalTotal gets a reference to the given int32 and assigns it to the PolicyViolationsOperationalTotal field.
func (o *DependencyMetrics) SetPolicyViolationsOperationalTotal(v int32) {
	o.PolicyViolationsOperationalTotal = &v
}

// GetPolicyViolationsOperationalAudited returns the PolicyViolationsOperationalAudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsOperationalAudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsOperationalAudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsOperationalAudited
}

// GetPolicyViolationsOperationalAuditedOk returns a tuple with the PolicyViolationsOperationalAudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsOperationalAuditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsOperationalAudited) {
		return nil, false
	}
	return o.PolicyViolationsOperationalAudited, true
}

// HasPolicyViolationsOperationalAudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsOperationalAudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsOperationalAudited) {
		return true
	}

	return false
}

// SetPolicyViolationsOperationalAudited gets a reference to the given int32 and assigns it to the PolicyViolationsOperationalAudited field.
func (o *DependencyMetrics) SetPolicyViolationsOperationalAudited(v int32) {
	o.PolicyViolationsOperationalAudited = &v
}

// GetPolicyViolationsOperationalUnaudited returns the PolicyViolationsOperationalUnaudited field value if set, zero value otherwise.
func (o *DependencyMetrics) GetPolicyViolationsOperationalUnaudited() int32 {
	if o == nil || IsNil(o.PolicyViolationsOperationalUnaudited) {
		var ret int32
		return ret
	}
	return *o.PolicyViolationsOperationalUnaudited
}

// GetPolicyViolationsOperationalUnauditedOk returns a tuple with the PolicyViolationsOperationalUnaudited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DependencyMetrics) GetPolicyViolationsOperationalUnauditedOk() (*int32, bool) {
	if o == nil || IsNil(o.PolicyViolationsOperationalUnaudited) {
		return nil, false
	}
	return o.PolicyViolationsOperationalUnaudited, true
}

// HasPolicyViolationsOperationalUnaudited returns a boolean if a field has been set.
func (o *DependencyMetrics) HasPolicyViolationsOperationalUnaudited() bool {
	if o != nil && !IsNil(o.PolicyViolationsOperationalUnaudited) {
		return true
	}

	return false
}

// SetPolicyViolationsOperationalUnaudited gets a reference to the given int32 and assigns it to the PolicyViolationsOperationalUnaudited field.
func (o *DependencyMetrics) SetPolicyViolationsOperationalUnaudited(v int32) {
	o.PolicyViolationsOperationalUnaudited = &v
}

func (o DependencyMetrics) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DependencyMetrics) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Critical) {
		toSerialize["critical"] = o.Critical
	}
	if !IsNil(o.High) {
		toSerialize["high"] = o.High
	}
	if !IsNil(o.Medium) {
		toSerialize["medium"] = o.Medium
	}
	if !IsNil(o.Low) {
		toSerialize["low"] = o.Low
	}
	if !IsNil(o.Unassigned) {
		toSerialize["unassigned"] = o.Unassigned
	}
	if !IsNil(o.Vulnerabilities) {
		toSerialize["vulnerabilities"] = o.Vulnerabilities
	}
	if !IsNil(o.Suppressed) {
		toSerialize["suppressed"] = o.Suppressed
	}
	if !IsNil(o.InheritedRiskScore) {
		toSerialize["inherited_risk_score"] = o.InheritedRiskScore
	}
	if !IsNil(o.FindingsTotal) {
		toSerialize["findings_total"] = o.FindingsTotal
	}
	if !IsNil(o.FindingsAudited) {
		toSerialize["findings_audited"] = o.FindingsAudited
	}
	if !IsNil(o.FindingsUnaudited) {
		toSerialize["findings_unaudited"] = o.FindingsUnaudited
	}
	if !IsNil(o.PolicyViolationsFail) {
		toSerialize["policy_violations_fail"] = o.PolicyViolationsFail
	}
	if !IsNil(o.PolicyViolationsWarn) {
		toSerialize["policy_violations_warn"] = o.PolicyViolationsWarn
	}
	if !IsNil(o.PolicyViolationsInfo) {
		toSerialize["policy_violations_info"] = o.PolicyViolationsInfo
	}
	if !IsNil(o.PolicyViolationsTotal) {
		toSerialize["policy_violations_total"] = o.PolicyViolationsTotal
	}
	if !IsNil(o.PolicyViolationsAudited) {
		toSerialize["policy_violations_audited"] = o.PolicyViolationsAudited
	}
	if !IsNil(o.PolicyViolationsUnaudited) {
		toSerialize["policy_violations_unaudited"] = o.PolicyViolationsUnaudited
	}
	if !IsNil(o.PolicyViolationsSecurityTotal) {
		toSerialize["policy_violations_security_total"] = o.PolicyViolationsSecurityTotal
	}
	if !IsNil(o.PolicyViolationsSecurityAudited) {
		toSerialize["policy_violations_security_audited"] = o.PolicyViolationsSecurityAudited
	}
	if !IsNil(o.PolicyViolationsSecurityUnaudited) {
		toSerialize["policy_violations_security_unaudited"] = o.PolicyViolationsSecurityUnaudited
	}
	if !IsNil(o.PolicyViolationsLicenseTotal) {
		toSerialize["policy_violations_license_total"] = o.PolicyViolationsLicenseTotal
	}
	if !IsNil(o.PolicyViolationsLicenseAudited) {
		toSerialize["policy_violations_license_audited"] = o.PolicyViolationsLicenseAudited
	}
	if !IsNil(o.PolicyViolationsLicenseUnaudited) {
		toSerialize["policy_violations_license_unaudited"] = o.PolicyViolationsLicenseUnaudited
	}
	if !IsNil(o.PolicyViolationsOperationalTotal) {
		toSerialize["policy_violations_operational_total"] = o.PolicyViolationsOperationalTotal
	}
	if !IsNil(o.PolicyViolationsOperationalAudited) {
		toSerialize["policy_violations_operational_audited"] = o.PolicyViolationsOperationalAudited
	}
	if !IsNil(o.PolicyViolationsOperationalUnaudited) {
		toSerialize["policy_violations_operational_unaudited"] = o.PolicyViolationsOperationalUnaudited
	}
	return toSerialize, nil
}

type NullableDependencyMetrics struct {
	value *DependencyMetrics
	isSet bool
}

func (v NullableDependencyMetrics) Get() *DependencyMetrics {
	return v.value
}

func (v *NullableDependencyMetrics) Set(val *DependencyMetrics) {
	v.value = val
	v.isSet = true
}

func (v NullableDependencyMetrics) IsSet() bool {
	return v.isSet
}

func (v *NullableDependencyMetrics) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDependencyMetrics(val *DependencyMetrics) *NullableDependencyMetrics {
	return &NullableDependencyMetrics{value: val, isSet: true}
}

func (v NullableDependencyMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDependencyMetrics) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
