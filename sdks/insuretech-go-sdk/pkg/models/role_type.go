package models

// RoleType represents a role_type
type RoleType string

// RoleType values
const (
	RoleTypeROLETYPEUNSPECIFIED RoleType = "ROLE_TYPE_UNSPECIFIED"
	RoleTypeROLETYPESYSTEM  = "ROLE_TYPE_SYSTEM"
	RoleTypeROLETYPETENANT  = "ROLE_TYPE_TENANT"
	RoleTypeROLETYPECUSTOM  = "ROLE_TYPE_CUSTOM"
)
