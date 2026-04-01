package models


// ManualPaymentReviewRequest represents a manual_payment_review_request
type ManualPaymentReviewRequest struct {
	Approved bool `json:"approved,omitempty"`
	PaymentId string `json:"payment_id"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	ReviewNotes string `json:"review_notes,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
}
