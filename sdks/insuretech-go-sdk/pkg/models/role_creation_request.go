package models


// RoleCreationRequest represents a role_creation_request
type RoleCreationRequest struct {
	CreatedBy string `json:"created_by,omitempty"`
	Description string `json:"description,omitempty"`
	Name string `json:"name"`
	Portal *Portal `json:"portal,omitempty"`
}
