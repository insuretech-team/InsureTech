package models


// OrderView represents a order_view
type OrderView struct {
	CustomerName string `json:"customer_name,omitempty"`
	Order *Order `json:"order,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	QuotationNumber string `json:"quotation_number,omitempty"`
}
