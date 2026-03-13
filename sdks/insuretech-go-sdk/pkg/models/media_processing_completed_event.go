package models

import (
	"time"
)

// MediaProcessingCompletedEvent represents a media_processing_completed_event
type MediaProcessingCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	JobId string `json:"job_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	ProcessingType string `json:"processing_type,omitempty"`
	ResultFileId string `json:"result_file_id,omitempty"`
}
