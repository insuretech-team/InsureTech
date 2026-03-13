package models


// FraudAlertsByEntityListingResponse represents a fraud_alerts_by_entity_listing_response
type FraudAlertsByEntityListingResponse struct {
	Alerts []*FraudAlert `json:"alerts,omitempty"`
}
