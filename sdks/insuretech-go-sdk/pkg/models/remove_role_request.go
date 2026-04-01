package models


// RemoveRoleRequest represents a remove_role_request
type RemoveRoleRequest struct {
	Domain string `json:"domain,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	RoleId string `json:"role_id"`
	UserId string `json:"user_id"`
}
