package models


// FraudAlertUpdateResponse represents a fraud_alert_update_response
type FraudAlertUpdateResponse struct {
	Alert *FraudAlert `json:"alert,omitempty"`
}
