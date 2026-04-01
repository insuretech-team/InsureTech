package models

import (
	"time"
)

// MediaVirusScanCompletedEvent represents a media_virus_scan_completed_event
type MediaVirusScanCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	ScanResult string `json:"scan_result,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	VirusScanStatus string `json:"virus_scan_status,omitempty"`
}
