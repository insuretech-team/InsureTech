package models

import (
	"time"
)

// FileUploadFinalizedEvent represents a file_upload_finalized_event
type FileUploadFinalizedEvent struct {
	FileId string `json:"file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	SizeBytes string `json:"size_bytes,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	FinalizedBy string `json:"finalized_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
