package models


// OptionalCoverage represents a optional_coverage
type OptionalCoverage struct {
	CoverageId string `json:"coverage_id,omitempty"`
	Name string `json:"name,omitempty"`
	SelectedDeductible float64 `json:"selected_deductible,omitempty"`
	SelectedLimit float64 `json:"selected_limit,omitempty"`
}
