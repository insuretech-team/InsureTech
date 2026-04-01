package models


// VoiceSessionStartRequest represents a voice_session_start_request
type VoiceSessionStartRequest struct {
	Language string `json:"language,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	UserId string `json:"user_id"`
}
