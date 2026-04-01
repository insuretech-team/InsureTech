package models


// QuoteParameters represents a quote_parameters
type QuoteParameters struct {
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
	AssetValue float64 `json:"asset_value,omitempty"`
	CoverageDurationMonths int `json:"coverage_duration_months,omitempty"`
	CoveragePlan string `json:"coverage_plan,omitempty"`
	CoverageType string `json:"coverage_type,omitempty"`
	OptionalCoverages []*OptionalCoverage `json:"optional_coverages,omitempty"`
}
