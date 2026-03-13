package models


// RoleDeletionRequest represents a role_deletion_request
type RoleDeletionRequest struct {
	RoleId string `json:"role_id"`
	DeletedBy string `json:"deleted_by,omitempty"`
}
