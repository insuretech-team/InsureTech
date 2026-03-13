package models

// MaritalStatus represents a marital_status
type MaritalStatus string

// MaritalStatus values
const (
	MaritalStatusMARITALSTATUSUNSPECIFIED MaritalStatus = "MARITAL_STATUS_UNSPECIFIED"
	MaritalStatusMARITALSTATUSSINGLE  = "MARITAL_STATUS_SINGLE"
	MaritalStatusMARITALSTATUSMARRIED  = "MARITAL_STATUS_MARRIED"
	MaritalStatusMARITALSTATUSDIVORCED  = "MARITAL_STATUS_DIVORCED"
	MaritalStatusMARITALSTATUSWIDOWED  = "MARITAL_STATUS_WIDOWED"
)
