package models


// PremiumCalculation represents a premium_calculation
type PremiumCalculation struct {
	BasePremium *Money `json:"base_premium,omitempty"`
	Breakdown []*QuotingPremiumBreakdown `json:"breakdown,omitempty"`
	Currency string `json:"currency,omitempty"`
	DiscountsTotal *Money `json:"discounts_total,omitempty"`
	Fees *Money `json:"fees,omitempty"`
	OptionalCoveragesTotal *Money `json:"optional_coverages_total,omitempty"`
	RiskAdjustment *Money `json:"risk_adjustment,omitempty"`
	Taxes *Money `json:"taxes,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
}
