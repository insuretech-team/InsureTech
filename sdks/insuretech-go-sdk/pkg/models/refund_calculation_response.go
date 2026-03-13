package models


// RefundCalculationResponse represents a refund_calculation_response
type RefundCalculationResponse struct {
	CalculationDetails string `json:"calculation_details,omitempty"`
	Error *Error `json:"error,omitempty"`
	TotalPremiumPaid string `json:"total_premium_paid,omitempty"`
	PremiumUsed string `json:"premium_used,omitempty"`
	CancellationCharge string `json:"cancellation_charge,omitempty"`
	RefundableAmount *Money `json:"refundable_amount,omitempty"`
}
