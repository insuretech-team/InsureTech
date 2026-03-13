package models


// PurchaseOrderRetrievalResponse represents a purchase_order_retrieval_response
type PurchaseOrderRetrievalResponse struct {
	PurchaseOrder *PurchaseOrderView `json:"purchase_order,omitempty"`
	Error *Error `json:"error,omitempty"`
}
