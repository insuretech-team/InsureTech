package models

import (
	"time"
)

// VoiceSessionEndedEvent represents a voice_session_ended_event
type VoiceSessionEndedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
