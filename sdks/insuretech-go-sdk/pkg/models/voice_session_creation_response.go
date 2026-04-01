package models


// VoiceSessionCreationResponse represents a voice_session_creation_response
type VoiceSessionCreationResponse struct {
	Status string `json:"status,omitempty"`
	VoiceSessionId string `json:"voice_session_id,omitempty"`
}
