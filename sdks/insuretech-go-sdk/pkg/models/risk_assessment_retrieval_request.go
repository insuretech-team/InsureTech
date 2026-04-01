package models


// RiskAssessmentRetrievalRequest represents a risk_assessment_retrieval_request
type RiskAssessmentRetrievalRequest struct {
	DeviceId string `json:"device_id"`
	PolicyId string `json:"policy_id"`
}
