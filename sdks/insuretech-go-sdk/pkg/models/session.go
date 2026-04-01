package models

import (
	"time"
)

// Session represents a session
type Session struct {
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
	AccessTokenJti string `json:"access_token_jti,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType *DeviceType `json:"device_type,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	LastActivityAt time.Time `json:"last_activity_at,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
	RefreshTokenJti string `json:"refresh_token_jti,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionTokenHash string `json:"session_token_hash,omitempty"`
	SessionTokenLookup string `json:"session_token_lookup,omitempty"`
	SessionType *SessionType `json:"session_type,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
