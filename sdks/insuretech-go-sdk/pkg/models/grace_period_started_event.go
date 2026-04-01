package models

import (
	"time"
)

// GracePeriodStartedEvent represents a grace_period_started_event
type GracePeriodStartedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	EventId string `json:"event_id,omitempty"`
	GracePeriodId string `json:"grace_period_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
