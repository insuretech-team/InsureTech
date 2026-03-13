package models


// EmailLoginResponse represents a email_login_response
type EmailLoginResponse struct {
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	User *User `json:"user,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Error *Error `json:"error,omitempty"`
}
