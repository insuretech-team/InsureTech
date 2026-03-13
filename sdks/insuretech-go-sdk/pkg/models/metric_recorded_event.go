package models

import (
	"time"
)

// MetricRecordedEvent represents a metric_recorded_event
type MetricRecordedEvent struct {
	MetricId string `json:"metric_id,omitempty"`
	MetricName string `json:"metric_name,omitempty"`
	MetricType string `json:"metric_type,omitempty"`
	Value float64 `json:"value,omitempty"`
	Dimensions map[string]interface{} `json:"dimensions,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
