package models


// ReservesCalculationRequest represents a reserves_calculation_request
type ReservesCalculationRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	ClaimId string `json:"claim_id"`
	Input *ReserveInput `json:"input,omitempty"`
}
