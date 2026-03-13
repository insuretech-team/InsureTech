package models


// PaymentRetrievalResponse represents a payment_retrieval_response
type PaymentRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	Payment *Payment `json:"payment,omitempty"`
}
