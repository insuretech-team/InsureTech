package models


// JWKSRetrievalResponse represents a jwksretrieval_response
type JWKSRetrievalResponse struct {
	Keys []*JWK `json:"keys,omitempty"`
	Error *Error `json:"error,omitempty"`
}
