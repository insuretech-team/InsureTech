package models


// StreamTelemetryResponse represents a stream_telemetry_response
type StreamTelemetryResponse struct {
	Received bool `json:"received,omitempty"`
	TelemetryId string `json:"telemetry_id,omitempty"`
}
