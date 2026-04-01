package models


// KycDocumentUploadRequest represents a kyc_document_upload_request
type KycDocumentUploadRequest struct {
	DocumentNumber string `json:"document_number,omitempty"`
	DocumentType string `json:"document_type,omitempty"`
	DocumentUrl string `json:"document_url,omitempty"`
	KycVerificationId string `json:"kyc_verification_id"`
}
