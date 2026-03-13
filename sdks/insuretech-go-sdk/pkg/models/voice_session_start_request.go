package models


// VoiceSessionStartRequest represents a voice_session_start_request
type VoiceSessionStartRequest struct {
	UserId string `json:"user_id"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Language string `json:"language,omitempty"`
}
