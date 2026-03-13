package models


// DiscontinueProductResponse represents a discontinue_product_response
type DiscontinueProductResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
