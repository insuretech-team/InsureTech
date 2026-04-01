package models

import (
	"time"
)

// PolicyIssuedEvent represents a policy_issued_event
type PolicyIssuedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
}
