package models


// PurchaseOrdersListingRequest represents a purchase_orders_listing_request
type PurchaseOrdersListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	BusinessId string `json:"business_id"`
	Status *PurchaseOrderStatus `json:"status,omitempty"`
}
