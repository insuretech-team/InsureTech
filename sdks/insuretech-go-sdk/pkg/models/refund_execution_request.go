package models


// RefundExecutionRequest represents a refund_execution_request
type RefundExecutionRequest struct {
	MfsTransactionId string `json:"mfs_transaction_id"`
	Amount *Money `json:"amount,omitempty"`
	Reason string `json:"reason,omitempty"`
}
