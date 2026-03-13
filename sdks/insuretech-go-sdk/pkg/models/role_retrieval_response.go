package models


// RoleRetrievalResponse represents a role_retrieval_response
type RoleRetrievalResponse struct {
	Role *Role `json:"role,omitempty"`
	Error *Error `json:"error,omitempty"`
}
