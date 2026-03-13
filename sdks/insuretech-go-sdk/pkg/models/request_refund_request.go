package models


// RequestRefundRequest represents a request_refund_request
type RequestRefundRequest struct {
	PolicyId string `json:"policy_id"`
	Reason string `json:"reason,omitempty"`
	ReasonDetails string `json:"reason_details,omitempty"`
}
