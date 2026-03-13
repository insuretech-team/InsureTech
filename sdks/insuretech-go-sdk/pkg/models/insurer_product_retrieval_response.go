package models


// InsurerProductRetrievalResponse represents a insurer_product_retrieval_response
type InsurerProductRetrievalResponse struct {
	InsurerProduct *InsurerProduct `json:"insurer_product,omitempty"`
	Error *Error `json:"error,omitempty"`
}
