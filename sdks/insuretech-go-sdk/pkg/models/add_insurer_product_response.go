package models


// AddInsurerProductResponse represents a add_insurer_product_response
type AddInsurerProductResponse struct {
	InsurerProductId string `json:"insurer_product_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
