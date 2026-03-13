package models


// RoleCreationResponse represents a role_creation_response
type RoleCreationResponse struct {
	Role *Role `json:"role,omitempty"`
	Error *Error `json:"error,omitempty"`
}
