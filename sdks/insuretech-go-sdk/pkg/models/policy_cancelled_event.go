package models

import (
	"time"
)

// PolicyCancelledEvent represents a policy_cancelled_event
type PolicyCancelledEvent struct {
	PolicyNumber string `json:"policy_number,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
