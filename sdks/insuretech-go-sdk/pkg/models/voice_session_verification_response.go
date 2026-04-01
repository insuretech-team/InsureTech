package models


// VoiceSessionVerificationResponse represents a voice_session_verification_response
type VoiceSessionVerificationResponse struct {
	Authenticated bool `json:"authenticated,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
