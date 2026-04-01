package models


// MfsInitiatePaymentResponse represents a mfs_initiate_payment_response
type MfsInitiatePaymentResponse struct {
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	PaymentUrl string `json:"payment_url,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}
