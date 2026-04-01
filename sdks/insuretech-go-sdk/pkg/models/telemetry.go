package models

import (
	"time"
)

// Telemetry represents a telemetry
type Telemetry struct {
	DeviceId string `json:"device_id"`
	Latitude float64 `json:"latitude,omitempty"`
	Location *Location `json:"location,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Metrics map[string]interface{} `json:"metrics,omitempty"`
	RawData string `json:"raw_data,omitempty"`
	TelemetryId string `json:"telemetry_id"`
	Timestamp time.Time `json:"timestamp"`
	Type *TelemetryType `json:"type"`
}
