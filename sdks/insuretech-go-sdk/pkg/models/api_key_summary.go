package models

import (
	"time"
)

// APIKeySummary represents a api_key_summary
type APIKeySummary struct {
	KeyId string `json:"key_id,omitempty"`
	Name string `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	Status string `json:"status,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
}
