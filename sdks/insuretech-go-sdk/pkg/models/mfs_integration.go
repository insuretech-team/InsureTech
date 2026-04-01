package models


// MFSIntegration represents a mfs_integration
type MFSIntegration struct {
	ApiBaseUrl string `json:"api_base_url"`
	ApiCredentials string `json:"api_credentials,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Config string `json:"config,omitempty"`
	Id string `json:"id"`
	IsActive bool `json:"is_active,omitempty"`
	MerchantId string `json:"merchant_id,omitempty"`
	Provider interface{} `json:"provider"`
}
