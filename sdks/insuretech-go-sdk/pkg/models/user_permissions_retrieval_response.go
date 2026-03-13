package models


// UserPermissionsRetrievalResponse represents a user_permissions_retrieval_response
type UserPermissionsRetrievalResponse struct {
	Permissions []*EffectivePermission `json:"permissions,omitempty"`
	Error *Error `json:"error,omitempty"`
}
