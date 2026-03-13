package models


// InsurerCreationResponse represents a insurer_creation_response
type InsurerCreationResponse struct {
	InsurerId string `json:"insurer_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
