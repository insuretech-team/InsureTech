package models

import (
	"time"
)

// InsurerProduct represents a insurer_product
type InsurerProduct struct {
	AuditInfo interface{} `json:"audit_info"`
	MinSumAssured *Money `json:"min_sum_assured,omitempty"`
	MaxSumAssured *Money `json:"max_sum_assured,omitempty"`
	MinPremium *Money `json:"min_premium,omitempty"`
	MaxEntryAge int `json:"max_entry_age,omitempty"`
	MedicalRequired bool `json:"medical_required,omitempty"`
	MedicalThreshold *Money `json:"medical_threshold,omitempty"`
	Id string `json:"id"`
	ProductId string `json:"product_id"`
	Code string `json:"code"`
	Status interface{} `json:"status"`
	MaxPremium *Money `json:"max_premium,omitempty"`
	Features string `json:"features,omitempty"`
	InsurerId string `json:"insurer_id"`
	Name string `json:"name"`
	MaxTermYears int `json:"max_term_years,omitempty"`
	PremiumPaymentModes []string `json:"premium_payment_modes,omitempty"`
	CommissionConfigId string `json:"commission_config_id,omitempty"`
	Exclusions string `json:"exclusions,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	MinEntryAge int `json:"min_entry_age,omitempty"`
	MaxMaturityAge int `json:"max_maturity_age,omitempty"`
	MinTermYears int `json:"min_term_years,omitempty"`
	FreeLookPeriodDays int `json:"free_look_period_days,omitempty"`
}
