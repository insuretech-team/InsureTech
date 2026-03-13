package models

import (
	"time"
)

// AuthnVoiceSessionRetrievalResponse represents a authn_voice_session_retrieval_response
type AuthnVoiceSessionRetrievalResponse struct {
	EndedAt time.Time `json:"ended_at,omitempty"`
	Error *Error `json:"error,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
	Language string `json:"language,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
}
