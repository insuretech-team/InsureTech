package models

import (
	"time"
)

// AggregatedMetric represents a aggregated_metric
type AggregatedMetric struct {
	Dimensions map[string]interface{} `json:"dimensions,omitempty"`
	TimeBucket string `json:"time_bucket"`
	AggregatedMetricId string `json:"aggregated_metric_id"`
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	Aggregation *MetricAggregation `json:"aggregation"`
	Value float64 `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
