package models


// CoverageComparison represents a coverage_comparison
type CoverageComparison struct {
	CoverageId string `json:"coverage_id,omitempty"`
	IsIncluded bool `json:"is_included,omitempty"`
	Name string `json:"name,omitempty"`
	Premium *Money `json:"premium,omitempty"`
}
