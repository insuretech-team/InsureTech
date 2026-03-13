package models


// TenantContext represents a tenant_context
type TenantContext struct {
	TenantId string `json:"tenant_id,omitempty"`
	OrganizationId string `json:"organization_id,omitempty"`
}
