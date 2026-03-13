package models

import (
	"time"
)

// PortalConfig represents a portal_config
type PortalConfig struct {
	MfaMethods []string `json:"mfa_methods,omitempty"`
	SessionTtlSeconds int `json:"session_ttl_seconds"`
	MaxConcurrentSessions int `json:"max_concurrent_sessions"`
	Portal interface{} `json:"portal"`
	MfaRequired bool `json:"mfa_required"`
	AccessTokenTtlSeconds int `json:"access_token_ttl_seconds"`
	RefreshTokenTtlSeconds int `json:"refresh_token_ttl_seconds"`
	IdleTimeoutSeconds int `json:"idle_timeout_seconds"`
	AllowConcurrentSessions bool `json:"allow_concurrent_sessions"`
	UpdatedBy string `json:"updated_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	PortalConfigId string `json:"portal_config_id"`
}
