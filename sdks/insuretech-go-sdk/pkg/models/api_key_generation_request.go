package models


// ApiKeyGenerationRequest represents a api_key_generation_request
type ApiKeyGenerationRequest struct {
	ExpiresInDays string `json:"expires_in_days,omitempty"`
	IpWhitelist []string `json:"ip_whitelist,omitempty"`
	Name string `json:"name"`
	OwnerId string `json:"owner_id"`
	OwnerType string `json:"owner_type,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}
