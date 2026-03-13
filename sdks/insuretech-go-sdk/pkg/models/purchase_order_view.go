package models


// PurchaseOrderView represents a purchase_order_view
type PurchaseOrderView struct {
	PurchaseOrder *PurchaseOrder `json:"purchase_order,omitempty"`
	DepartmentName string `json:"department_name,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
}
