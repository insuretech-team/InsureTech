package models

import (
	"time"
)

// CommissionPayoutCreatedEvent represents a commission_payout_created_event
type CommissionPayoutCreatedEvent struct {
	CommissionCount int `json:"commission_count,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PayoutId string `json:"payout_id,omitempty"`
	PayoutNumber string `json:"payout_number,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
}
