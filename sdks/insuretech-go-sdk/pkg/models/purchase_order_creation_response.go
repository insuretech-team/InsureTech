package models


// PurchaseOrderCreationResponse represents a purchase_order_creation_response
type PurchaseOrderCreationResponse struct {
	PurchaseOrder *PurchaseOrderView `json:"purchase_order,omitempty"`
}
