package models

import (
	"time"
)

// MediaOCRCompletedEvent represents a media_ocrcompleted_event
type MediaOCRCompletedEvent struct {
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	OcrText string `json:"ocr_text,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
