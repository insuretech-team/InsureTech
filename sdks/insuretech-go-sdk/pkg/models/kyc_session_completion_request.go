package models


// KYCSessionCompletionRequest represents a kyc_session_completion_request
type KYCSessionCompletionRequest struct {
	UserId string `json:"user_id"`
	SessionId string `json:"session_id"`
}
