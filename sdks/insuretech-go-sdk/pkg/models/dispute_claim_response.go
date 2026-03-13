package models


// DisputeClaimResponse represents a dispute_claim_response
type DisputeClaimResponse struct {
	Error *Error `json:"error,omitempty"`
	DisputeId string `json:"dispute_id,omitempty"`
	Message string `json:"message,omitempty"`
}
