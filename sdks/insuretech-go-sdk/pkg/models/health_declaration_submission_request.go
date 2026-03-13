package models


// HealthDeclarationSubmissionRequest represents a health_declaration_submission_request
type HealthDeclarationSubmissionRequest struct {
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	PreExistingConditions string `json:"pre_existing_conditions,omitempty"`
	Smoker bool `json:"smoker,omitempty"`
	AlcoholConsumer bool `json:"alcohol_consumer,omitempty"`
	OccupationRiskLevel string `json:"occupation_risk_level,omitempty"`
	QuoteId string `json:"quote_id"`
	HeightCm int `json:"height_cm,omitempty"`
	WeightKg string `json:"weight_kg,omitempty"`
}
