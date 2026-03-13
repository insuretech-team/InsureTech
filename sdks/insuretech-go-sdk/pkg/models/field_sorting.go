package models


// FieldSorting represents a field_sorting
type FieldSorting struct {
	Direction *SortDirection `json:"direction,omitempty"`
	Field string `json:"field,omitempty"`
}
