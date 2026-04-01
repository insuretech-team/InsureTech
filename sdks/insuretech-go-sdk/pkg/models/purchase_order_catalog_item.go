package models


// PurchaseOrderCatalogItem represents a purchase_order_catalog_item
type PurchaseOrderCatalogItem struct {
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
}
