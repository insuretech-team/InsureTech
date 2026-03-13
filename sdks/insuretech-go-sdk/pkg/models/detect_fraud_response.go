package models


// DetectFraudResponse represents a detect_fraud_response
type DetectFraudResponse struct {
	IsSuspicious bool `json:"is_suspicious,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	RiskIndicators []string `json:"risk_indicators,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
	Error *Error `json:"error,omitempty"`
	AnalysisId string `json:"analysis_id,omitempty"`
}
