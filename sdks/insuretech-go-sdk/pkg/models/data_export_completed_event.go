package models

import (
	"time"
)

// DataExportCompletedEvent represents a data_export_completed_event
type DataExportCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	RecordCount string `json:"record_count,omitempty"`
	GenerationTimeSeconds int `json:"generation_time_seconds,omitempty"`
	ExportId string `json:"export_id,omitempty"`
	ExportUrl string `json:"export_url,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
