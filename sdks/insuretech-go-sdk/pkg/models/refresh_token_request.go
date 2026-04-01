package models


// RefreshTokenRequest represents a refresh_token_request
type RefreshTokenRequest struct {
	DeviceId string `json:"device_id"`
	IpAddress string `json:"ip_address,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
