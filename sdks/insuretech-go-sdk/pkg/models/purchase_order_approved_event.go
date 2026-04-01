package models

import (
	"time"
)

// PurchaseOrderApprovedEvent represents a purchase_order_approved_event
type PurchaseOrderApprovedEvent struct {
	ApprovedBy string `json:"approved_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
