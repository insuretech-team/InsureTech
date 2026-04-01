package models

import (
	"time"
)

// KPIThresholdBreachedEvent represents a kpithreshold_breached_event
type KPIThresholdBreachedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	CurrentValue float64 `json:"current_value,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MetricName string `json:"metric_name,omitempty"`
	NotificationSentTo string `json:"notification_sent_to,omitempty"`
	Severity string `json:"severity,omitempty"`
	ThresholdType string `json:"threshold_type,omitempty"`
	ThresholdValue float64 `json:"threshold_value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
