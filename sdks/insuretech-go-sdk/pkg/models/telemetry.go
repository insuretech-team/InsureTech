package models

import (
	"time"
)

// Telemetry represents a telemetry
type Telemetry struct {
	TelemetryId string `json:"telemetry_id"`
	DeviceId string `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	RawData string `json:"raw_data,omitempty"`
	Location *Location `json:"location,omitempty"`
	Type *TelemetryType `json:"type"`
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}
