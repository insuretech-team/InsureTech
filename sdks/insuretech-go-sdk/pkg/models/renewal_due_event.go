package models

import (
	"time"
)

// RenewalDueEvent represents a renewal_due_event
type RenewalDueEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	RenewalDueDate string `json:"renewal_due_date,omitempty"`
	RenewalPremium *Money `json:"renewal_premium,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
