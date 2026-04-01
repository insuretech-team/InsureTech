package models

import (
	"time"
)

// ClaimSubmittedEvent represents a claim_submitted_event
type ClaimSubmittedEvent struct {
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	ClaimType string `json:"claim_type,omitempty"`
	ClaimedAmount *Money `json:"claimed_amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
