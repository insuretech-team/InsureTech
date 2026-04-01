package models

import (
	"time"
)

// AnomalyDetectedEvent represents a anomaly_detected_event
type AnomalyDetectedEvent struct {
	AnomalyDetails map[string]interface{} `json:"anomaly_details,omitempty"`
	AnomalyType string `json:"anomaly_type,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Severity string `json:"severity,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
