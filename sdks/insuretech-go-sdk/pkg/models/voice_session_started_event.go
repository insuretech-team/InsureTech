package models

import (
	"time"
)

// VoiceSessionStartedEvent represents a voice_session_started_event
type VoiceSessionStartedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Language string `json:"language,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
