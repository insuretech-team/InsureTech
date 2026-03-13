package models


// AuthnEndVoiceSessionRequest represents a authn_end_voice_session_request
type AuthnEndVoiceSessionRequest struct {
	Status string `json:"status,omitempty"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
	VoiceSessionId string `json:"voice_session_id"`
}
