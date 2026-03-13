package models


// PaymentMethodsListingResponse represents a payment_methods_listing_response
type PaymentMethodsListingResponse struct {
	PaymentMethods []*PaymentPaymentMethod `json:"payment_methods,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
