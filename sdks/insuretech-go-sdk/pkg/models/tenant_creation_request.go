package models


// TenantCreationRequest represents a tenant_creation_request
type TenantCreationRequest struct {
	Branding string `json:"branding,omitempty"`
	Code string `json:"code,omitempty"`
	Config string `json:"config,omitempty"`
	Name string `json:"name"`
	ParentTenantId string `json:"parent_tenant_id"`
	Type string `json:"type"`
}
