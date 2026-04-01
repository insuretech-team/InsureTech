package models

import (
	"time"
)

// InsurerProduct represents a insurer_product
type InsurerProduct struct {
	AuditInfo interface{} `json:"audit_info"`
	Code string `json:"code"`
	CommissionConfigId string `json:"commission_config_id,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	Exclusions string `json:"exclusions,omitempty"`
	Features string `json:"features,omitempty"`
	FreeLookPeriodDays int `json:"free_look_period_days,omitempty"`
	Id string `json:"id"`
	InsurerId string `json:"insurer_id"`
	MaxEntryAge int `json:"max_entry_age,omitempty"`
	MaxMaturityAge int `json:"max_maturity_age,omitempty"`
	MaxPremium *Money `json:"max_premium,omitempty"`
	MaxSumAssured *Money `json:"max_sum_assured,omitempty"`
	MaxTermYears int `json:"max_term_years,omitempty"`
	MedicalRequired bool `json:"medical_required,omitempty"`
	MedicalThreshold *Money `json:"medical_threshold,omitempty"`
	MinEntryAge int `json:"min_entry_age,omitempty"`
	MinPremium *Money `json:"min_premium,omitempty"`
	MinSumAssured *Money `json:"min_sum_assured,omitempty"`
	MinTermYears int `json:"min_term_years,omitempty"`
	Name string `json:"name"`
	PremiumPaymentModes []string `json:"premium_payment_modes,omitempty"`
	ProductId string `json:"product_id"`
	Status interface{} `json:"status"`
}
