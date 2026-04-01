package models


// UploadURLRetrievalRequest represents a upload_url_retrieval_request
type UploadURLRetrievalRequest struct {
	ContentType string `json:"content_type,omitempty"`
	ExpiresInMinutes int `json:"expires_in_minutes,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id"`
	ReferenceType string `json:"reference_type,omitempty"`
	TenantId string `json:"tenant_id"`
}
