package models

import (
	"time"
)

// RatingFormulaEvaluationResponse represents a rating_formula_evaluation_response
type RatingFormulaEvaluationResponse struct {
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationId string `json:"calculation_id,omitempty"`
	Errors []string `json:"errors,omitempty"`
	ExecutionTimeMs string `json:"execution_time_ms,omitempty"`
	OutputVariables map[string]interface{} `json:"output_variables,omitempty"`
	Result float64 `json:"result,omitempty"`
	Success bool `json:"success,omitempty"`
}
