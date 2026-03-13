package models


// FileDeletionRequest represents a file_deletion_request
type FileDeletionRequest struct {
	TenantId string `json:"tenant_id"`
	FileId string `json:"file_id"`
}
