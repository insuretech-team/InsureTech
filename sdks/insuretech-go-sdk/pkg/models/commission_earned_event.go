package models

import (
	"time"
)

// CommissionEarnedEvent represents a commission_earned_event
type CommissionEarnedEvent struct {
	Amount *Money `json:"amount,omitempty"`
	CommissionId string `json:"commission_id,omitempty"`
	CommissionNumber string `json:"commission_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	RecipientType string `json:"recipient_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
