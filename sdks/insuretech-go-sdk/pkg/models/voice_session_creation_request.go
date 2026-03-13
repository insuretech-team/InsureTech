package models


// VoiceSessionCreationRequest represents a voice_session_creation_request
type VoiceSessionCreationRequest struct {
	UserId string `json:"user_id"`
	Language string `json:"language,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
