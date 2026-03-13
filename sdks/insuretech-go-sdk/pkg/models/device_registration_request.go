package models


// DeviceRegistrationRequest represents a device_registration_request
type DeviceRegistrationRequest struct {
	DeviceSerial string `json:"device_serial,omitempty"`
	Type *IotDeviceType `json:"type"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model string `json:"model,omitempty"`
	PolicyId string `json:"policy_id"`
	OwnerId string `json:"owner_id"`
}
