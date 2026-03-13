package models

// ClaimProcessingType represents a claim_processing_type
type ClaimProcessingType string

// ClaimProcessingType values
const (
	ClaimProcessingTypeCLAIMPROCESSINGTYPEUNSPECIFIED ClaimProcessingType = "CLAIM_PROCESSING_TYPE_UNSPECIFIED"
	ClaimProcessingTypeCLAIMPROCESSINGTYPEMANUAL  = "CLAIM_PROCESSING_TYPE_MANUAL"
	ClaimProcessingTypeCLAIMPROCESSINGTYPEAUTOADJUDICATED  = "CLAIM_PROCESSING_TYPE_AUTO_ADJUDICATED"
	ClaimProcessingTypeCLAIMPROCESSINGTYPEAIASSISTED  = "CLAIM_PROCESSING_TYPE_AI_ASSISTED"
)
