package models


// InsurerProductUpdateResponse represents a insurer_product_update_response
type InsurerProductUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
