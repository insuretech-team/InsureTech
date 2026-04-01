package models

import (
	"time"
)

// OrdersInitiatePaymentResponse represents a orders_initiate_payment_response
type OrdersInitiatePaymentResponse struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentUrl string `json:"payment_url,omitempty"`
}
