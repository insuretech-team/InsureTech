package models


// InsurerInsurerProductUpdateRequest represents a insurer_insurer_product_update_request
type InsurerInsurerProductUpdateRequest struct {
	InsurerProductId string `json:"insurer_product_id"`
	Status string `json:"status,omitempty"`
	EffectiveTo string `json:"effective_to,omitempty"`
	Features string `json:"features,omitempty"`
}
