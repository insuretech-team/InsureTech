package models


// VoiceSessionStartResponse represents a voice_session_start_response
type VoiceSessionStartResponse struct {
	SessionId string `json:"session_id,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
