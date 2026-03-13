package models


// InsurerConfig represents a insurer_config
type InsurerConfig struct {
	AutoUnderwritingEnabled bool `json:"auto_underwriting_enabled,omitempty"`
	PaymentTerms string `json:"payment_terms,omitempty"`
	InsurerId string `json:"insurer_id"`
	ApiVersion string `json:"api_version,omitempty"`
	AuthCredentials string `json:"auth_credentials,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
	BusinessModel string `json:"business_model,omitempty"`
	RealTimeClaimNotification bool `json:"real_time_claim_notification,omitempty"`
	ClaimSettlementDays int `json:"claim_settlement_days,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	UnderwritingThreshold int `json:"underwriting_threshold,omitempty"`
	Id string `json:"id"`
	ApiBaseUrl string `json:"api_base_url,omitempty"`
	AuthType *AuthenticationType `json:"auth_type,omitempty"`
	WebhookUrl string `json:"webhook_url,omitempty"`
}
