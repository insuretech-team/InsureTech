package models

import (
	"time"
)

// BusinessMetrics represents a business_metrics
type BusinessMetrics struct {
	Dimensions map[string]interface{} `json:"dimensions,omitempty"`
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	RecordedAt time.Time `json:"recorded_at"`
	Type *MetricType `json:"type"`
	Value float64 `json:"value"`
}
