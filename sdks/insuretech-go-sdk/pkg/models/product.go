package models

import (
	"time"
)

// Product represents a product
type Product struct {
	BasePremiumCurrency string `json:"base_premium_currency,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	MaxSumInsured *Money `json:"max_sum_insured,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	AvailableRiders []*ProductsRider `json:"available_riders,omitempty"`
	MinSumInsuredCurrency string `json:"min_sum_insured_currency,omitempty"`
	MaxSumInsuredCurrency string `json:"max_sum_insured_currency,omitempty"`
	Category *ProductCategory `json:"category,omitempty"`
	MinTenureMonths int `json:"min_tenure_months,omitempty"`
	Exclusions []string `json:"exclusions,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Description string `json:"description,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	MinSumInsured *Money `json:"min_sum_insured,omitempty"`
	MaxTenureMonths int `json:"max_tenure_months,omitempty"`
	Status *ProductsProductStatus `json:"status,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	ProductAttributes string `json:"product_attributes,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	PricingConfig *PricingConfig `json:"pricing_config,omitempty"`
	Plans []*ProductPlan `json:"plans,omitempty"`
}
