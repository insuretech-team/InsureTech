package models


// RiskAssessmentResponse represents a risk_assessment_response
type RiskAssessmentResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
	RiskCategory string `json:"risk_category,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	RecommendedPremium *Money `json:"recommended_premium,omitempty"`
	Error *Error `json:"error,omitempty"`
}
