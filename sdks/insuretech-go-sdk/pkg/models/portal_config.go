package models

import (
	"time"
)

// PortalConfig represents a portal_config
type PortalConfig struct {
	AccessTokenTtlSeconds int `json:"access_token_ttl_seconds"`
	AllowConcurrentSessions bool `json:"allow_concurrent_sessions"`
	IdleTimeoutSeconds int `json:"idle_timeout_seconds"`
	MaxConcurrentSessions int `json:"max_concurrent_sessions"`
	MfaMethods []string `json:"mfa_methods,omitempty"`
	MfaRequired bool `json:"mfa_required"`
	Portal interface{} `json:"portal"`
	PortalConfigId string `json:"portal_config_id"`
	RefreshTokenTtlSeconds int `json:"refresh_token_ttl_seconds"`
	SessionTtlSeconds int `json:"session_ttl_seconds"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
