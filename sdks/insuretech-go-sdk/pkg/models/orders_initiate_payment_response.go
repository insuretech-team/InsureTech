package models

import (
	"time"
)

// OrdersInitiatePaymentResponse represents a orders_initiate_payment_response
type OrdersInitiatePaymentResponse struct {
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error *Error `json:"error,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentUrl string `json:"payment_url,omitempty"`
}
