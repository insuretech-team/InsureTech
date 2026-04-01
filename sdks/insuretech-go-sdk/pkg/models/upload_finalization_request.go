package models

import (
	"time"
)

// UploadFinalizationRequest represents a upload_finalization_request
type UploadFinalizationRequest struct {
	ContentType string `json:"content_type,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	FileId string `json:"file_id"`
	FileType *FileType `json:"file_type,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id"`
	ReferenceType string `json:"reference_type,omitempty"`
	TenantId string `json:"tenant_id"`
}
