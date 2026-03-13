package models


// OrderCancellationRequest represents a order_cancellation_request
type OrderCancellationRequest struct {
	OrderId string `json:"order_id"`
	Reason string `json:"reason,omitempty"`
}
