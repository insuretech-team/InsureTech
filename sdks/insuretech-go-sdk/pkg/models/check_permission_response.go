package models


// CheckPermissionResponse represents a check_permission_response
type CheckPermissionResponse struct {
	AppliedPolicies []string `json:"applied_policies,omitempty"`
	Error *Error `json:"error,omitempty"`
	Allowed bool `json:"allowed,omitempty"`
	Reason string `json:"reason,omitempty"`
}
