package models


// RevokeRoleResponse represents a revoke_role_response
type RevokeRoleResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
