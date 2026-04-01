package models

import (
	"time"
)

// ApiKey represents a api_key
type ApiKey struct {
	AuditInfo interface{} `json:"audit_info"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Id string `json:"id"`
	IpWhitelist []string `json:"ip_whitelist,omitempty"`
	KeyHash string `json:"key_hash"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
	Name string `json:"name"`
	OwnerId string `json:"owner_id"`
	OwnerType *ApiKeyOwnerType `json:"owner_type"`
	RateLimitPerMinute int `json:"rate_limit_per_minute"`
	Scopes []string `json:"scopes,omitempty"`
	Status interface{} `json:"status"`
}
