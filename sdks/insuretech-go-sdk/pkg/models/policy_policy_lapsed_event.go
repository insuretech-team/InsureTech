package models

import (
	"time"
)

// PolicyPolicyLapsedEvent represents a policy_policy_lapsed_event
type PolicyPolicyLapsedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
