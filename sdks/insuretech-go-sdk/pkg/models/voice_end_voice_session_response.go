package models


// VoiceEndVoiceSessionResponse represents a voice_end_voice_session_response
type VoiceEndVoiceSessionResponse struct {
	DurationSeconds int `json:"duration_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
