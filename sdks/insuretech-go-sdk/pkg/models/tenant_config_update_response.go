package models


// TenantConfigUpdateResponse represents a tenant_config_update_response
type TenantConfigUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
