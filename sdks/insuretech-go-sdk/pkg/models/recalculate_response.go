package models


// RecalculateResponse represents a recalculate_response
type RecalculateResponse struct {
	Calculation *ActuarialCalculation `json:"calculation,omitempty"`
	Errors []string `json:"errors,omitempty"`
	NewCalculationId string `json:"new_calculation_id,omitempty"`
	Success bool `json:"success,omitempty"`
}
