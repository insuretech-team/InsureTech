package models

// IotRiskLevel represents a iot_risk_level
type IotRiskLevel string

// IotRiskLevel values
const (
	IotRiskLevelRISKLEVELUNSPECIFIED IotRiskLevel = "RISK_LEVEL_UNSPECIFIED"
	IotRiskLevelRISKLEVELLOW  = "RISK_LEVEL_LOW"
	IotRiskLevelRISKLEVELMEDIUM  = "RISK_LEVEL_MEDIUM"
	IotRiskLevelRISKLEVELHIGH  = "RISK_LEVEL_HIGH"
	IotRiskLevelRISKLEVELCRITICAL  = "RISK_LEVEL_CRITICAL"
)
