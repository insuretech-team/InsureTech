package models


// RiskAssessmentRetrievalResponse represents a risk_assessment_retrieval_response
type RiskAssessmentRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	Assessment *RiskAssessment `json:"assessment,omitempty"`
}
