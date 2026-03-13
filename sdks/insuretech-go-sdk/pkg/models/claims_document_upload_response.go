package models


// ClaimsDocumentUploadResponse represents a claims_document_upload_response
type ClaimsDocumentUploadResponse struct {
	DocumentId string `json:"document_id,omitempty"`
	DocumentUrl string `json:"document_url,omitempty"`
	FileHash string `json:"file_hash,omitempty"`
	Error *Error `json:"error,omitempty"`
}
