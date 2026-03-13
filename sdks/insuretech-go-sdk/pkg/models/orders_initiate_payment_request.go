package models


// OrdersInitiatePaymentRequest represents a orders_initiate_payment_request
type OrdersInitiatePaymentRequest struct {
	CallbackUrl string `json:"callback_url,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	OrderId string `json:"order_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
}
