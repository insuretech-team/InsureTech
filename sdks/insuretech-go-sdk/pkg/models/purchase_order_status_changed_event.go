package models

import (
	"time"
)

// PurchaseOrderStatusChangedEvent represents a purchase_order_status_changed_event
type PurchaseOrderStatusChangedEvent struct {
	OrganisationId string `json:"organisation_id,omitempty"`
	OldStatus *PurchaseOrderStatus `json:"old_status,omitempty"`
	NewStatus *PurchaseOrderStatus `json:"new_status,omitempty"`
	ChangedBy string `json:"changed_by,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
}
