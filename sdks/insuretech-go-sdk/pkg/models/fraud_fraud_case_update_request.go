package models


// FraudFraudCaseUpdateRequest represents a fraud_fraud_case_update_request
type FraudFraudCaseUpdateRequest struct {
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
	FraudCaseId string `json:"fraud_case_id"`
	Status string `json:"status,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}
