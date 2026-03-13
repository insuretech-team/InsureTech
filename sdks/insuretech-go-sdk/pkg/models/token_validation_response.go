package models

import (
	"time"
)

// TokenValidationResponse represents a token_validation_response
type TokenValidationResponse struct {
	Portal string `json:"portal,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	ApiKeyScopes []string `json:"api_key_scopes,omitempty"`
	Valid bool `json:"valid,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
	Error *Error `json:"error,omitempty"`
	SessionType string `json:"session_type,omitempty"`
}
