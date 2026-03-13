package models


// ProductUpdateResponse represents a product_update_response
type ProductUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
