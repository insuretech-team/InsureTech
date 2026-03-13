package models


// FraudCaseUpdateRequest represents a fraud_case_update_request
type FraudCaseUpdateRequest struct {
	FraudCaseId string `json:"fraud_case_id"`
	Status string `json:"status,omitempty"`
	Outcome string `json:"outcome,omitempty"`
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}
