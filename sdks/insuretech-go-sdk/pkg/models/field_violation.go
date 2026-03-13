package models


// FieldViolation represents a field_violation
type FieldViolation struct {
	Field string `json:"field,omitempty"`
	Code string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	RejectedValue string `json:"rejected_value,omitempty"`
}
