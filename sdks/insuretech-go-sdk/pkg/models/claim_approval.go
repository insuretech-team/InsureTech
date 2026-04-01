package models

import (
	"time"
)

// ClaimApproval represents a claim_approval
type ClaimApproval struct {
	ApprovalId string `json:"approval_id,omitempty"`
	ApprovalLevel int `json:"approval_level,omitempty"`
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedCurrency string `json:"approved_currency,omitempty"`
	ApproverId string `json:"approver_id,omitempty"`
	ApproverRole string `json:"approver_role,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Decision *ApprovalDecision `json:"decision,omitempty"`
	Notes string `json:"notes,omitempty"`
}
