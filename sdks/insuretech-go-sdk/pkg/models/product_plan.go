package models

import (
	"time"
)

// ProductPlan represents a product_plan
type ProductPlan struct {
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	MaxSumInsured *Money `json:"max_sum_insured,omitempty"`
	Attributes string `json:"attributes,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	MinSumInsured *Money `json:"min_sum_insured,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	PlanDescription string `json:"plan_description,omitempty"`
}
