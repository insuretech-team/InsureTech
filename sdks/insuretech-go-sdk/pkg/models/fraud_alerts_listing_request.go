package models


// FraudAlertsListingRequest represents a fraud_alerts_listing_request
type FraudAlertsListingRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	Status string `json:"status"`
}
