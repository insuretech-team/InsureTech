package models

import (
	"time"
)

// ReportGeneratedEvent represents a report_generated_event
type ReportGeneratedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ReportId string `json:"report_id,omitempty"`
	ReportType string `json:"report_type,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	ReportUrl string `json:"report_url,omitempty"`
	GenerationTimeSeconds int `json:"generation_time_seconds,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	ReportName string `json:"report_name,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	GeneratedBy string `json:"generated_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
