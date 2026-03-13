package models


// VoiceSessionCreationResponse represents a voice_session_creation_response
type VoiceSessionCreationResponse struct {
	VoiceSessionId string `json:"voice_session_id,omitempty"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
