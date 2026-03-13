package models


// TOTPVerificationResponse represents a t_otp_verification_response
type TOTPVerificationResponse struct {
	SessionType string `json:"session_type,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
	Verified bool `json:"verified,omitempty"`
	Error *Error `json:"error,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	Message string `json:"message,omitempty"`
}
