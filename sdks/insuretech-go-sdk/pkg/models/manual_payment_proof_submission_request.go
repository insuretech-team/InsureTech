package models


// ManualPaymentProofSubmissionRequest represents a manual_payment_proof_submission_request
type ManualPaymentProofSubmissionRequest struct {
	ManualProofFileId string `json:"manual_proof_file_id"`
	Notes string `json:"notes,omitempty"`
	PaymentId string `json:"payment_id"`
	SubmittedBy string `json:"submitted_by,omitempty"`
}
