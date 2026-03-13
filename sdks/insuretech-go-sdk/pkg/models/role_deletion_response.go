package models


// RoleDeletionResponse represents a role_deletion_response
type RoleDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
