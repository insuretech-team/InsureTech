package models


// TelemetrySendingResponse represents a telemetry_sending_response
type TelemetrySendingResponse struct {
	TelemetryId string `json:"telemetry_id,omitempty"`
	Accepted bool `json:"accepted,omitempty"`
	Error *Error `json:"error,omitempty"`
}
