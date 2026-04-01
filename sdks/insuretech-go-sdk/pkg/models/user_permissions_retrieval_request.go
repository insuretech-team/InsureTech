package models


// UserPermissionsRetrievalRequest represents a user_permissions_retrieval_request
type UserPermissionsRetrievalRequest struct {
	Domain string `json:"domain,omitempty"`
	UserId string `json:"user_id"`
}
