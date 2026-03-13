package models

import (
	"time"
)

// UploadFinalizationRequest represents a upload_finalization_request
type UploadFinalizationRequest struct {
	ReferenceType string `json:"reference_type,omitempty"`
	TenantId string `json:"tenant_id"`
	ContentType string `json:"content_type,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	FileId string `json:"file_id"`
	Filename string `json:"filename,omitempty"`
	ReferenceId string `json:"reference_id"`
}
