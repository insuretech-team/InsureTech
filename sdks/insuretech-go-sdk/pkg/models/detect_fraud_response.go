package models


// DetectFraudResponse represents a detect_fraud_response
type DetectFraudResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	IsSuspicious bool `json:"is_suspicious,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
	RiskIndicators []string `json:"risk_indicators,omitempty"`
}
