package models

import (
	"time"
)

// TokenRefreshedEvent represents a token_refreshed_event
type TokenRefreshedEvent struct {
	DeviceType string `json:"device_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	NewAccessTokenJti string `json:"new_access_token_jti,omitempty"`
	NewRefreshTokenJti string `json:"new_refresh_token_jti,omitempty"`
	OldAccessTokenJti string `json:"old_access_token_jti,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
