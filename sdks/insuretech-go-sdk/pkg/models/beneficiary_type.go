package models

// BeneficiaryType represents a beneficiary_type
type BeneficiaryType string

// BeneficiaryType values
const (
	BeneficiaryTypeBENEFICIARYTYPEUNSPECIFIED BeneficiaryType = "BENEFICIARY_TYPE_UNSPECIFIED"
	BeneficiaryTypeBENEFICIARYTYPEINDIVIDUAL  = "BENEFICIARY_TYPE_INDIVIDUAL"
	BeneficiaryTypeBENEFICIARYTYPEBUSINESS  = "BENEFICIARY_TYPE_BUSINESS"
)
