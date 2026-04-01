package models


// TenantConfigUpdateRequest represents a tenant_config_update_request
type TenantConfigUpdateRequest struct {
	ConfigKey string `json:"config_key,omitempty"`
	ConfigValue string `json:"config_value,omitempty"`
	TenantId string `json:"tenant_id"`
	Type string `json:"type"`
}
