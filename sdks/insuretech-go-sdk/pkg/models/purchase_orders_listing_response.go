package models


// PurchaseOrdersListingResponse represents a purchase_orders_listing_response
type PurchaseOrdersListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	PurchaseOrders []*PurchaseOrderView `json:"purchase_orders,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
