package models


// KycDocumentUploadResponse represents a kyc_document_upload_response
type KycDocumentUploadResponse struct {
	DocumentVerificationId string `json:"document_verification_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
