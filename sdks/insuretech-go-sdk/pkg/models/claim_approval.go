package models

import (
	"time"
)

// ClaimApproval represents a claim_approval
type ClaimApproval struct {
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	ApprovalId string `json:"approval_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	Decision *ApprovalDecision `json:"decision,omitempty"`
	Notes string `json:"notes,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ApprovedCurrency string `json:"approved_currency,omitempty"`
	ApproverId string `json:"approver_id,omitempty"`
	ApproverRole string `json:"approver_role,omitempty"`
	ApprovalLevel int `json:"approval_level,omitempty"`
}
