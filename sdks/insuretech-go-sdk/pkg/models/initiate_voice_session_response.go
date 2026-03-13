package models


// InitiateVoiceSessionResponse represents a initiate_voice_session_response
type InitiateVoiceSessionResponse struct {
	SessionId string `json:"session_id,omitempty"`
	Challenge string `json:"challenge,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
