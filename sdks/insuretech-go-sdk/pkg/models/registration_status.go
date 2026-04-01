package models

// RegistrationStatus represents a registration_status
type RegistrationStatus string

// RegistrationStatus values
const (
	RegistrationStatusREGISTRATIONSTATUSUNSPECIFIED RegistrationStatus = "REGISTRATION_STATUS_UNSPECIFIED"
	RegistrationStatusREGISTRATIONSTATUSACTIVE  = "REGISTRATION_STATUS_ACTIVE"
	RegistrationStatusREGISTRATIONSTATUSSUSPENDED  = "REGISTRATION_STATUS_SUSPENDED"
	RegistrationStatusREGISTRATIONSTATUSEXPIRED  = "REGISTRATION_STATUS_EXPIRED"
)
