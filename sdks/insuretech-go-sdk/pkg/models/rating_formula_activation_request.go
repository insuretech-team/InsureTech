package models

import (
	"time"
)

// RatingFormulaActivationRequest represents a rating_formula_activation_request
type RatingFormulaActivationRequest struct {
	ActivatedBy string `json:"activated_by,omitempty"`
	EffectiveDate time.Time `json:"effective_date,omitempty"`
	FormulaId string `json:"formula_id"`
}
