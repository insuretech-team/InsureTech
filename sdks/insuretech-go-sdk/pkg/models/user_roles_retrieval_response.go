package models


// UserRolesRetrievalResponse represents a user_roles_retrieval_response
type UserRolesRetrievalResponse struct {
	Roles []*Role `json:"roles,omitempty"`
	Error *Error `json:"error,omitempty"`
}
