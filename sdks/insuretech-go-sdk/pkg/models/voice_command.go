package models

import (
	"time"
)

// VoiceCommand represents a voice_command
type VoiceCommand struct {
	AudioUrl string `json:"audio_url,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	CommandText string `json:"command_text,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	ExecutedAt time.Time `json:"executed_at"`
	Id string `json:"id"`
	Parameters string `json:"parameters,omitempty"`
	ResponseAudioUrl string `json:"response_audio_url,omitempty"`
	ResponseText string `json:"response_text,omitempty"`
	Status interface{} `json:"status"`
	Type *CommandType `json:"type"`
	VoiceSessionId string `json:"voice_session_id"`
}
