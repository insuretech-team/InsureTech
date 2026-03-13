package models

import (
	"time"
)

// RefundRefundProcessedEvent represents a refund_refund_processed_event
type RefundRefundProcessedEvent struct {
	PaymentReference string `json:"payment_reference,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
}
