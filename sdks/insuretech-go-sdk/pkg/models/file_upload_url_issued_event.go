package models

import (
	"time"
)

// FileUploadURLIssuedEvent represents a file_upload_url_issued_event
type FileUploadURLIssuedEvent struct {
	RequestedBy string `json:"requested_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Filename string `json:"filename,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
}
