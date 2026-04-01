package models


// RequestRefundResponse represents a request_refund_response
type RequestRefundResponse struct {
	RefundId string `json:"refund_id,omitempty"`
	RefundNumber string `json:"refund_number,omitempty"`
}
