package models


// ClaimApprovalRequest represents a claim_approval_request
type ClaimApprovalRequest struct {
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	ApproverId string `json:"approver_id"`
	ClaimId string `json:"claim_id"`
	Notes string `json:"notes,omitempty"`
}
