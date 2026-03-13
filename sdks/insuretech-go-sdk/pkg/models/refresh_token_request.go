package models


// RefreshTokenRequest represents a refresh_token_request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	DeviceId string `json:"device_id"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
