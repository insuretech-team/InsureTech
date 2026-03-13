package models

import (
	"time"
)

// PartnerBenefits represents a partner_benefits
type PartnerBenefits struct {
	MaxDiscount float64 `json:"max_discount,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	DiscountType string `json:"discount_type,omitempty"`
	CashlessLimit string `json:"cashless_limit,omitempty"`
	CashlessEnabled bool `json:"cashless_enabled,omitempty"`
	AutoApprovalThreshold string `json:"auto_approval_threshold,omitempty"`
	PreAuthorizationRequired bool `json:"pre_authorization_required,omitempty"`
	ServiceLocations []string `json:"service_locations,omitempty"`
	Notes string `json:"notes,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	DiscountEnabled bool `json:"discount_enabled,omitempty"`
	DiscountPercentage float64 `json:"discount_percentage,omitempty"`
	MinDiscount float64 `json:"min_discount,omitempty"`
	AuthorizationValidityDays int `json:"authorization_validity_days,omitempty"`
	RequiredDocuments []string `json:"required_documents,omitempty"`
	NationwideCoverage bool `json:"nationwide_coverage,omitempty"`
}
