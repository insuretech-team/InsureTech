package models

import (
	"time"
)

// APIKeyCreationRequest represents a api_key_creation_request
type APIKeyCreationRequest struct {
	Name string `json:"name"`
	OwnerId string `json:"owner_id"`
	OwnerType string `json:"owner_type,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
