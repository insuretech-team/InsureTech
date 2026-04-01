package models

import (
	"time"
)

// PartnerBenefits represents a partner_benefits
type PartnerBenefits struct {
	AuthorizationValidityDays int `json:"authorization_validity_days,omitempty"`
	AutoApprovalThreshold string `json:"auto_approval_threshold,omitempty"`
	CashlessEnabled bool `json:"cashless_enabled,omitempty"`
	CashlessLimit string `json:"cashless_limit,omitempty"`
	DiscountEnabled bool `json:"discount_enabled,omitempty"`
	DiscountPercentage float64 `json:"discount_percentage,omitempty"`
	DiscountType string `json:"discount_type,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	MaxDiscount float64 `json:"max_discount,omitempty"`
	MinDiscount float64 `json:"min_discount,omitempty"`
	NationwideCoverage bool `json:"nationwide_coverage,omitempty"`
	Notes string `json:"notes,omitempty"`
	PreAuthorizationRequired bool `json:"pre_authorization_required,omitempty"`
	RequiredDocuments []string `json:"required_documents,omitempty"`
	ServiceLocations []string `json:"service_locations,omitempty"`
}
