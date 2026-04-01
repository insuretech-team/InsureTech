package models


// CommissionCalculationResponse represents a commission_calculation_response
type CommissionCalculationResponse struct {
	Amount *Money `json:"amount,omitempty"`
	CalculationBreakdown string `json:"calculation_breakdown,omitempty"`
	CommissionId string `json:"commission_id,omitempty"`
	CommissionNumber string `json:"commission_number,omitempty"`
}
