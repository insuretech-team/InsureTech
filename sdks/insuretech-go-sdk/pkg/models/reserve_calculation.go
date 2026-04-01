package models

import (
	"time"
)

// ReserveCalculation represents a reserve_calculation
type ReserveCalculation struct {
	CalculationMethod string `json:"calculation_method,omitempty"`
	CaseReserve *Money `json:"case_reserve,omitempty"`
	ClaimId string `json:"claim_id"`
	ConfidenceLevel float64 `json:"confidence_level,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	ExpenseReserve *Money `json:"expense_reserve,omitempty"`
	IbnrReserve *Money `json:"ibnr_reserve,omitempty"`
	LowerBound *Money `json:"lower_bound,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PolicyId string `json:"policy_id"`
	ReserveId string `json:"reserve_id"`
	ReserveType *ReserveType `json:"reserve_type"`
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	Status interface{} `json:"status"`
	TotalReserve *Money `json:"total_reserve,omitempty"`
	TriangleDataJson string `json:"triangle_data_json,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	UpperBound *Money `json:"upper_bound,omitempty"`
}
