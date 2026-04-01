package models


// FraudFraudCaseCreationRequest represents a fraud_fraud_case_creation_request
type FraudFraudCaseCreationRequest struct {
	FraudAlertId string `json:"fraud_alert_id"`
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	InvestigatorId string `json:"investigator_id"`
	Priority string `json:"priority,omitempty"`
}
