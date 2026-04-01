package models

import (
	"time"
)

// ProductPlan represents a product_plan
type ProductPlan struct {
	Attributes string `json:"attributes,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	MaxSumInsured *Money `json:"max_sum_insured,omitempty"`
	MinSumInsured *Money `json:"min_sum_insured,omitempty"`
	PlanDescription string `json:"plan_description,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
