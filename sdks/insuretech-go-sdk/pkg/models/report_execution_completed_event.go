package models

import (
	"time"
)

// ReportExecutionCompletedEvent represents a report_execution_completed_event
type ReportExecutionCompletedEvent struct {
	FileUrl string `json:"file_url,omitempty"`
	FileFormat string `json:"file_format,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ReportExecutionId string `json:"report_execution_id,omitempty"`
}
