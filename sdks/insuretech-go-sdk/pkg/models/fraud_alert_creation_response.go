package models


// FraudAlertCreationResponse represents a fraud_alert_creation_response
type FraudAlertCreationResponse struct {
	Alert *FraudAlert `json:"alert,omitempty"`
}
