package models


// InsurerConfig represents a insurer_config
type InsurerConfig struct {
	ApiBaseUrl string `json:"api_base_url,omitempty"`
	ApiVersion string `json:"api_version,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	AuthCredentials string `json:"auth_credentials,omitempty"`
	AuthType *AuthenticationType `json:"auth_type,omitempty"`
	AutoUnderwritingEnabled bool `json:"auto_underwriting_enabled,omitempty"`
	BusinessModel string `json:"business_model,omitempty"`
	ClaimSettlementDays int `json:"claim_settlement_days,omitempty"`
	Id string `json:"id"`
	InsurerId string `json:"insurer_id"`
	PaymentTerms string `json:"payment_terms,omitempty"`
	RealTimeClaimNotification bool `json:"real_time_claim_notification,omitempty"`
	UnderwritingThreshold int `json:"underwriting_threshold,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
	WebhookUrl string `json:"webhook_url,omitempty"`
}
