package models


// LogoutRequest represents a logout_request
type LogoutRequest struct {
	AccessToken string `json:"access_token,omitempty"`
	LogoutReason string `json:"logout_reason,omitempty"`
	SessionId string `json:"session_id"`
}
