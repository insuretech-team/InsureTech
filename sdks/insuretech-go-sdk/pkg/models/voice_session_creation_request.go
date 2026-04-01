package models


// VoiceSessionCreationRequest represents a voice_session_creation_request
type VoiceSessionCreationRequest struct {
	Language string `json:"language,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	UserId string `json:"user_id"`
}
