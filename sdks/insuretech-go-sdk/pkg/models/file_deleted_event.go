package models

import (
	"time"
)

// FileDeletedEvent represents a file_deleted_event
type FileDeletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
