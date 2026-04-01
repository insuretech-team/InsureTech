package models

import (
	"time"
)

// PolicyPolicyLapsedEvent represents a policy_policy_lapsed_event
type PolicyPolicyLapsedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
