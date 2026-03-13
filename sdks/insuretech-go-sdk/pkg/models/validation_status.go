package models

// ValidationStatus represents a validation_status
type ValidationStatus string

// ValidationStatus values
const (
	ValidationStatusVALIDATIONSTATUSUNSPECIFIED ValidationStatus = "VALIDATION_STATUS_UNSPECIFIED"
	ValidationStatusVALIDATIONSTATUSPENDING  = "VALIDATION_STATUS_PENDING"
	ValidationStatusVALIDATIONSTATUSVALIDATED  = "VALIDATION_STATUS_VALIDATED"
	ValidationStatusVALIDATIONSTATUSREJECTED  = "VALIDATION_STATUS_REJECTED"
)
