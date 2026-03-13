package models

import (
	"time"
)

// TelemetryReceivedEvent represents a telemetry_received_event
type TelemetryReceivedEvent struct {
	DeviceTimestamp time.Time `json:"device_timestamp,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
	DataSizeBytes string `json:"data_size_bytes,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TelemetryId string `json:"telemetry_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	TelemetryType string `json:"telemetry_type,omitempty"`
}
