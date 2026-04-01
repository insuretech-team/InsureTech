package models


// ReserveInput represents a reserve_input
type ReserveInput struct {
	CalculationMethod string `json:"calculation_method,omitempty"`
	CaseReserve float64 `json:"case_reserve,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ConfidenceLevel float64 `json:"confidence_level,omitempty"`
	DevelopmentPeriods int `json:"development_periods,omitempty"`
	PaidClaims float64 `json:"paid_claims,omitempty"`
	ReportedClaims float64 `json:"reported_claims,omitempty"`
	TriangleDataJson string `json:"triangle_data_json,omitempty"`
}
