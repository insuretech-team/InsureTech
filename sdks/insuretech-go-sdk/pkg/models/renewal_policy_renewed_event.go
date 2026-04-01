package models

import (
	"time"
)

// RenewalPolicyRenewedEvent represents a renewal_policy_renewed_event
type RenewalPolicyRenewedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewPolicyId string `json:"new_policy_id,omitempty"`
	OldPolicyId string `json:"old_policy_id,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
