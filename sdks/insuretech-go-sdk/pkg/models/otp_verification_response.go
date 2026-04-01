package models


// OTPVerificationResponse represents a otp_verification_response
type OTPVerificationResponse struct {
	DeviceCredential string `json:"device_credential,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Verified bool `json:"verified,omitempty"`
}
