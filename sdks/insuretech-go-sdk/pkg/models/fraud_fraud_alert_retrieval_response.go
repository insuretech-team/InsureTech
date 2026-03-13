package models


// FraudFraudAlertRetrievalResponse represents a fraud_fraud_alert_retrieval_response
type FraudFraudAlertRetrievalResponse struct {
	FraudAlert *FraudAlert `json:"fraud_alert,omitempty"`
	Error *Error `json:"error,omitempty"`
}
