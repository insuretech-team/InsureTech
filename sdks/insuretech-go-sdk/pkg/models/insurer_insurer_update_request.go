package models


// InsurerInsurerUpdateRequest represents a insurer_insurer_update_request
type InsurerInsurerUpdateRequest struct {
	InsurerId string `json:"insurer_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
	Status string `json:"status,omitempty"`
}
