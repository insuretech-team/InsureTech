package models


// InitiateVoiceSessionResponse represents a initiate_voice_session_response
type InitiateVoiceSessionResponse struct {
	Challenge string `json:"challenge,omitempty"`
	SessionId string `json:"session_id,omitempty"`
}
