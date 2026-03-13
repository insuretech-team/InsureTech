package models


// InsurerInsurerCreationResponse represents a insurer_insurer_creation_response
type InsurerInsurerCreationResponse struct {
	InsurerId string `json:"insurer_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
