package models


// PaymentConfirmationRequest represents a payment_confirmation_request
type PaymentConfirmationRequest struct {
	OrderId string `json:"order_id"`
	PaymentId string `json:"payment_id"`
	TransactionId string `json:"transaction_id"`
}
