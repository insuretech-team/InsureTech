package models


// OrderStatusRetrievalResponse represents a order_status_retrieval_response
type OrderStatusRetrievalResponse struct {
	OrderId string `json:"order_id,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Error *Error `json:"error,omitempty"`
}
