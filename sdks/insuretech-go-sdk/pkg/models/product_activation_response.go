package models


// ProductActivationResponse represents a product_activation_response
type ProductActivationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
