package models


// LoginRequest represents a login_request
type LoginRequest struct {
	DeviceId string `json:"device_id"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Password string `json:"password"`
}
