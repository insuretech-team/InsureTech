package models


// AuthnJWKSRetrievalResponse represents a authn_jwksretrieval_response
type AuthnJWKSRetrievalResponse struct {
	Keys []*AuthnJWK `json:"keys,omitempty"`
}
