package models


// TenantConfig represents a tenant_config
type TenantConfig struct {
	AuditInfo interface{} `json:"audit_info"`
	ConfigKey string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Description string `json:"description,omitempty"`
	Id string `json:"id"`
	TenantId string `json:"tenant_id"`
	Type *ConfigType `json:"type"`
}
