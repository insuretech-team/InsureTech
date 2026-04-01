package models

import (
	"time"
)

// RefreshTokenResponse represents a refresh_token_response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	AccessTokenExpiresIn int `json:"access_token_expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in,omitempty"`
	SessionExpiresAt time.Time `json:"session_expires_at,omitempty"`
	SessionId string `json:"session_id,omitempty"`
}
