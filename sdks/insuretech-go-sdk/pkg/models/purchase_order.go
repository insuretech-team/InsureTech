package models

import (
	"time"
)

// PurchaseOrder represents a purchase_order
type PurchaseOrder struct {
	BusinessId string `json:"business_id,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EmployeeCount int `json:"employee_count,omitempty"`
	EstimatedPremium *Money `json:"estimated_premium,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	Notes string `json:"notes,omitempty"`
	NumberOfDependents int `json:"number_of_dependents,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	PurchaseOrderNumber string `json:"purchase_order_number,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	Status *PurchaseOrderStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
