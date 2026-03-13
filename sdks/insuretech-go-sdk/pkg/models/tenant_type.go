package models

// TenantType represents a tenant_type
type TenantType string

// TenantType values
const (
	TenantTypeTENANTTYPEUNSPECIFIED TenantType = "TENANT_TYPE_UNSPECIFIED"
	TenantTypeTENANTTYPEPLATFORM  = "TENANT_TYPE_PLATFORM"
	TenantTypeTENANTTYPEPARTNER  = "TENANT_TYPE_PARTNER"
	TenantTypeTENANTTYPEINSURER  = "TENANT_TYPE_INSURER"
	TenantTypeTENANTTYPEAGENT  = "TENANT_TYPE_AGENT"
)
