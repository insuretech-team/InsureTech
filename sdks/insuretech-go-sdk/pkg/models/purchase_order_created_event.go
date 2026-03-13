package models

import (
	"time"
)

// PurchaseOrderCreatedEvent represents a purchase_order_created_event
type PurchaseOrderCreatedEvent struct {
	PlanId string `json:"plan_id,omitempty"`
	EmployeeCount int `json:"employee_count,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
}
