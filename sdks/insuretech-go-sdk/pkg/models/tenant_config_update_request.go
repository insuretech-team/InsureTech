package models


// TenantConfigUpdateRequest represents a tenant_config_update_request
type TenantConfigUpdateRequest struct {
	TenantId string `json:"tenant_id"`
	ConfigKey string `json:"config_key,omitempty"`
	ConfigValue string `json:"config_value,omitempty"`
	Type string `json:"type"`
}
