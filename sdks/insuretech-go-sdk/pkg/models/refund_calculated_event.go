package models

import (
	"time"
)

// RefundCalculatedEvent represents a refund_calculated_event
type RefundCalculatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	RefundableAmount *Money `json:"refundable_amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
