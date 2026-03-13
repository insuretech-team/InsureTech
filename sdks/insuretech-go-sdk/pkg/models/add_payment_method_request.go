package models


// AddPaymentMethodRequest represents a add_payment_method_request
type AddPaymentMethodRequest struct {
	UserId string `json:"user_id"`
	PaymentMethod *PaymentPaymentMethod `json:"payment_method,omitempty"`
}
