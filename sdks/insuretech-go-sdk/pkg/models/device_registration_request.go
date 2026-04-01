package models


// DeviceRegistrationRequest represents a device_registration_request
type DeviceRegistrationRequest struct {
	DeviceSerial string `json:"device_serial,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model string `json:"model,omitempty"`
	OwnerId string `json:"owner_id"`
	PolicyId string `json:"policy_id"`
	Type *IoTDeviceType `json:"type"`
}
