package models

// ContactType represents a contact_type
type ContactType string

// ContactType values
const (
	ContactTypeCONTACTTYPEUNSPECIFIED ContactType = "CONTACT_TYPE_UNSPECIFIED"
	ContactTypeCONTACTTYPEINDIVIDUAL  = "CONTACT_TYPE_INDIVIDUAL"
	ContactTypeCONTACTTYPEBUSINESS  = "CONTACT_TYPE_BUSINESS"
	ContactTypeCONTACTTYPEFAMILY  = "CONTACT_TYPE_FAMILY"
)
