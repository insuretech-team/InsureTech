package models


// MfsInitiatePaymentRequest represents a mfs_initiate_payment_request
type MfsInitiatePaymentRequest struct {
	Provider string `json:"provider,omitempty"`
	PaymentId string `json:"payment_id"`
	Amount *Money `json:"amount,omitempty"`
	CustomerMsisdn string `json:"customer_msisdn,omitempty"`
}
