package models

import (
	"time"
)

// PurchaseOrderCreatedEvent represents a purchase_order_created_event
type PurchaseOrderCreatedEvent struct {
	CreatedBy string `json:"created_by,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EmployeeCount int `json:"employee_count,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
