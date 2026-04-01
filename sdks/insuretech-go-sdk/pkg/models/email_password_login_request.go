package models


// EmailPasswordLoginRequest represents a email_password_login_request
type EmailPasswordLoginRequest struct {
	DeviceId string `json:"device_id"`
	DeviceName string `json:"device_name,omitempty"`
	Email string `json:"email"`
	Password string `json:"password"`
}
