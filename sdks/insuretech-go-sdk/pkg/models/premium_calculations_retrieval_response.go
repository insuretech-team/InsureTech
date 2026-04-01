package models


// PremiumCalculationsRetrievalResponse represents a premium_calculations_retrieval_response
type PremiumCalculationsRetrievalResponse struct {
	Calculations []*VehiclePremiumCalculation `json:"calculations,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
