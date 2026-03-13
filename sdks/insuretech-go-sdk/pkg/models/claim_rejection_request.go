package models


// ClaimRejectionRequest represents a claim_rejection_request
type ClaimRejectionRequest struct {
	ClaimId string `json:"claim_id"`
	ApproverId string `json:"approver_id"`
	Reason string `json:"reason,omitempty"`
}
