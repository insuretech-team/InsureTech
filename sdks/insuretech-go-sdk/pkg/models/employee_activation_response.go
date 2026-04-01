package models


// EmployeeActivationResponse represents a employee_activation_response
type EmployeeActivationResponse struct {
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
