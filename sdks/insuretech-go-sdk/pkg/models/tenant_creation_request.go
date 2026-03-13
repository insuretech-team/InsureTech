package models


// TenantCreationRequest represents a tenant_creation_request
type TenantCreationRequest struct {
	Branding string `json:"branding,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
	ParentTenantId string `json:"parent_tenant_id"`
	Config string `json:"config,omitempty"`
}
