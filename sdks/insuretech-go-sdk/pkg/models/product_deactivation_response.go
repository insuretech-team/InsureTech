package models


// ProductDeactivationResponse represents a product_deactivation_response
type ProductDeactivationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
