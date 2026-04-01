package models


// ClaimEvaluationResponse represents a claim_evaluation_response
type ClaimEvaluationResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Findings []string `json:"findings,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
	RequiredDocuments []string `json:"required_documents,omitempty"`
	SuggestedAmount *Money `json:"suggested_amount,omitempty"`
}
