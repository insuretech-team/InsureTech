package models

import (
	"time"
)

// ReportGenerationFailedEvent represents a report_generation_failed_event
type ReportGenerationFailedEvent struct {
	ReportId string `json:"report_id,omitempty"`
	ReportName string `json:"report_name,omitempty"`
	ReportType string `json:"report_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
