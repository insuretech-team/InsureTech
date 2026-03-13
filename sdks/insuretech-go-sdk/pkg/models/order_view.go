package models


// OrderView represents a order_view
type OrderView struct {
	ProductName string `json:"product_name,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	QuotationNumber string `json:"quotation_number,omitempty"`
	Order *Order `json:"order,omitempty"`
}
