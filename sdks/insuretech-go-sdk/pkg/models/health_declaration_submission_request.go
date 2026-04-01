package models


// HealthDeclarationSubmissionRequest represents a health_declaration_submission_request
type HealthDeclarationSubmissionRequest struct {
	AlcoholConsumer bool `json:"alcohol_consumer,omitempty"`
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	HeightCm int `json:"height_cm,omitempty"`
	OccupationRiskLevel string `json:"occupation_risk_level,omitempty"`
	PreExistingConditions string `json:"pre_existing_conditions,omitempty"`
	QuoteId string `json:"quote_id"`
	Smoker bool `json:"smoker,omitempty"`
	WeightKg string `json:"weight_kg,omitempty"`
}
