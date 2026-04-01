package models

import (
	"time"
)

// ReportExecutionCompletedEvent represents a report_execution_completed_event
type ReportExecutionCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileFormat string `json:"file_format,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	ReportExecutionId string `json:"report_execution_id,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
