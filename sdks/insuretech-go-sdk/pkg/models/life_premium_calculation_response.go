package models


// LifePremiumCalculationResponse represents a life_premium_calculation_response
type LifePremiumCalculationResponse struct {
	AgeAddition string `json:"age_addition,omitempty"`
	AppliedBonuses []string `json:"applied_bonuses,omitempty"`
	AppliedConditions []string `json:"applied_conditions,omitempty"`
	BasePremium string `json:"base_premium,omitempty"`
	BonusDiscount string `json:"bonus_discount,omitempty"`
	Breakdown []*LifePremiumBreakdown `json:"breakdown,omitempty"`
	ConditionAddition string `json:"condition_addition,omitempty"`
	ConditionMultiplier float64 `json:"condition_multiplier,omitempty"`
	TotalPremium string `json:"total_premium,omitempty"`
}
