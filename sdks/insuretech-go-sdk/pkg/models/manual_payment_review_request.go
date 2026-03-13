package models


// ManualPaymentReviewRequest represents a manual_payment_review_request
type ManualPaymentReviewRequest struct {
	PaymentId string `json:"payment_id"`
	Approved bool `json:"approved,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	ReviewNotes string `json:"review_notes,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}
