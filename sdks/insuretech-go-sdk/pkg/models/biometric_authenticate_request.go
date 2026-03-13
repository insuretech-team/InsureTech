package models


// BiometricAuthenticateRequest represents a biometric_authenticate_request
type BiometricAuthenticateRequest struct {
	DeviceId string `json:"device_id"`
	DeviceType string `json:"device_type,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	BiometricToken string `json:"biometric_token,omitempty"`
}
