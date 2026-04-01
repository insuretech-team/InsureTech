package models

import (
	"time"
)

// RefundRefundProcessedEvent represents a refund_refund_processed_event
type RefundRefundProcessedEvent struct {
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
