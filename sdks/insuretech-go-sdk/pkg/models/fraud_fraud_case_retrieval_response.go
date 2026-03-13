package models


// FraudFraudCaseRetrievalResponse represents a fraud_fraud_case_retrieval_response
type FraudFraudCaseRetrievalResponse struct {
	FraudCase *FraudCase `json:"fraud_case,omitempty"`
	Error *Error `json:"error,omitempty"`
}
