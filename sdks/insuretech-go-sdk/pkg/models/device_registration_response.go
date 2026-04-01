package models


// DeviceRegistrationResponse represents a device_registration_response
type DeviceRegistrationResponse struct {
	ActivationToken string `json:"activation_token,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
}
