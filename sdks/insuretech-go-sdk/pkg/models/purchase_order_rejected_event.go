package models

import (
	"time"
)

// PurchaseOrderRejectedEvent represents a purchase_order_rejected_event
type PurchaseOrderRejectedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	RejectedBy string `json:"rejected_by,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
