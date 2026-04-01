package models

import (
	"time"
)

// ReportExecutionFailedEvent represents a report_execution_failed_event
type ReportExecutionFailedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ReportExecutionId string `json:"report_execution_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
