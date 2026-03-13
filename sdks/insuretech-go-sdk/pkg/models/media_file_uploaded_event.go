package models

import (
	"time"
)

// MediaFileUploadedEvent represents a media_file_uploaded_event
type MediaFileUploadedEvent struct {
	EventId string `json:"event_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	MediaType string `json:"media_type,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	FileId string `json:"file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
