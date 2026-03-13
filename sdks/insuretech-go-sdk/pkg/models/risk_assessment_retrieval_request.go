package models


// RiskAssessmentRetrievalRequest represents a risk_assessment_retrieval_request
type RiskAssessmentRetrievalRequest struct {
	PolicyId string `json:"policy_id"`
	DeviceId string `json:"device_id"`
}
