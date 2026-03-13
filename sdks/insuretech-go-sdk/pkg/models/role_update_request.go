package models


// RoleUpdateRequest represents a role_update_request
type RoleUpdateRequest struct {
	RoleId string `json:"role_id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
