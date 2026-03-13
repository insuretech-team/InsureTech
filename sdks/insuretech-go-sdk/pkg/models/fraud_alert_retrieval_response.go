package models


// FraudAlertRetrievalResponse represents a fraud_alert_retrieval_response
type FraudAlertRetrievalResponse struct {
	FraudAlert *FraudAlert `json:"fraud_alert,omitempty"`
	Error *Error `json:"error,omitempty"`
}
