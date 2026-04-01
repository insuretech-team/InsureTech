package models


// OrderCreationResponse represents a order_creation_response
type OrderCreationResponse struct {
	Order *OrderView `json:"order,omitempty"`
}
