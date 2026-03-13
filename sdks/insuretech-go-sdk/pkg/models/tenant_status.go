package models

// TenantStatus represents a tenant_status
type TenantStatus string

// TenantStatus values
const (
	TenantStatusTENANTSTATUSUNSPECIFIED TenantStatus = "TENANT_STATUS_UNSPECIFIED"
	TenantStatusTENANTSTATUSACTIVE  = "TENANT_STATUS_ACTIVE"
	TenantStatusTENANTSTATUSSUSPENDED  = "TENANT_STATUS_SUSPENDED"
	TenantStatusTENANTSTATUSINACTIVE  = "TENANT_STATUS_INACTIVE"
)
