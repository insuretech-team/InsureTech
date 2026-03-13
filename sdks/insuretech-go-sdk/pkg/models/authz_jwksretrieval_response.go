package models


// AuthzJWKSRetrievalResponse represents a authz_jwksretrieval_response
type AuthzJWKSRetrievalResponse struct {
	Keys []*AuthzJWK `json:"keys,omitempty"`
	Error *Error `json:"error,omitempty"`
}
