package models

import (
	"time"
)

// PurchaseOrderStatusChangedEvent represents a purchase_order_status_changed_event
type PurchaseOrderStatusChangedEvent struct {
	ChangedBy string `json:"changed_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus *PurchaseOrderStatus `json:"new_status,omitempty"`
	OldStatus *PurchaseOrderStatus `json:"old_status,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
