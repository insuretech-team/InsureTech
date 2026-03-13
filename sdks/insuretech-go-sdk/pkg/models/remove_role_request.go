package models


// RemoveRoleRequest represents a remove_role_request
type RemoveRoleRequest struct {
	RoleId string `json:"role_id"`
	Domain string `json:"domain,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	UserId string `json:"user_id"`
}
