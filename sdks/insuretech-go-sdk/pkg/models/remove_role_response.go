package models


// RemoveRoleResponse represents a remove_role_response
type RemoveRoleResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
