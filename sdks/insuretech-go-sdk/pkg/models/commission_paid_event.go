package models

import (
	"time"
)

// CommissionPaidEvent represents a commission_paid_event
type CommissionPaidEvent struct {
	EventId string `json:"event_id,omitempty"`
	PayoutId string `json:"payout_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
