package models


// KYCSessionCompletionRequest represents a kyc_session_completion_request
type KYCSessionCompletionRequest struct {
	SessionId string `json:"session_id"`
	UserId string `json:"user_id"`
}
