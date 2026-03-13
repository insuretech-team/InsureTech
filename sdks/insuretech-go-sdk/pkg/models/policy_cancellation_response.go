package models


// PolicyCancellationResponse represents a policy_cancellation_response
type PolicyCancellationResponse struct {
	Message string `json:"message,omitempty"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
	Error *Error `json:"error,omitempty"`
}
