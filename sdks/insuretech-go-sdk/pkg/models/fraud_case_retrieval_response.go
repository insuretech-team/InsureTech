package models


// FraudCaseRetrievalResponse represents a fraud_case_retrieval_response
type FraudCaseRetrievalResponse struct {
	FraudCase *FraudCase `json:"fraud_case,omitempty"`
	Error *Error `json:"error,omitempty"`
}
