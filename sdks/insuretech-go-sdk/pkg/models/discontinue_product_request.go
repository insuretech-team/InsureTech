package models


// DiscontinueProductRequest represents a discontinue_product_request
type DiscontinueProductRequest struct {
	ProductId string `json:"product_id"`
	Reason string `json:"reason,omitempty"`
}
