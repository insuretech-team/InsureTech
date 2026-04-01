package models


// RecalculateRequest represents a recalculate_request
type RecalculateRequest struct {
	CalculationId string `json:"calculation_id"`
	Reason string `json:"reason,omitempty"`
	RecalculatedBy string `json:"recalculated_by,omitempty"`
	UpdatedParameters map[string]interface{} `json:"updated_parameters,omitempty"`
}
