package models


// RiskAssessmentResponse represents a risk_assessment_response
type RiskAssessmentResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	RecommendedPremium *Money `json:"recommended_premium,omitempty"`
	RiskCategory string `json:"risk_category,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
}
