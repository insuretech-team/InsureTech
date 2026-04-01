package models


// CommissionSummary represents a commission_summary
type CommissionSummary struct {
	Count int `json:"count,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	Type string `json:"type,omitempty"`
}
