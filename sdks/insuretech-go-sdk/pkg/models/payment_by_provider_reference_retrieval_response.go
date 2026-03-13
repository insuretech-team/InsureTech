package models


// PaymentByProviderReferenceRetrievalResponse represents a payment_by_provider_reference_retrieval_response
type PaymentByProviderReferenceRetrievalResponse struct {
	Payment *Payment `json:"payment,omitempty"`
	Error *Error `json:"error,omitempty"`
}
