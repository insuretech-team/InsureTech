package models


// PortalConfigUpdateRequest represents a portal_config_update_request
type PortalConfigUpdateRequest struct {
	MfaMethods []string `json:"mfa_methods,omitempty"`
	RefreshTokenTtlSeconds int `json:"refresh_token_ttl_seconds,omitempty"`
	AllowConcurrentSessions bool `json:"allow_concurrent_sessions,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	Portal *Portal `json:"portal"`
	MfaRequired bool `json:"mfa_required,omitempty"`
	AccessTokenTtlSeconds int `json:"access_token_ttl_seconds,omitempty"`
	SessionTtlSeconds int `json:"session_ttl_seconds,omitempty"`
	IdleTimeoutSeconds int `json:"idle_timeout_seconds,omitempty"`
	MaxConcurrentSessions int `json:"max_concurrent_sessions,omitempty"`
}
