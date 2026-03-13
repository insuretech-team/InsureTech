package models

import (
	"time"
)

// ApiKey represents a api_key
type ApiKey struct {
	Scopes []string `json:"scopes,omitempty"`
	Status interface{} `json:"status"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
	IpWhitelist []string `json:"ip_whitelist,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	KeyHash string `json:"key_hash"`
	Name string `json:"name"`
	OwnerType *ApiKeyOwnerType `json:"owner_type"`
	RateLimitPerMinute int `json:"rate_limit_per_minute"`
	Id string `json:"id"`
	OwnerId string `json:"owner_id"`
}
