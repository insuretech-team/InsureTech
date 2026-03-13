package models

import (
	"time"
)

// MediaProcessingFailedEvent represents a media_processing_failed_event
type MediaProcessingFailedEvent struct {
	RetryCount int `json:"retry_count,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	JobId string `json:"job_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	ProcessingType string `json:"processing_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
