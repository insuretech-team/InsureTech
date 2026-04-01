package models


// InsurerInsurerUpdateRequest represents a insurer_insurer_update_request
type InsurerInsurerUpdateRequest struct {
	Email string `json:"email"`
	InsurerId string `json:"insurer_id"`
	Name string `json:"name"`
	Phone string `json:"phone,omitempty"`
	Status string `json:"status,omitempty"`
}
