package models


// ProductsProductUpdateResponse represents a products_product_update_response
type ProductsProductUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
