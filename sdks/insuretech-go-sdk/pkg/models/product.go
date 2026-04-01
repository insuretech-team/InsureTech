package models

import (
	"time"
)

// Product represents a product
type Product struct {
	AvailableRiders []*ProductsRider `json:"available_riders,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	BasePremiumCurrency string `json:"base_premium_currency,omitempty"`
	Category *ProductCategory `json:"category,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	Exclusions []string `json:"exclusions,omitempty"`
	MaxSumInsured *Money `json:"max_sum_insured,omitempty"`
	MaxSumInsuredCurrency string `json:"max_sum_insured_currency,omitempty"`
	MaxTenureMonths int `json:"max_tenure_months,omitempty"`
	MinSumInsured *Money `json:"min_sum_insured,omitempty"`
	MinSumInsuredCurrency string `json:"min_sum_insured_currency,omitempty"`
	MinTenureMonths int `json:"min_tenure_months,omitempty"`
	Plans []*ProductPlan `json:"plans,omitempty"`
	PricingConfig *PricingConfig `json:"pricing_config,omitempty"`
	ProductAttributes string `json:"product_attributes,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Status *ProductsProductStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
