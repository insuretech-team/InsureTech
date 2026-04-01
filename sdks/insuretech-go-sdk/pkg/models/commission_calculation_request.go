package models


// CommissionCalculationRequest represents a commission_calculation_request
type CommissionCalculationRequest struct {
	CommissionType string `json:"commission_type,omitempty"`
	PolicyId string `json:"policy_id"`
	RecipientId string `json:"recipient_id"`
	RecipientType string `json:"recipient_type,omitempty"`
}
