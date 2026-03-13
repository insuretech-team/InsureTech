package models

import (
	"time"
)

// PurchaseOrder represents a purchase_order
type PurchaseOrder struct {
	PlanId string `json:"plan_id,omitempty"`
	EmployeeCount int `json:"employee_count,omitempty"`
	Status *PurchaseOrderStatus `json:"status,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	NumberOfDependents int `json:"number_of_dependents,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	EstimatedPremium *Money `json:"estimated_premium,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PurchaseOrderNumber string `json:"purchase_order_number,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Notes string `json:"notes,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
