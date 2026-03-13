package models

import (
	"time"
)

// Session represents a session
type Session struct {
	DeviceType *AuthnDeviceType `json:"device_type,omitempty"`
	SessionTokenLookup string `json:"session_token_lookup,omitempty"`
	AccessTokenJti string `json:"access_token_jti,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	LastActivityAt time.Time `json:"last_activity_at,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	RefreshTokenJti string `json:"refresh_token_jti,omitempty"`
	CsrfToken string `json:"csrf_token,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionTokenHash string `json:"session_token_hash,omitempty"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	SessionType *SessionType `json:"session_type,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
}
