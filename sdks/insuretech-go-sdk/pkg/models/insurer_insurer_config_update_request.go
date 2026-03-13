package models


// InsurerInsurerConfigUpdateRequest represents a insurer_insurer_config_update_request
type InsurerInsurerConfigUpdateRequest struct {
	BusinessModel string `json:"business_model,omitempty"`
	AutoUnderwritingEnabled bool `json:"auto_underwriting_enabled,omitempty"`
	RealTimeClaimNotification bool `json:"real_time_claim_notification,omitempty"`
	InsurerId string `json:"insurer_id"`
	ApiBaseUrl string `json:"api_base_url,omitempty"`
	AuthType string `json:"auth_type,omitempty"`
	AuthCredentials string `json:"auth_credentials,omitempty"`
}
