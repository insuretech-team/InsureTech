package models


// Condition represents a condition
type Condition struct {
	Attribute string `json:"attribute,omitempty"`
	Operator *Operator `json:"operator,omitempty"`
	Values []string `json:"values,omitempty"`
}
