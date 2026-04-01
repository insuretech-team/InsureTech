package models

import (
	"time"
)

// FileUploadFinalizedEvent represents a file_upload_finalized_event
type FileUploadFinalizedEvent struct {
	ContentType string `json:"content_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	FinalizedBy string `json:"finalized_by,omitempty"`
	SizeBytes string `json:"size_bytes,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
