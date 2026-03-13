package models


// FileRetrievalRequest represents a file_retrieval_request
type FileRetrievalRequest struct {
	FileId string `json:"file_id"`
	TenantId string `json:"tenant_id"`
}
