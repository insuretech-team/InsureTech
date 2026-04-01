package models


// EmployeeActivationRequest represents a employee_activation_request
type EmployeeActivationRequest struct {
	Email string `json:"email"`
	EmployeeId string `json:"employee_id"`
	OrganisationCode string `json:"organisation_code,omitempty"`
	OrganisationId string `json:"organisation_id"`
}
