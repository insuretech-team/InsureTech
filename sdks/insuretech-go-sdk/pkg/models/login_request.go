package models


// LoginRequest represents a login_request
type LoginRequest struct {
	DeviceName string `json:"device_name,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	DeviceType string `json:"device_type,omitempty"`
}
