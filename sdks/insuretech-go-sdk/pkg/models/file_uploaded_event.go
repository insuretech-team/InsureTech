package models

import (
	"time"
)

// FileUploadedEvent represents a file_uploaded_event
type FileUploadedEvent struct {
	Bucket string `json:"bucket,omitempty"`
	CdnUrl string `json:"cdn_url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	SizeBytes string `json:"size_bytes,omitempty"`
	Source string `json:"source,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
	Url string `json:"url,omitempty"`
}
