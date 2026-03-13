package models


// PurchaseOrderCatalogItem represents a purchase_order_catalog_item
type PurchaseOrderCatalogItem struct {
	PlanName string `json:"plan_name,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
}
