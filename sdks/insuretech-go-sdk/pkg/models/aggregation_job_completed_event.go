package models

import (
	"time"
)

// AggregationJobCompletedEvent represents a aggregation_job_completed_event
type AggregationJobCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExecutionTimeSeconds int `json:"execution_time_seconds,omitempty"`
	JobId string `json:"job_id,omitempty"`
	JobType string `json:"job_type,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	RecordsProcessed string `json:"records_processed,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
