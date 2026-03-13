package models


// RefundRetrievalResponse represents a refund_retrieval_response
type RefundRetrievalResponse struct {
	Refund *Refund `json:"refund,omitempty"`
	Error *Error `json:"error,omitempty"`
}
