package models


// LogoutRequest represents a logout_request
type LogoutRequest struct {
	SessionId string `json:"session_id"`
	AccessToken string `json:"access_token,omitempty"`
	LogoutReason string `json:"logout_reason,omitempty"`
}
