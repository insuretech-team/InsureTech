package models


// FraudAlertCreationRequest represents a fraud_alert_creation_request
type FraudAlertCreationRequest struct {
	Alert *FraudAlert `json:"alert"`
}
