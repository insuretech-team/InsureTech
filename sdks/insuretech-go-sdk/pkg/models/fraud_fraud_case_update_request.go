package models


// FraudFraudCaseUpdateRequest represents a fraud_fraud_case_update_request
type FraudFraudCaseUpdateRequest struct {
	Evidence map[string]interface{} `json:"evidence,omitempty"`
	FraudCaseId string `json:"fraud_case_id"`
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	Outcome string `json:"outcome,omitempty"`
	Status string `json:"status,omitempty"`
}
