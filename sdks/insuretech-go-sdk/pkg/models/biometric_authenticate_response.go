package models


// BiometricAuthenticateResponse represents a biometric_authenticate_response
type BiometricAuthenticateResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	User *User `json:"user,omitempty"`
	Error *Error `json:"error,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
}
