package models

import (
	"time"
)

// MediaOCRCompletedEvent represents a media_ocrcompleted_event
type MediaOCRCompletedEvent struct {
	MediaId string `json:"media_id,omitempty"`
	OcrText string `json:"ocr_text,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
