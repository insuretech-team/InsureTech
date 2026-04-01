package models

import (
	"time"
)

// CommissionPaidEvent represents a commission_paid_event
type CommissionPaidEvent struct {
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PayoutId string `json:"payout_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
