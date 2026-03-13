package models


// DeviceDeactivationRequest represents a device_deactivation_request
type DeviceDeactivationRequest struct {
	Reason string `json:"reason,omitempty"`
	DeviceId string `json:"device_id"`
}
