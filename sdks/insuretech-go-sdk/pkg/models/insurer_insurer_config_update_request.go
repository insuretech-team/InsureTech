package models


// InsurerInsurerConfigUpdateRequest represents a insurer_insurer_config_update_request
type InsurerInsurerConfigUpdateRequest struct {
	ApiBaseUrl string `json:"api_base_url,omitempty"`
	AuthCredentials string `json:"auth_credentials,omitempty"`
	AuthType string `json:"auth_type,omitempty"`
	AutoUnderwritingEnabled bool `json:"auto_underwriting_enabled,omitempty"`
	BusinessModel string `json:"business_model,omitempty"`
	InsurerId string `json:"insurer_id"`
	RealTimeClaimNotification bool `json:"real_time_claim_notification,omitempty"`
}
