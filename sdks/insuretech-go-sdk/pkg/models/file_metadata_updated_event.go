package models

import (
	"time"
)

// FileMetadataUpdatedEvent represents a file_metadata_updated_event
type FileMetadataUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	UpdatedFields []string `json:"updated_fields,omitempty"`
}
