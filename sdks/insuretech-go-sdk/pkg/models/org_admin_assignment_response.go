package models


// OrgAdminAssignmentResponse represents a org_admin_assignment_response
type OrgAdminAssignmentResponse struct {
	Member *OrgMember `json:"member,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
