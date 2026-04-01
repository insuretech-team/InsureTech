package models

import (
	"time"
)

// StoredFile represents a stored_file
type StoredFile struct {
	Bucket string `json:"bucket,omitempty"`
	CdnUrl string `json:"cdn_url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	FileId string `json:"file_id,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	SizeBytes string `json:"size_bytes,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
	Url string `json:"url,omitempty"`
}
