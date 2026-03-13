package models

import (
	"time"
)

// PaymentInitiatePaymentResponse represents a payment_initiate_payment_response
type PaymentInitiatePaymentResponse struct {
	PaymentUrl string `json:"payment_url,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	Status string `json:"status,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error *Error `json:"error,omitempty"`
	Provider string `json:"provider,omitempty"`
	GatewayPageUrl string `json:"gateway_page_url,omitempty"`
	TranId string `json:"tran_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}
