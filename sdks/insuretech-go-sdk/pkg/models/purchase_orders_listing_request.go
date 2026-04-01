package models


// PurchaseOrdersListingRequest represents a purchase_orders_listing_request
type PurchaseOrdersListingRequest struct {
	BusinessId string `json:"business_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Status *PurchaseOrderStatus `json:"status,omitempty"`
}
