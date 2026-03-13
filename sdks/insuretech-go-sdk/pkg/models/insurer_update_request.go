package models


// InsurerUpdateRequest represents a insurer_update_request
type InsurerUpdateRequest struct {
	Phone string `json:"phone,omitempty"`
	Status string `json:"status,omitempty"`
	InsurerId string `json:"insurer_id"`
	Name string `json:"name"`
	Email string `json:"email"`
}
