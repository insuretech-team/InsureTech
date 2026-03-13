package models


// ProductCreationResponse represents a product_creation_response
type ProductCreationResponse struct {
	ProductId string `json:"product_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
