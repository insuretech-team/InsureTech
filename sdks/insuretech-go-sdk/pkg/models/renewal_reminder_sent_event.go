package models

import (
	"time"
)

// RenewalReminderSentEvent represents a renewal_reminder_sent_event
type RenewalReminderSentEvent struct {
	Channel string `json:"channel,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	DaysBeforeRenewal int `json:"days_before_renewal,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	ReminderId string `json:"reminder_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
