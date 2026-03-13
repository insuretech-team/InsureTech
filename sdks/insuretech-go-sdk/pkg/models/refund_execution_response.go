package models


// RefundExecutionResponse represents a refund_execution_response
type RefundExecutionResponse struct {
	RefundTransactionId string `json:"refund_transaction_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
