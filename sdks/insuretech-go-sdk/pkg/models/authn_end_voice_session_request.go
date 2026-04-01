package models


// AuthnEndVoiceSessionRequest represents a authn_end_voice_session_request
type AuthnEndVoiceSessionRequest struct {
	VoiceSessionId string `json:"voice_session_id"`
	Status string `json:"status"`
	DurationSeconds int `json:"duration_seconds,omitempty"`
}
