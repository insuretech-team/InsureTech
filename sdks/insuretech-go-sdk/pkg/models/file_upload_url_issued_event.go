package models

import (
	"time"
)

// FileUploadURLIssuedEvent represents a file_upload_url_issued_event
type FileUploadURLIssuedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	FileId string `json:"file_id,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
