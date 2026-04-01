package models


// FileDeletionRequest represents a file_deletion_request
type FileDeletionRequest struct {
	FileId string `json:"file_id"`
	TenantId string `json:"tenant_id"`
}
