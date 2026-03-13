package models


// ProductsProductRetrievalResponse represents a products_product_retrieval_response
type ProductsProductRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	Product *Product `json:"product,omitempty"`
}
