package models


// AddPaymentMethodRequest represents a add_payment_method_request
type AddPaymentMethodRequest struct {
	PaymentMethod *PaymentMethodDetails `json:"payment_method,omitempty"`
	UserId string `json:"user_id"`
}
