package models


// VoiceSessionVerificationResponse represents a voice_session_verification_response
type VoiceSessionVerificationResponse struct {
	UserId string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Authenticated bool `json:"authenticated,omitempty"`
}
