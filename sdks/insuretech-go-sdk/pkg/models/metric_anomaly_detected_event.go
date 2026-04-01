package models

import (
	"time"
)

// MetricAnomalyDetectedEvent represents a metric_anomaly_detected_event
type MetricAnomalyDetectedEvent struct {
	AnomalyType string `json:"anomaly_type,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CurrentValue float64 `json:"current_value,omitempty"`
	DeviationPercentage float64 `json:"deviation_percentage,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExpectedValue float64 `json:"expected_value,omitempty"`
	MetricName string `json:"metric_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
