package models

import (
	"time"
)

// RenewalDueEvent represents a renewal_due_event
type RenewalDueEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	RenewalDueDate string `json:"renewal_due_date,omitempty"`
	RenewalPremium *Money `json:"renewal_premium,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
