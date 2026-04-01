package models


// ProductsPremiumCalculationResponse represents a products_premium_calculation_response
type ProductsPremiumCalculationResponse struct {
	BasePremium *Money `json:"base_premium,omitempty"`
	Breakdown []*ProductsPremiumBreakdown `json:"breakdown,omitempty"`
	RiderPremium *Money `json:"rider_premium,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
}
