package models


// RequestRefundResponse represents a request_refund_response
type RequestRefundResponse struct {
	RefundId string `json:"refund_id,omitempty"`
	RefundNumber string `json:"refund_number,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
