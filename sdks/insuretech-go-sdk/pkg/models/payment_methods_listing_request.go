package models


// PaymentMethodsListingRequest represents a payment_methods_listing_request
type PaymentMethodsListingRequest struct {
	UserId string `json:"user_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
