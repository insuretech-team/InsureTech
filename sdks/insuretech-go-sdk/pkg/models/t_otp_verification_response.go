package models


// TOTPVerificationResponse represents a t_otp_verification_response
type TOTPVerificationResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Verified bool `json:"verified,omitempty"`
}
