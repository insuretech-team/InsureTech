package models


// ApiKeyGenerationRequest represents a api_key_generation_request
type ApiKeyGenerationRequest struct {
	OwnerId string `json:"owner_id"`
	Scopes []string `json:"scopes,omitempty"`
	RateLimitPerMinute int `json:"rate_limit_per_minute,omitempty"`
	ExpiresInDays string `json:"expires_in_days,omitempty"`
	IpWhitelist []string `json:"ip_whitelist,omitempty"`
	Name string `json:"name"`
	OwnerType string `json:"owner_type,omitempty"`
}
