package models

// Portal represents a portal
type Portal string

// Portal values
const (
	PortalPORTALUNSPECIFIED Portal = "PORTAL_UNSPECIFIED"
	PortalPORTALSYSTEM  = "PORTAL_SYSTEM"
	PortalPORTALBUSINESS  = "PORTAL_BUSINESS"
	PortalPORTALB2B  = "PORTAL_B2B"
	PortalPORTALAGENT  = "PORTAL_AGENT"
	PortalPORTALREGULATOR  = "PORTAL_REGULATOR"
	PortalPORTALB2C  = "PORTAL_B2C"
)
