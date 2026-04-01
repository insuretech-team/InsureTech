package models


// SaveReserveCalculationRequest represents a save_reserve_calculation_request
type SaveReserveCalculationRequest struct {
	Reserve *ReserveCalculation `json:"reserve"`
}
