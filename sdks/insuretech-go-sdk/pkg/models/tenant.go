package models


// Tenant represents a tenant
type Tenant struct {
	AuditInfo interface{} `json:"audit_info"`
	Branding string `json:"branding,omitempty"`
	Code string `json:"code"`
	Config string `json:"config,omitempty"`
	Name string `json:"name"`
	ParentTenantId string `json:"parent_tenant_id,omitempty"`
	Status interface{} `json:"status"`
	TenantId string `json:"tenant_id"`
	Type *TenantType `json:"type"`
}
