package models

// VerificationType represents a verification_type
type VerificationType string

// VerificationType values
const (
	VerificationTypeVERIFICATIONTYPEUNSPECIFIED VerificationType = "VERIFICATION_TYPE_UNSPECIFIED"
	VerificationTypeVERIFICATIONTYPEKYC  = "VERIFICATION_TYPE_KYC"
	VerificationTypeVERIFICATIONTYPEKYB  = "VERIFICATION_TYPE_KYB"
)
