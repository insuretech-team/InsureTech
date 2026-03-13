package models


// CommissionRetrievalResponse represents a commission_retrieval_response
type CommissionRetrievalResponse struct {
	Commission *Commission `json:"commission,omitempty"`
	Error *Error `json:"error,omitempty"`
}
