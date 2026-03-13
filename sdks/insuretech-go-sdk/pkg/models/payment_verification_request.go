package models


// PaymentVerificationRequest represents a payment_verification_request
type PaymentVerificationRequest struct {
	PaymentId string `json:"payment_id"`
	TransactionId string `json:"transaction_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Provider string `json:"provider,omitempty"`
	ValId string `json:"val_id"`
	TranId string `json:"tran_id"`
	SessionKey string `json:"session_key,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	ForceProviderRequery bool `json:"force_provider_requery,omitempty"`
}
