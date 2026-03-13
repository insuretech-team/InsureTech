package models

import (
	"time"
)

// PolicyRevivedEvent represents a policy_revived_event
type PolicyRevivedEvent struct {
	GracePeriodId string `json:"grace_period_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
