package models


// MfsInitiatePaymentRequest represents a mfs_initiate_payment_request
type MfsInitiatePaymentRequest struct {
	Amount *Money `json:"amount,omitempty"`
	CustomerMsisdn string `json:"customer_msisdn,omitempty"`
	PaymentId string `json:"payment_id"`
	Provider string `json:"provider,omitempty"`
}
