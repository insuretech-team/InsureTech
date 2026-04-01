package models

import (
	"time"
)

// RefundRequestedEvent represents a refund_requested_event
type RefundRequestedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	RefundNumber string `json:"refund_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
