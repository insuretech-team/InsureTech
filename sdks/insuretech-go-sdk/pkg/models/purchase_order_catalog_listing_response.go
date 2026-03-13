package models


// PurchaseOrderCatalogListingResponse represents a purchase_order_catalog_listing_response
type PurchaseOrderCatalogListingResponse struct {
	Items []*PurchaseOrderCatalogItem `json:"items,omitempty"`
	Error *Error `json:"error,omitempty"`
}
