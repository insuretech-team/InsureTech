package models


// RoleAssignmentResponse represents a role_assignment_response
type RoleAssignmentResponse struct {
	UserRole *UserRole `json:"user_role,omitempty"`
	Error *Error `json:"error,omitempty"`
}
