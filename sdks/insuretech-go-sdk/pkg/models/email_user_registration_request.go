package models


// EmailUserRegistrationRequest represents a email_user_registration_request
type EmailUserRegistrationRequest struct {
	DeviceId string `json:"device_id"`
	Email string `json:"email"`
	FullName string `json:"full_name,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Password string `json:"password"`
	UserType string `json:"user_type,omitempty"`
}
