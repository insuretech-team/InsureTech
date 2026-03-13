package models

import (
	"time"
)

// VoiceCommandExecutedEvent represents a voice_command_executed_event
type VoiceCommandExecutedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	VoiceCommandId string `json:"voice_command_id,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
	CommandType string `json:"command_type,omitempty"`
	Status string `json:"status,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
