package models


// TenantUpdateRequest represents a tenant_update_request
type TenantUpdateRequest struct {
	Branding string `json:"branding,omitempty"`
	Config string `json:"config,omitempty"`
	Name string `json:"name"`
	Status string `json:"status,omitempty"`
	TenantId string `json:"tenant_id"`
}
