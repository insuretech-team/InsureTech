package models


// OrderRetrievalResponse represents a order_retrieval_response
type OrderRetrievalResponse struct {
	Order *OrderView `json:"order,omitempty"`
}
