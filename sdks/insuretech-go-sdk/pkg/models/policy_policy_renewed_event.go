package models

import (
	"time"
)

// PolicyPolicyRenewedEvent represents a policy_policy_renewed_event
type PolicyPolicyRenewedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OldPolicyId string `json:"old_policy_id,omitempty"`
	NewPolicyId string `json:"new_policy_id,omitempty"`
	NewPolicyNumber string `json:"new_policy_number,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
