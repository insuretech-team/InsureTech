package models

import (
	"time"
)

// ReportExecutionStartedEvent represents a report_execution_started_event
type ReportExecutionStartedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExecutedBy string `json:"executed_by,omitempty"`
	ReportDefinitionId string `json:"report_definition_id,omitempty"`
	ReportExecutionId string `json:"report_execution_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
