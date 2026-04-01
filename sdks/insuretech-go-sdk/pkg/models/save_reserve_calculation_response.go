package models


// SaveReserveCalculationResponse represents a save_reserve_calculation_response
type SaveReserveCalculationResponse struct {
	Errors []string `json:"errors,omitempty"`
	Reserve *ReserveCalculation `json:"reserve,omitempty"`
	Success bool `json:"success,omitempty"`
}
