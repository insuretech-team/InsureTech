package models


// CheckFraudResponse represents a check_fraud_response
type CheckFraudResponse struct {
	FraudAlertId string `json:"fraud_alert_id,omitempty"`
	FraudScore int `json:"fraud_score,omitempty"`
	IsFraudDetected bool `json:"is_fraud_detected,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	TriggeredRules []string `json:"triggered_rules,omitempty"`
}
