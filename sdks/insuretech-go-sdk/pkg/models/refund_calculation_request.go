package models


// RefundCalculationRequest represents a refund_calculation_request
type RefundCalculationRequest struct {
	PolicyId string `json:"policy_id"`
	Reason string `json:"reason,omitempty"`
}
