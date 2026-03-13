package models

import (
	"time"
)

// ReportExecutionStartedEvent represents a report_execution_started_event
type ReportExecutionStartedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ReportExecutionId string `json:"report_execution_id,omitempty"`
	ReportDefinitionId string `json:"report_definition_id,omitempty"`
	ExecutedBy string `json:"executed_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
