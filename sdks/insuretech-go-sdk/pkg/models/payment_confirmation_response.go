package models


// PaymentConfirmationResponse represents a payment_confirmation_response
type PaymentConfirmationResponse struct {
	OrderId string `json:"order_id,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
