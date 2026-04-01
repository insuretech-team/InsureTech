package models


// LifeProductDeletionRequest represents a life_product_deletion_request
type LifeProductDeletionRequest struct {
	Permanent bool `json:"permanent,omitempty"`
	ProductId string `json:"product_id"`
}
