package models

import (
	"time"
)

// ReservesCalculationResponse represents a reserves_calculation_response
type ReservesCalculationResponse struct {
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	Errors []string `json:"errors,omitempty"`
	ReserveId string `json:"reserve_id,omitempty"`
	Result *ReserveResult `json:"result,omitempty"`
	Success bool `json:"success,omitempty"`
}
