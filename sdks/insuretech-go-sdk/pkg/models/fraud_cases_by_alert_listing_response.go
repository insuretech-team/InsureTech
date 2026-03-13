package models


// FraudCasesByAlertListingResponse represents a fraud_cases_by_alert_listing_response
type FraudCasesByAlertListingResponse struct {
	Cases []*FraudCase `json:"cases,omitempty"`
}
