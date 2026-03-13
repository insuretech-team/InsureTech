package models


// InsuranceFraudAlertRetrievalResponse represents a insurance_fraud_alert_retrieval_response
type InsuranceFraudAlertRetrievalResponse struct {
	Alert *FraudAlert `json:"alert,omitempty"`
}
