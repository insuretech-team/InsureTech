package models


// InsurerInsurerProductUpdateResponse represents a insurer_insurer_product_update_response
type InsurerInsurerProductUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
