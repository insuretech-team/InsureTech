package models


// PaymentMethodsListingResponse represents a payment_methods_listing_response
type PaymentMethodsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	PaymentMethods []*PaymentMethodDetails `json:"payment_methods,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
