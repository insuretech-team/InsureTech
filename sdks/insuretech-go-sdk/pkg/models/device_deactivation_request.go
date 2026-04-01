package models


// DeviceDeactivationRequest represents a device_deactivation_request
type DeviceDeactivationRequest struct {
	DeviceId string `json:"device_id"`
	Reason string `json:"reason,omitempty"`
}
