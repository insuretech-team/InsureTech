package models


// DeviceDeactivationResponse represents a device_deactivation_response
type DeviceDeactivationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
