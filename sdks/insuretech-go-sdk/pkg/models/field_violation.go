package models


// FieldViolation represents a field_violation
type FieldViolation struct {
	Code string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Field string `json:"field,omitempty"`
	RejectedValue string `json:"rejected_value,omitempty"`
}
