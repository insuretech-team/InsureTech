package models

import (
	"time"
)

// MediaProcessingJobCreatedEvent represents a media_processing_job_created_event
type MediaProcessingJobCreatedEvent struct {
	Priority int `json:"priority,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	JobId string `json:"job_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	ProcessingType string `json:"processing_type,omitempty"`
}
