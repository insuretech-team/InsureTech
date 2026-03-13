package models

// FraudRiskLevel represents a fraud_risk_level
type FraudRiskLevel string

// FraudRiskLevel values
const (
	FraudRiskLevelRISKLEVELUNSPECIFIED FraudRiskLevel = "RISK_LEVEL_UNSPECIFIED"
	FraudRiskLevelRISKLEVELLOW  = "RISK_LEVEL_LOW"
	FraudRiskLevelRISKLEVELMEDIUM  = "RISK_LEVEL_MEDIUM"
	FraudRiskLevelRISKLEVELHIGH  = "RISK_LEVEL_HIGH"
	FraudRiskLevelRISKLEVELCRITICAL  = "RISK_LEVEL_CRITICAL"
)
