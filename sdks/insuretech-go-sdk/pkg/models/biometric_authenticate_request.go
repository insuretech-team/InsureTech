package models


// BiometricAuthenticateRequest represents a biometric_authenticate_request
type BiometricAuthenticateRequest struct {
	BiometricToken string `json:"biometric_token,omitempty"`
	DeviceId string `json:"device_id"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
}
