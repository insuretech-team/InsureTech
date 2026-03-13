package models

import (
	"time"
)

// DataExportRequestedEvent represents a data_export_requested_event
type DataExportRequestedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExportId string `json:"export_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	ExportType string `json:"export_type,omitempty"`
	Format string `json:"format,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
