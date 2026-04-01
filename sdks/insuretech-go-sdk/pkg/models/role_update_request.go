package models


// RoleUpdateRequest represents a role_update_request
type RoleUpdateRequest struct {
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name"`
	RoleId string `json:"role_id"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
