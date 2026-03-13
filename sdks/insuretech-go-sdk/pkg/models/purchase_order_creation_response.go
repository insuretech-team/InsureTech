package models


// PurchaseOrderCreationResponse represents a purchase_order_creation_response
type PurchaseOrderCreationResponse struct {
	PurchaseOrder *PurchaseOrderView `json:"purchase_order,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
