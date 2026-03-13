package models


// RegistrationRequest represents a registration_request
type RegistrationRequest struct {
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	DeviceType string `json:"device_type,omitempty"`
}
