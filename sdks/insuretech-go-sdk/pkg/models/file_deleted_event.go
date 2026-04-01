package models

import (
	"time"
)

// FileDeletedEvent represents a file_deleted_event
type FileDeletedEvent struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileId string `json:"file_id,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
