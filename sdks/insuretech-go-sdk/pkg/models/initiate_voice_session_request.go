package models


// InitiateVoiceSessionRequest represents a initiate_voice_session_request
type InitiateVoiceSessionRequest struct {
	UserId string `json:"user_id"`
	PortalId string `json:"portal_id"`
}
