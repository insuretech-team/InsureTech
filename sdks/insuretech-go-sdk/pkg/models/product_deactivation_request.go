package models


// ProductDeactivationRequest represents a product_deactivation_request
type ProductDeactivationRequest struct {
	ProductId string `json:"product_id"`
	Reason string `json:"reason,omitempty"`
}
