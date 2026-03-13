package models


// RoleCreationRequest represents a role_creation_request
type RoleCreationRequest struct {
	Portal *Portal `json:"portal,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Name string `json:"name"`
}
