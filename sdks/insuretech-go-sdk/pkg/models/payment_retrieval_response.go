package models


// PaymentRetrievalResponse represents a payment_retrieval_response
type PaymentRetrievalResponse struct {
	Payment *Payment `json:"payment,omitempty"`
}
