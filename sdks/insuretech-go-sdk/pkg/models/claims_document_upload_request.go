package models


// ClaimsDocumentUploadRequest represents a claims_document_upload_request
type ClaimsDocumentUploadRequest struct {
	ClaimId string `json:"claim_id"`
	DocumentType string `json:"document_type,omitempty"`
	FileData string `json:"file_data,omitempty"`
	FileName string `json:"file_name,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}
