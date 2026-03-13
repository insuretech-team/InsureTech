package models

import (
	"time"
)

// MediaVirusScanCompletedEvent represents a media_virus_scan_completed_event
type MediaVirusScanCompletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	VirusScanStatus string `json:"virus_scan_status,omitempty"`
	ScanResult string `json:"scan_result,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
