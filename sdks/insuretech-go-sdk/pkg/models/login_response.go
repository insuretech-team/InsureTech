package models


// LoginResponse represents a login_response
type LoginResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	MfaMethod string `json:"mfa_method,omitempty"`
	MfaRequired bool `json:"mfa_required,omitempty"`
	MfaSessionToken string `json:"mfa_session_token,omitempty"`
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	User *User `json:"user,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
