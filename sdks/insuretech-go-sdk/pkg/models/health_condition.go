package models


// HealthCondition represents a health_condition
type HealthCondition struct {
	ConditionCode string `json:"condition_code,omitempty"`
	ConditionName string `json:"condition_name,omitempty"`
	DiagnosisDate string `json:"diagnosis_date,omitempty"`
	Notes string `json:"notes,omitempty"`
	Severity string `json:"severity,omitempty"`
}
