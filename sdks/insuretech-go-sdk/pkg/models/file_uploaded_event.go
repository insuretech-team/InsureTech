package models

import (
	"time"
)

// FileUploadedEvent represents a file_uploaded_event
type FileUploadedEvent struct {
	StorageKey string `json:"storage_key,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	FileId string `json:"file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Filename string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Url string `json:"url,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
	Source string `json:"source,omitempty"`
	EventId string `json:"event_id,omitempty"`
	SizeBytes string `json:"size_bytes,omitempty"`
	Bucket string `json:"bucket,omitempty"`
	CdnUrl string `json:"cdn_url,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
}
