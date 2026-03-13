package models


// UserRolesRetrievalRequest represents a user_roles_retrieval_request
type UserRolesRetrievalRequest struct {
	UserId string `json:"user_id"`
	TenantId string `json:"tenant_id"`
}
