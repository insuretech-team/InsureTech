package models


// ProvisionEmployeeUserResponse represents a provision_employee_user_response
type ProvisionEmployeeUserResponse struct {
	Created bool `json:"created,omitempty"`
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
