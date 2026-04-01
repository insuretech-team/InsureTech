package models


// ReserveResult represents a reserve_result
type ReserveResult struct {
	CaseReserve float64 `json:"case_reserve,omitempty"`
	ExpenseReserve float64 `json:"expense_reserve,omitempty"`
	IbnerReserve float64 `json:"ibner_reserve,omitempty"`
	IbnrReserve float64 `json:"ibnr_reserve,omitempty"`
	LowerBound float64 `json:"lower_bound,omitempty"`
	MethodUsed string `json:"method_used,omitempty"`
	TotalReserve float64 `json:"total_reserve,omitempty"`
	UpperBound float64 `json:"upper_bound,omitempty"`
}
