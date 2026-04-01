package models


// PaymentVerificationRequest represents a payment_verification_request
type PaymentVerificationRequest struct {
	ForceProviderRequery bool `json:"force_provider_requery,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	PaymentId string `json:"payment_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Provider string `json:"provider,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	TranId string `json:"tran_id"`
	TransactionId string `json:"transaction_id"`
	ValId string `json:"val_id"`
}
