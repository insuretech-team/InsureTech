package models


// ClaimRetrievalResponse represents a claim_retrieval_response
type ClaimRetrievalResponse struct {
	Claim *Claim `json:"claim,omitempty"`
	Error *Error `json:"error,omitempty"`
}
