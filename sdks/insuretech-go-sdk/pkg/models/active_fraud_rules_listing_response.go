package models


// ActiveFraudRulesListingResponse represents a active_fraud_rules_listing_response
type ActiveFraudRulesListingResponse struct {
	Rules []*FraudRule `json:"rules,omitempty"`
}
