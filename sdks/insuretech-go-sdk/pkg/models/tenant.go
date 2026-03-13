package models


// Tenant represents a tenant
type Tenant struct {
	TenantId string `json:"tenant_id"`
	Name string `json:"name"`
	Type *TenantType `json:"type"`
	ParentTenantId string `json:"parent_tenant_id,omitempty"`
	Branding string `json:"branding,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Code string `json:"code"`
	Status interface{} `json:"status"`
	Config string `json:"config,omitempty"`
}
