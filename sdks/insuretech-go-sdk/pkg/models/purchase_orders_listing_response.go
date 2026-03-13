package models


// PurchaseOrdersListingResponse represents a purchase_orders_listing_response
type PurchaseOrdersListingResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	PurchaseOrders []*PurchaseOrderView `json:"purchase_orders,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}
