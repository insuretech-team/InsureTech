package models


// Coverage represents a coverage
type Coverage struct {
	CoverageId string `json:"coverage_id,omitempty"`
	Deductible *Money `json:"deductible,omitempty"`
	Description string `json:"description,omitempty"`
	IsIncluded bool `json:"is_included,omitempty"`
	IsOptional bool `json:"is_optional,omitempty"`
	Limit *Money `json:"limit,omitempty"`
	Name string `json:"name,omitempty"`
	Premium *Money `json:"premium,omitempty"`
}
