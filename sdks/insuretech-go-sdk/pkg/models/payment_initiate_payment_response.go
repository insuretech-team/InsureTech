package models

import (
	"time"
)

// PaymentInitiatePaymentResponse represents a payment_initiate_payment_response
type PaymentInitiatePaymentResponse struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	GatewayPageUrl string `json:"gateway_page_url,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentUrl string `json:"payment_url,omitempty"`
	Provider string `json:"provider,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	Status string `json:"status,omitempty"`
	TranId string `json:"tran_id,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}
