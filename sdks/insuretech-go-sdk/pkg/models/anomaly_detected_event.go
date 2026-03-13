package models

import (
	"time"
)

// AnomalyDetectedEvent represents a anomaly_detected_event
type AnomalyDetectedEvent struct {
	EventId string `json:"event_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	AnomalyType string `json:"anomaly_type,omitempty"`
	AnomalyDetails map[string]interface{} `json:"anomaly_details,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	Severity string `json:"severity,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
}
