package models


// AuthnEndVoiceSessionResponse represents a authn_end_voice_session_response
type AuthnEndVoiceSessionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
