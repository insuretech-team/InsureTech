package models

import (
	"time"
)

// RenewalPolicyLapsedEvent represents a renewal_policy_lapsed_event
type RenewalPolicyLapsedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	GracePeriodId string `json:"grace_period_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
