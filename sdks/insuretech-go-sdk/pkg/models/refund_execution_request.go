package models


// RefundExecutionRequest represents a refund_execution_request
type RefundExecutionRequest struct {
	Amount *Money `json:"amount,omitempty"`
	MfsTransactionId string `json:"mfs_transaction_id"`
	Reason string `json:"reason,omitempty"`
}
