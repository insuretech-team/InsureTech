package models

import (
	"time"
)

// MetricAnomalyDetectedEvent represents a metric_anomaly_detected_event
type MetricAnomalyDetectedEvent struct {
	AnomalyType string `json:"anomaly_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MetricName string `json:"metric_name,omitempty"`
	CurrentValue float64 `json:"current_value,omitempty"`
	ExpectedValue float64 `json:"expected_value,omitempty"`
	DeviationPercentage float64 `json:"deviation_percentage,omitempty"`
}
