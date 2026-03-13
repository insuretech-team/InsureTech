package models

import (
	"time"
)

// AggregationJobCompletedEvent represents a aggregation_job_completed_event
type AggregationJobCompletedEvent struct {
	ExecutionTimeSeconds int `json:"execution_time_seconds,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	JobId string `json:"job_id,omitempty"`
	RecordsProcessed string `json:"records_processed,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	JobType string `json:"job_type,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
}
