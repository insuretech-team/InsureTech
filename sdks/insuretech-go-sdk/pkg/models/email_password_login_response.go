package models


// EmailPasswordLoginResponse represents a email_password_login_response
type EmailPasswordLoginResponse struct {
	CsrfToken string `json:"csrf_token,omitempty"`
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	User *User `json:"user,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
