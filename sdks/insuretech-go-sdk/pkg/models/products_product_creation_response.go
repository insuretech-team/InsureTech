package models


// ProductsProductCreationResponse represents a products_product_creation_response
type ProductsProductCreationResponse struct {
	ProductId string `json:"product_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
