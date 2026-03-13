package models


// ProductRetrievalResponse represents a product_retrieval_response
type ProductRetrievalResponse struct {
	Product *Product `json:"product,omitempty"`
	Error *Error `json:"error,omitempty"`
}
