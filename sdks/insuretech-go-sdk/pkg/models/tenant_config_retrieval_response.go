package models


// TenantConfigRetrievalResponse represents a tenant_config_retrieval_response
type TenantConfigRetrievalResponse struct {
	Configs []*TenantConfig `json:"configs,omitempty"`
}
