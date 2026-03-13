package models


// KycDocumentUploadRequest represents a kyc_document_upload_request
type KycDocumentUploadRequest struct {
	KycVerificationId string `json:"kyc_verification_id"`
	DocumentType string `json:"document_type,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
	DocumentUrl string `json:"document_url,omitempty"`
}
