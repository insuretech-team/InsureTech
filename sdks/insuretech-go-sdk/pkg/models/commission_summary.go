package models


// CommissionSummary represents a commission_summary
type CommissionSummary struct {
	Type string `json:"type,omitempty"`
	Count int `json:"count,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
}
