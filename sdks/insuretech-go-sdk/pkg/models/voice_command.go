package models

import (
	"time"
)

// VoiceCommand represents a voice_command
type VoiceCommand struct {
	VoiceSessionId string `json:"voice_session_id"`
	Type *CommandType `json:"type"`
	CommandText string `json:"command_text,omitempty"`
	Parameters string `json:"parameters,omitempty"`
	Status interface{} `json:"status"`
	ResponseAudioUrl string `json:"response_audio_url,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	ExecutedAt time.Time `json:"executed_at"`
	Id string `json:"id"`
	AudioUrl string `json:"audio_url,omitempty"`
	ResponseText string `json:"response_text,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
}
