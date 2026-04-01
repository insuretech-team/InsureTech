package models

import (
	"time"
)

// VoiceCommandExecutedEvent represents a voice_command_executed_event
type VoiceCommandExecutedEvent struct {
	CommandType string `json:"command_type,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	VoiceCommandId string `json:"voice_command_id,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
