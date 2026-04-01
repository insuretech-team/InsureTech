package models

import (
	"time"
)

// VoiceSession represents a voice_session
type VoiceSession struct {
	AuditInfo interface{} `json:"audit_info"`
	Context string `json:"context,omitempty"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
	EndedAt time.Time `json:"ended_at,omitempty"`
	Id string `json:"id,omitempty"`
	Intent string `json:"intent,omitempty"`
	Language string `json:"language,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Status *SessionStatus `json:"status,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
