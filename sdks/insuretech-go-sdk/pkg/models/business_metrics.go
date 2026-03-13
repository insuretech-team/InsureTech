package models

import (
	"time"
)

// BusinessMetrics represents a business_metrics
type BusinessMetrics struct {
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	Type *MetricType `json:"type"`
	Value float64 `json:"value"`
	Dimensions map[string]interface{} `json:"dimensions,omitempty"`
	RecordedAt time.Time `json:"recorded_at"`
}
