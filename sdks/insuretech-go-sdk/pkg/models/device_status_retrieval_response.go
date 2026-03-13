package models

import (
	"time"
)

// DeviceStatusRetrievalResponse represents a device_status_retrieval_response
type DeviceStatusRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	Device *IoTDevice `json:"device,omitempty"`
	LastTelemetryAt time.Time `json:"last_telemetry_at,omitempty"`
	TelemetryCountToday int `json:"telemetry_count_today,omitempty"`
}
