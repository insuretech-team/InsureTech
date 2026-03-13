package models


// OrderCancellationResponse represents a order_cancellation_response
type OrderCancellationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
}
