package models


// ClaimRejectionRequest represents a claim_rejection_request
type ClaimRejectionRequest struct {
	ApproverId string `json:"approver_id"`
	ClaimId string `json:"claim_id"`
	Reason string `json:"reason,omitempty"`
}
