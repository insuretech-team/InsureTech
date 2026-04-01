package models

import (
	"time"
)

// TokenValidationResponse represents a token_validation_response
type TokenValidationResponse struct {
	ApiKeyScopes []string `json:"api_key_scopes,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
	Valid bool `json:"valid,omitempty"`
}
