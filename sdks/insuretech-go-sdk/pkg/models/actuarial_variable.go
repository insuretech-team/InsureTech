package models


// ActuarialVariable represents a actuarial_variable
type ActuarialVariable struct {
	AllowedValues []string `json:"allowed_values,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Description string `json:"description,omitempty"`
	IsRequired bool `json:"is_required,omitempty"`
	MaxValue float64 `json:"max_value,omitempty"`
	MinValue float64 `json:"min_value,omitempty"`
	VariableName string `json:"variable_name,omitempty"`
	VariableType string `json:"variable_type,omitempty"`
}
