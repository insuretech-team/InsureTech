package models


// PurchaseOrderView represents a purchase_order_view
type PurchaseOrderView struct {
	DepartmentName string `json:"department_name,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	PurchaseOrder *PurchaseOrder `json:"purchase_order,omitempty"`
}
