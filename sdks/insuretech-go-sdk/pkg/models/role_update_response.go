package models


// RoleUpdateResponse represents a role_update_response
type RoleUpdateResponse struct {
	Role *Role `json:"role,omitempty"`
	Error *Error `json:"error,omitempty"`
}
