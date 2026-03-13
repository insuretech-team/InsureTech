package models

import (
	"time"
)

// VoiceSession represents a voice_session
type VoiceSession struct {
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Language string `json:"language,omitempty"`
	Status *SessionStatus `json:"status,omitempty"`
	Intent string `json:"intent,omitempty"`
	Context string `json:"context,omitempty"`
	EndedAt time.Time `json:"ended_at,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
}
