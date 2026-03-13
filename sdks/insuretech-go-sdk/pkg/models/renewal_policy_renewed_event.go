package models

import (
	"time"
)

// RenewalPolicyRenewedEvent represents a renewal_policy_renewed_event
type RenewalPolicyRenewedEvent struct {
	OldPolicyId string `json:"old_policy_id,omitempty"`
	NewPolicyId string `json:"new_policy_id,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
