package models


// FraudAlertsListingResponse represents a fraud_alerts_listing_response
type FraudAlertsListingResponse struct {
	FraudAlerts []*FraudAlert `json:"fraud_alerts,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
