package models


// RiskAssessmentRequest represents a risk_assessment_request
type RiskAssessmentRequest struct {
	ApplicantData map[string]interface{} `json:"applicant_data,omitempty"`
	ApplicantId string `json:"applicant_id"`
	ProductId string `json:"product_id"`
}
