package models


// InsuranceFraudCaseCreationResponse represents a insurance_fraud_case_creation_response
type InsuranceFraudCaseCreationResponse struct {
	FraudCase *FraudCase `json:"fraud_case,omitempty"`
}
