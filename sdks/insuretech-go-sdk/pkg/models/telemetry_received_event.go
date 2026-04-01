package models

import (
	"time"
)

// TelemetryReceivedEvent represents a telemetry_received_event
type TelemetryReceivedEvent struct {
	DataSizeBytes string `json:"data_size_bytes,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	DeviceTimestamp time.Time `json:"device_timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
	TelemetryId string `json:"telemetry_id,omitempty"`
	TelemetryType string `json:"telemetry_type,omitempty"`
}
