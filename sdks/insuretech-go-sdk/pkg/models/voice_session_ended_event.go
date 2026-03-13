package models

import (
	"time"
)

// VoiceSessionEndedEvent represents a voice_session_ended_event
type VoiceSessionEndedEvent struct {
	Status string `json:"status,omitempty"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
