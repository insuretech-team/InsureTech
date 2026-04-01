package models

import (
	"time"
)

// APIKeyCreationRequest represents a api_key_creation_request
type APIKeyCreationRequest struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Name string `json:"name"`
	OwnerId string `json:"owner_id"`
	OwnerType string `json:"owner_type,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}
