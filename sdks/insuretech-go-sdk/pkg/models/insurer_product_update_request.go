package models


// InsurerProductUpdateRequest represents a insurer_product_update_request
type InsurerProductUpdateRequest struct {
	EffectiveTo string `json:"effective_to,omitempty"`
	Features string `json:"features,omitempty"`
	InsurerProductId string `json:"insurer_product_id"`
	Status string `json:"status,omitempty"`
}
