package models


// EmailUserRegistrationRequest represents a email_user_registration_request
type EmailUserRegistrationRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type,omitempty"`
	FullName string `json:"full_name,omitempty"`
	DeviceId string `json:"device_id"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
