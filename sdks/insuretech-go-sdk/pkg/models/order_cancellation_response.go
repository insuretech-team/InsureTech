package models


// OrderCancellationResponse represents a order_cancellation_response
type OrderCancellationResponse struct {
	OrderId string `json:"order_id,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
}
