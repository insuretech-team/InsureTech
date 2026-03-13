package models

// OrganisationStatus represents a organisation_status
type OrganisationStatus string

// OrganisationStatus values
const (
	OrganisationStatusORGANISATIONSTATUSUNSPECIFIED OrganisationStatus = "ORGANISATION_STATUS_UNSPECIFIED"
	OrganisationStatusORGANISATIONSTATUSACTIVE  = "ORGANISATION_STATUS_ACTIVE"
	OrganisationStatusORGANISATIONSTATUSINACTIVE  = "ORGANISATION_STATUS_INACTIVE"
	OrganisationStatusORGANISATIONSTATUSSUSPENDED  = "ORGANISATION_STATUS_SUSPENDED"
)
