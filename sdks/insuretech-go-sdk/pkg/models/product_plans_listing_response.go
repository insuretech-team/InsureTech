package models


// ProductPlansListingResponse represents a product_plans_listing_response
type ProductPlansListingResponse struct {
	Plans []*ProductPlan `json:"plans,omitempty"`
}
