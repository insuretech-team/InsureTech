package models

import (
	"time"
)

// RefundApprovedEvent represents a refund_approved_event
type RefundApprovedEvent struct {
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
}
