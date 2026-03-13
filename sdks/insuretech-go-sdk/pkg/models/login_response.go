package models


// LoginResponse represents a login_response
type LoginResponse struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	User *User `json:"user,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	MfaRequired bool `json:"mfa_required,omitempty"`
	MfaMethod string `json:"mfa_method,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
	MfaSessionToken string `json:"mfa_session_token,omitempty"`
	Error *Error `json:"error,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
