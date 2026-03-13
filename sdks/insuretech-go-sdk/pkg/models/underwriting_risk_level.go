package models

// UnderwritingRiskLevel represents a underwriting_risk_level
type UnderwritingRiskLevel string

// UnderwritingRiskLevel values
const (
	UnderwritingRiskLevelRISKLEVELUNSPECIFIED UnderwritingRiskLevel = "RISK_LEVEL_UNSPECIFIED"
	UnderwritingRiskLevelRISKLEVELLOW  = "RISK_LEVEL_LOW"
	UnderwritingRiskLevelRISKLEVELMEDIUM  = "RISK_LEVEL_MEDIUM"
	UnderwritingRiskLevelRISKLEVELHIGH  = "RISK_LEVEL_HIGH"
	UnderwritingRiskLevelRISKLEVELVERYHIGH  = "RISK_LEVEL_VERY_HIGH"
)
