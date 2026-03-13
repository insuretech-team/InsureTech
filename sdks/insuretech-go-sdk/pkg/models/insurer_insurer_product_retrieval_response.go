package models


// InsurerInsurerProductRetrievalResponse represents a insurer_insurer_product_retrieval_response
type InsurerInsurerProductRetrievalResponse struct {
	InsurerProduct *InsurerProduct `json:"insurer_product,omitempty"`
	Error *Error `json:"error,omitempty"`
}
