package models

import (
	"time"
)

// FileUpdateRequest represents a file_update_request
type FileUpdateRequest struct {
	TenantId string `json:"tenant_id"`
	FileId string `json:"file_id"`
	Filename string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	ClearExpiresAt bool `json:"clear_expires_at,omitempty"`
	ReferenceId string `json:"reference_id"`
	IsPublic bool `json:"is_public,omitempty"`
}
