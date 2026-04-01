package models

import (
	"time"
)

// AggregatedMetric represents a aggregated_metric
type AggregatedMetric struct {
	AggregatedMetricId string `json:"aggregated_metric_id"`
	Aggregation *MetricAggregation `json:"aggregation"`
	Dimensions map[string]interface{} `json:"dimensions,omitempty"`
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	TimeBucket string `json:"time_bucket"`
	Timestamp time.Time `json:"timestamp"`
	Value float64 `json:"value"`
}
