package models


// DocumentVerificationRequest represents a document_verification_request
type DocumentVerificationRequest struct {
	RejectionReason string `json:"rejection_reason,omitempty"`
	UserDocumentId string `json:"user_document_id"`
	VerificationStatus string `json:"verification_status,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
