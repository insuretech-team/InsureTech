package models


// TenantContext represents a tenant_context
type TenantContext struct {
	OrganizationId string `json:"organization_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
}
