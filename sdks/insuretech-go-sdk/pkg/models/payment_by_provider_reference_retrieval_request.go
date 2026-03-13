package models


// PaymentByProviderReferenceRetrievalRequest represents a payment_by_provider_reference_retrieval_request
type PaymentByProviderReferenceRetrievalRequest struct {
	ProviderReference string `json:"provider_reference,omitempty"`
	Provider string `json:"provider"`
}
