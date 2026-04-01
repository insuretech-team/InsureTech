package models

import (
	"time"
)

// ClaimRejectedEvent represents a claim_rejected_event
type ClaimRejectedEvent struct {
	ApproverId string `json:"approver_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
