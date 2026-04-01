package models

import (
	"time"
)

// FileUploadRequest represents a file_upload_request
type FileUploadRequest struct {
	Content string `json:"content,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id"`
	ReferenceType string `json:"reference_type,omitempty"`
	TenantId string `json:"tenant_id"`
}
