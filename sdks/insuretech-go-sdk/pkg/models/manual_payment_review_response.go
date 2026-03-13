package models


// ManualPaymentReviewResponse represents a manual_payment_review_response
type ManualPaymentReviewResponse struct {
	PaymentId string `json:"payment_id,omitempty"`
	Status string `json:"status,omitempty"`
	Error *Error `json:"error,omitempty"`
}
