package models

import (
	"time"
)

// ClaimApprovedEvent represents a claim_approved_event
type ClaimApprovedEvent struct {
	ApprovalLevel string `json:"approval_level,omitempty"`
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	ApproverId string `json:"approver_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
