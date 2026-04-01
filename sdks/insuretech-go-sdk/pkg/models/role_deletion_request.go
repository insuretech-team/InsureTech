package models


// RoleDeletionRequest represents a role_deletion_request
type RoleDeletionRequest struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	RoleId string `json:"role_id"`
}
