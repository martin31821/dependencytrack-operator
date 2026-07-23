/*
OWASP Dependency-Track

Extensions to generated DT API models for operator-specific fields.

These fields are supported by the Dependency-Track API but not yet
present in the OpenAPI-generated model structs.
*/

package dtapi

// UpdateNotificationRuleRequestExtensions adds Teams and Projects fields
// to the generated UpdateNotificationRuleRequest struct.
type UpdateNotificationRuleRequestExtensions struct {
	Teams    []Team    `json:"teams,omitempty"`
	Projects []Project `json:"projects,omitempty"`
}
