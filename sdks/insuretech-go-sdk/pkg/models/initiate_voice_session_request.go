package models


// InitiateVoiceSessionRequest represents a initiate_voice_session_request
type InitiateVoiceSessionRequest struct {
	PortalId string `json:"portal_id"`
	UserId string `json:"user_id"`
}
