package models


// ManualPaymentProofSubmissionRequest represents a manual_payment_proof_submission_request
type ManualPaymentProofSubmissionRequest struct {
	PaymentId string `json:"payment_id"`
	ManualProofFileId string `json:"manual_proof_file_id"`
	SubmittedBy string `json:"submitted_by,omitempty"`
	Notes string `json:"notes,omitempty"`
}
