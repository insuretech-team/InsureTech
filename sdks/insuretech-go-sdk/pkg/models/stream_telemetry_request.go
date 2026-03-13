package models


// StreamTelemetryRequest represents a stream_telemetry_request
type StreamTelemetryRequest struct {
	Telemetry *Telemetry `json:"telemetry"`
}
