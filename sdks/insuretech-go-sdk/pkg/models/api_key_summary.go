package models

import (
	"time"
)

// APIKeySummary represents a api_key_summary
type APIKeySummary struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	KeyId string `json:"key_id,omitempty"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
	Name string `json:"name,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	Status string `json:"status,omitempty"`
}
