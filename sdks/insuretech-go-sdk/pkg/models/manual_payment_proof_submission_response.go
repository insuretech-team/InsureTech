package models


// ManualPaymentProofSubmissionResponse represents a manual_payment_proof_submission_response
type ManualPaymentProofSubmissionResponse struct {
	PaymentId string `json:"payment_id,omitempty"`
	Status string `json:"status,omitempty"`
	Error *Error `json:"error,omitempty"`
}
