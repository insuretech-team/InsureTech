package models

// ContactStatus represents a contact_status
type ContactStatus string

// ContactStatus values
const (
	ContactStatusCONTACTSTATUSUNSPECIFIED ContactStatus = "CONTACT_STATUS_UNSPECIFIED"
	ContactStatusCONTACTSTATUSACTIVE  = "CONTACT_STATUS_ACTIVE"
	ContactStatusCONTACTSTATUSINACTIVE  = "CONTACT_STATUS_INACTIVE"
	ContactStatusCONTACTSTATUSPROSPECT  = "CONTACT_STATUS_PROSPECT"
	ContactStatusCONTACTSTATUSARCHIVED  = "CONTACT_STATUS_ARCHIVED"
)
