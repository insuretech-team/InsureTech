package models


// EmailLoginResponse represents a email_login_response
type EmailLoginResponse struct {
	CsrfToken string `json:"csrf_token,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	User *User `json:"user,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
