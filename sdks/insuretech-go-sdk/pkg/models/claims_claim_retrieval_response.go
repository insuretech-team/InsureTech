package models


// ClaimsClaimRetrievalResponse represents a claims_claim_retrieval_response
type ClaimsClaimRetrievalResponse struct {
	Claim *Claim `json:"claim,omitempty"`
	Error *Error `json:"error,omitempty"`
}
