package models


// InsuranceClaimRetrievalResponse represents a insurance_claim_retrieval_response
type InsuranceClaimRetrievalResponse struct {
	Claim *Claim `json:"claim,omitempty"`
}
