package models

import (
	"time"
)

// MediaFileUploadedEvent represents a media_file_uploaded_event
type MediaFileUploadedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	MediaType string `json:"media_type,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}
