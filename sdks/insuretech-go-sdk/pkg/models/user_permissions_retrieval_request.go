package models


// UserPermissionsRetrievalRequest represents a user_permissions_retrieval_request
type UserPermissionsRetrievalRequest struct {
	UserId string `json:"user_id"`
	Domain string `json:"domain,omitempty"`
}
